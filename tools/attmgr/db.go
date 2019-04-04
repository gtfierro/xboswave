package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/fsnotify/fsnotify"
	"github.com/immesys/wave/eapi/pb"
	"github.com/pkg/errors"
)

type DB struct {
	db          *sql.DB
	wave        pb.WAVEClient
	perspective *pb.Perspective
	//base64
	phash string
	shell *ishell.Shell
}

func NewDB(cfg *Config) (*DB, error) {
	var err error

	db := &DB{
		wave:        getConn(cfg.Agent),
		perspective: getPerspective(cfg.Perspective, "", "missing perspective"),
	}

	// get entity hash
	resp, err := db.wave.Inspect(context.Background(), &pb.InspectParams{
		Content: db.perspective.EntitySecret.DER,
	})
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("could not inspect file: %s\n", resp.Error.Message)
	}
	if resp.Entity == nil {
		return nil, fmt.Errorf("file was not an entity %s\n", cfg.Perspective)
	}
	db.phash = base64.URLEncoding.EncodeToString(resp.Entity.Hash)

	db.db, err = sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, err
	}

	// set up tables
	_, err = db.db.Exec(`CREATE TABLE IF NOT EXISTS attestations (
        id          INTEGER PRIMARY KEY,
        hash        TEXT UNIQUE NOT NULL,
        subject		TEXT NOT NULL,
        inserted    DATETIME DEFAULT CURRENT_TIMESTAMP,
        expires     DATETIME NOT NULL,
        policies    JSON
    );`)
	_, err = db.db.Exec(`CREATE TABLE IF NOT EXISTS policies (
        id          INTEGER PRIMARY KEY,
        namespace   TEXT NOT NULL,
        resource    TEXT NOT NULL,
        pset        TEXT NOT NULL,
        indirections INTEGER NOT NULL,
        permissions JSON
    );`)
	if err != nil {
		return nil, err
	}

	// setup interactive shell
	db.setupShell()

	return db, nil
}

func (db *DB) watch(dir string) {
	// load in existing
	files, err := filepath.Glob("*.pem")
	if err != nil {
		log.Fatal(err)
	}
	for _, filename := range files {
		log.Println("load detected file: ", db.LoadAttestationFile(filename))
	}

	// watch the directory!
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Create && strings.HasSuffix(event.Name, ".pem") {
					log.Println("Loading new .pem file ", event.Name)
					log.Println("load detected file: ", db.LoadAttestationFile(event.Name))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func (db *DB) insertAttestation(att *Attestation) error {

	if att == nil {
		return errors.New("Could not insert attestation that was not complete or fully decoded")
	}

	// check attesterhash
	//if att.Attester

	tx, err := db.db.BeginTx(context.Background(), &sql.TxOptions{
		ReadOnly: false,
	})
	if err != nil {
		return err
	}

	var ids []int
	//	// insert policy statements
	stmt := "INSERT INTO policies(namespace, resource, pset, indirections, permissions) VALUES (?, ?, ?, ?, ?)"
	for _, ps := range att.PolicyStatements {
		// sort permissions for consistency
		var s = sort.StringSlice(ps.Permissions)
		s.Sort()
		b, err := json.Marshal(s)
		if err != nil {
			return err
		}

		// see if policy already exists
		row := tx.QueryRow(formQueryStr("policies", "id", map[string]string{
			"namespace":    ps.Namespace,
			"resource":     ps.Resource,
			"pset":         ps.PermissionSet,
			"indirections": fmt.Sprintf("%d", ps.Indirections),
			"permissions":  string(b),
		}))
		var id int
		if err := row.Scan(&id); err != nil && err != sql.ErrNoRows {
			log.Println(errors.Wrap(err, "query policy"))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
			}
			return err
		} else if err == nil {
			ids = append(ids, id)
			continue
		}

		_, err = tx.Exec(stmt, ps.Namespace, ps.Resource, ps.PermissionSet, ps.Indirections, string(b))
		if err != nil {
			log.Println(errors.Wrap(err, "upsert policy"))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}

		row = tx.QueryRow(formQueryStr("policies", "id", map[string]string{
			"namespace":    ps.Namespace,
			"resource":     ps.Resource,
			"pset":         ps.PermissionSet,
			"indirections": fmt.Sprintf("%d", ps.Indirections),
			"permissions":  string(b),
		}))
		if err := row.Scan(&id); err != nil {
			log.Println(errors.Wrap(err, "query policy id"))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}
		ids = append(ids, id)
	}
	fmt.Println(ids)

	// insert attestations
	stmt = "INSERT INTO attestations(hash, expires, policies, subject) VALUES (?, ?, ?, ?) ON CONFLICT(hash) DO UPDATE SET policies=json_patch(policies, '%s')"
	for _, id := range ids {
		pol, err := json.Marshal(map[int]int{id: 0})
		if err != nil {
			log.Println(errors.Wrap(err, "insert att"))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}

		//pol := fmt.Sprintf("'{%d: 0}'", id)
		fmt.Printf(stmt+"\n", string(pol))
		fmt.Println(pol)

		_, err = tx.Exec(fmt.Sprintf(stmt, string(pol)), att.Hash, att.ValidUntil, string(pol), att.Subject)
		if err != nil {
			log.Println(errors.Wrap(err, "insert att"))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}
	}
	return tx.Commit()
}

func (db *DB) getUnique(attribute string, where map[string]string) ([]string, error) {
	stmt := formQueryStr("attestations", attribute, where)
	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	var res []string
	defer rows.Close()
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return res, err
		}
		res = append(res, s)
	}

	return res, nil
}

func (db *DB) getUniqueFromPolicy(attribute string, where map[string]string) ([]string, error) {
	stmt := formQueryStr("policies", attribute, where)
	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	var res []string
	defer rows.Close()
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return res, err
		}
		res = append(res, s)
	}

	return res, nil
}

func formQueryStr(table, attribute string, where map[string]string) string {
	stmt := fmt.Sprintf("SELECT distinct %s FROM %s", attribute, table)
	var filters []string
	if len(where) > 0 {
		stmt += " WHERE "
		for attribute, filter := range where {
			filters = append(filters, fmt.Sprintf(" %s='%s'", attribute, filter))
		}
		stmt += strings.Join(filters, " AND ")
	}
	return stmt
}

func (db *DB) readAttestations(rows *sql.Rows) ([]Attestation, error) {
	var ret []Attestation
	for rows.Next() {
		att := &Attestation{}
		var expires interface{}
		var policyids map[int]int
		var _policyids []byte
		if err := rows.Scan(&att.Hash, &att.Subject, &expires, &_policyids); err != nil {
			return nil, err
		}
		if expires != nil {
			att.ValidUntil = expires.(time.Time)
		}
		if err := json.Unmarshal(_policyids, &policyids); err != nil {
			return nil, err
		}
		fmt.Println(policyids)
		for policyid := range policyids {
			policies, err := db.listPolicy(&filter{polid: &policyid})
			if err != nil {
				return nil, err
			}
			att.PolicyStatements = append(att.PolicyStatements, policies...)
		}
		ret = append(ret, *att)

	}
	return ret, nil
}

func (db *DB) readPolicies(rows *sql.Rows) ([]PolicyStatement, error) {
	var ret []PolicyStatement
	for rows.Next() {
		pol := &PolicyStatement{}
		var _perm []byte
		if err := rows.Scan(&pol.Namespace, &pol.Resource, &pol.PermissionSet, &pol.Indirections, &_perm); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(_perm, &pol.Permissions); err != nil {
			return nil, err
		}
		ret = append(ret, *pol)
	}
	return ret, nil
}

func (db *DB) RunShell() {
	db.shell.Run()
}

func (db *DB) listAttestation(filter *filter) ([]Attestation, error) {
	stmt := `SELECT hash, subject, expires, policies
			 FROM attestations, json_each(policies)
			 LEFT JOIN policies ON policies.id = json_each.value
			 `

	where, err := filter.toSQL()
	if err != nil {
		return nil, err
	}
	if len(where) > 0 {
		stmt += " WHERE " + where
	}
	fmt.Println(stmt)

	rows, err := db.db.Query(stmt)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return db.readAttestations(rows)
}

func (db *DB) listPolicy(filter *filter) ([]PolicyStatement, error) {
	stmt := `SELECT namespace, resource, pset, indirections, permissions
			 FROM policies
			 `

	where, err := filter.toSQL()
	if err != nil {
		return nil, err
	}
	if len(where) > 0 {
		stmt += " WHERE " + where
	}
	fmt.Println(stmt)

	rows, err := db.db.Query(stmt)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return db.readPolicies(rows)
}

type filter struct {
	polid       *int
	attid       *int
	hash        *string
	policy      *int
	namespace   *string
	resource    *string
	pset        *string
	permissions []string
}

func (f *filter) toSQL() (string, error) {
	var filters []string
	if f == nil {
		return "", nil
	}

	if f.attid != nil {
		filters = append(filters, fmt.Sprintf("attestations.id=%d ", *f.attid))
	}
	if f.polid != nil {
		filters = append(filters, fmt.Sprintf("policies.id=%d ", *f.polid))
	}
	if f.hash != nil {
		filters = append(filters, fmt.Sprintf("attestations.hash='%s' ", *f.hash))
	}
	if f.policy != nil {
		filters = append(filters, fmt.Sprintf("attestations.policy=%d ", *f.policy))
	}
	if f.namespace != nil {
		filters = append(filters, fmt.Sprintf("policies.namespace='%s' ", *f.namespace))
	}
	if f.resource != nil {
		filters = append(filters, fmt.Sprintf("policies.resource='%s' ", *f.resource))
	}
	if f.pset != nil {
		filters = append(filters, fmt.Sprintf("policies.pset='%s' ", *f.pset))
	}

	if f.permissions != nil {
		// sort permissions for consistency
		var s = sort.StringSlice(f.permissions)
		s.Sort()
		b, err := json.Marshal(s)
		if err != nil {
			return "", err
		}
		filters = append(filters, fmt.Sprintf("policies.permissions='%s' ", string(b)))
	}
	return strings.Join(filters, " AND "), nil
}
