package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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
	if err != nil {
		return nil, err
	}
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
	_, err = db.db.Exec(`CREATE TABLE IF NOT EXISTS entities (
        id          INTEGER PRIMARY KEY,
        hash        TEXT UNIQUE NOT NULL,
        expires     DATETIME NOT NULL,
        name        TEXT
    );`)
	if err != nil {
		return nil, err
	}

	// resolve any names in the tables if we have new names
	db.resolveHashesToNames()
	db.removeDups()
	// TODO: load in default entity
	if err := db.LoadEntityFile(cfg.Perspective); err != nil {
		return nil, errors.Wrap(err, "Could not insert default entity")
	}

	// setup interactive shell
	db.setupShell()

	return db, nil
}

func (db *DB) removeDups() {
	var pending []*PolicyStatement
	stmt := `select count(*), namespace, pset, resource, permissions, indirections from policies group by namespace, pset, resource, permissions, indirections;`
	rows, err := db.db.Query(stmt)
	if err != nil {
		log.Error(err)
		rows.Close()
		return
	}

	for rows.Next() {
		var count int
		pol := &PolicyStatement{}
		var _perm []byte
		if err := rows.Scan(&count, &pol.Namespace, &pol.PermissionSet, &pol.Resource, &_perm, &pol.Indirections); err != nil {
			log.Error(err)
			continue
		}
		if err := json.Unmarshal(_perm, &pol.Permissions); err != nil {
			log.Error(err)
			continue
		}
		if count > 1 {
			fmt.Println(count, "dups for", pol)
			pending = append(pending, pol)
		}
	}
	rows.Close()

	for _, pol := range pending {
		fmt.Println(pol.Namespace, pol.PermissionSet, pol.Resource, pol.permissionString(), pol.Indirections)
		rows, err := db.db.Query("select id from policies where namespace=? AND pset=? AND resource=? AND permissions=? AND indirections=?;", pol.Namespace, pol.PermissionSet, pol.Resource, pol.permissionString(), pol.Indirections)

		//rows, err := db.db.Query("select id from policies where namespace=? and resource=? and permissions=?;", pol.Namespace, pol.Resource, pol.permissionString())
		if err != nil {
			log.Error(err)
			rows.Close()
			return
		}
		var ids []int
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				log.Error(err)
				continue
			}
			ids = append(ids, id)
		}
		rows.Close()

		new_pol := ids[0]
		for _, old_pol := range ids[1:] {
			_, err = db.db.Exec("delete from policies where id=?", old_pol)
			if err != nil {
				log.Error(err)
				continue
			}
			stmt := fmt.Sprintf("update attestations set policies=json_insert(json_remove(policies, '$.%d'), '$.%d', 0) where json_extract(policies, '$.%d')!='';", old_pol, new_pol, old_pol)
			_, err = db.db.Exec(stmt)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

func (db *DB) resolveHashesToNames() {
	db.syncgraph()
	stmt := `SELECT id, subject from attestations;`
	rows, err := db.db.Query(stmt)
	if err != nil {
		log.Error(err)
		rows.Close()
		return
	}
	var updateatts []struct {
		subject string
		id      int
	}
	for rows.Next() {
		att := &Attestation{}
		if err := rows.Scan(&att.id, &att.Subject); err != nil {
			log.Error(err)
			continue
		}
		name := db.getNameFromHash(att.Subject)
		if name != att.Subject {
			updateatts = append(updateatts, struct {
				subject string
				id      int
			}{name, att.id})
		}
	}
	rows.Close()

	for _, update := range updateatts {
		_, err := db.db.Exec("UPDATE attestations SET subject=? WHERE id=?", update.subject, update.id)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	var updatepols []struct {
		ns   string
		pset string
		id   int
	}

	stmt = `SELECT id, namespace, pset from policies;`
	prows, err := db.db.Query(stmt)
	defer prows.Close()
	if err != nil {
		log.Error(err)
		return
	}
	for prows.Next() {
		pol := &PolicyStatement{}
		if err := prows.Scan(&pol.id, &pol.Namespace, &pol.PermissionSet); err != nil {
			log.Error(err)
			continue
		}
		ns := db.getNameFromHash(pol.Namespace)
		pset := db.getNameFromHash(pol.PermissionSet)
		if ns != pol.Namespace || pset != pol.PermissionSet {
			updatepols = append(updatepols, struct {
				ns   string
				pset string
				id   int
			}{ns, pset, pol.id})
		}
	}
	for _, update := range updatepols {
		_, err := db.db.Exec("UPDATE policies SET namespace=?, pset=? WHERE id=?", update.ns, update.pset, update.id)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	var updateents []struct {
		name string
		id   int
	}

	stmt = `SELECT id, name, hash from entities;`
	erows, err := db.db.Query(stmt)
	defer erows.Close()
	if err != nil {
		log.Error(err)
		return
	}
	for erows.Next() {
		ent := &Entity{}
		if err := erows.Scan(&ent.id, &ent.Name, &ent.Hash); err != nil {
			log.Error(err)
			continue
		}
		n := db.getNameFromHash(ent.Hash)
		if n != ent.Name {
			updateents = append(updateents, struct {
				name string
				id   int
			}{n, ent.id})
		}
	}
	for _, update := range updateents {
		_, err := db.db.Exec("UPDATE entities SET name=? WHERE id=?", update.name, update.id)
		if err != nil {
			log.Error(err)
			continue
		}
	}

}

func (db *DB) watch(dir string) {
	// load in existing
	files, err := filepath.Glob("*.pem")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not glob *.pem"))
	}
	for _, filename := range files {
		if err := db.LoadAttestationFile(filename); err != nil {
			log.Warning("Could not load detected attestation: ", err)
		}
	}

	// load in existing entities
	files, err = filepath.Glob("*.ent")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not glob *.ent"))
	}
	for _, filename := range files {
		if err := db.LoadEntityFile(filename); err != nil {
			log.Warning("Could not load detected entity: ", err)
		}
	}

	// watch the directory!
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not create new file watcher"))
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
					if err := db.LoadAttestationFile(event.Name); err != nil {
						log.Warning("Could not load detected attestation: ", err)
					}
				}
				if event.Op == fsnotify.Create && strings.HasSuffix(event.Name, ".ent") {
					if err := db.LoadEntityFile(event.Name); err != nil {
						log.Warning("Could not load detected entity: ", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not add dir to file watch"))
	}
	<-done
}

func (db *DB) insertPolicy(pol *PolicyStatement) error {
	if pol == nil {
		return errors.New("could not insert empty policy")
	}
	if pol.Namespace == "" {
		return errors.New("Policy needs namespace")
	}
	if pol.PermissionSet == "" {
		return errors.New("Policy needs PermissionSet")
	}
	if len(pol.Permissions) == 0 {
		return errors.New("Policy needs permissions")
	}

	tx, err := db.db.BeginTx(context.Background(), &sql.TxOptions{
		ReadOnly: false,
	})
	if err != nil {
		return err
	}
	// sort permissions for consistency
	var s = sort.StringSlice(pol.Permissions)
	s.Sort()
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// resolve names (if any)
	pol.Namespace = db.getNameFromHash(pol.Namespace)
	pol.PermissionSet = db.getNameFromHash(pol.PermissionSet)

	row := tx.QueryRow(formQueryStr("policies", "id", map[string]string{
		"namespace":    pol.Namespace,
		"resource":     pol.Resource,
		"pset":         pol.PermissionSet,
		"indirections": fmt.Sprintf("%d", pol.Indirections),
		"permissions":  string(b),
	}))
	var id int
	if err := row.Scan(&id); err != nil && err != sql.ErrNoRows {
		log.Error(errors.Wrap(err, "query policy"))
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
		}
		return err
	} else if err == nil {
		return tx.Commit()
	}
	stmt := "INSERT INTO policies(namespace, resource, pset, indirections, permissions) VALUES (?, ?, ?, ?, ?)"
	_, err = tx.Exec(stmt, pol.Namespace, pol.Resource, pol.PermissionSet, pol.Indirections, string(b))
	if err != nil {
		log.Error(errors.Wrap(err, "upsert policy"))
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
		}
	}

	return tx.Commit()
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

		ps.Namespace = db.getNameFromHash(ps.Namespace)
		ps.PermissionSet = db.getNameFromHash(ps.PermissionSet)

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
			log.Error(errors.Wrap(err, "query policy"))
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
			log.Error(errors.Wrap(err, "upsert policy"))
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
			log.Error(errors.Wrap(err, "query policy id"))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}
		ids = append(ids, id)
	}

	// insert attestations
	stmt = "INSERT INTO attestations(hash, expires, policies, subject) VALUES (?, ?, ?, ?) ON CONFLICT(hash) DO UPDATE SET policies=json_patch(policies, '%s')"
	for _, id := range ids {
		pol, err := json.Marshal(map[int]int{id: 0})
		if err != nil {
			log.Error(errors.Wrap(err, "insert att"))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}

		fmt.Printf("Resolve subject from %s to %s\n", att.Subject, db.getNameFromHash(att.Subject))
		att.Subject = db.getNameFromHash(att.Subject)

		_, err = tx.Exec(fmt.Sprintf(stmt, string(pol)), att.Hash, att.ValidUntil, string(pol), att.Subject)
		if err != nil {
			log.Error(errors.Wrap(err, "insert att"))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}
	}
	return tx.Commit()
}

func (db *DB) insertEntity(ent *Entity) error {
	if ent == nil {
		return errors.New("could not insert empty entity")
	}
	if ent.Hash == "" {
		return errors.New("Entity needs hash")
	}

	tx, err := db.db.BeginTx(context.Background(), &sql.TxOptions{
		ReadOnly: false,
	})
	if err != nil {
		return err
	}

	// check if ent exists

	// resolve names (if any)
	if ent.Name == "" {
		ent.Name = db.getNameFromHash(ent.Hash)
	}

	stmt := "INSERT OR IGNORE INTO entities(name, hash, expires) VALUES (?, ?, ?)"
	_, err = tx.Exec(stmt, ent.Name, ent.Hash, ent.Expires)
	if err != nil {
		log.Error(errors.Wrap(err, "upsert policy"))
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("update drivers: unable to rollback: %v", rollbackErr)
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
	var seenhashes = make(map[string]struct{})
	for rows.Next() {
		att := &Attestation{}
		var expires interface{}
		var policyids map[int]int
		var _policyids []byte
		if err := rows.Scan(&att.Hash, &att.Subject, &expires, &_policyids); err != nil {
			return nil, err
		}
		if _, found := seenhashes[att.Hash]; found {
			continue
		}
		seenhashes[att.Hash] = struct{}{}

		if expires != nil {
			att.ValidUntil = expires.(time.Time)
		}
		if err := json.Unmarshal(_policyids, &policyids); err != nil {
			return nil, err
		}
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
		if err := rows.Scan(&pol.id, &pol.Namespace, &pol.Resource, &pol.PermissionSet, &pol.Indirections, &_perm); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(_perm, &pol.Permissions); err != nil {
			return nil, err
		}
		ret = append(ret, *pol)
	}
	return ret, nil
}

func (db *DB) readEntities(rows *sql.Rows) ([]Entity, error) {
	var ret []Entity
	for rows.Next() {
		ent := &Entity{}
		var expires interface{}
		if err := rows.Scan(&ent.id, &ent.Name, &ent.Hash, &expires); err != nil {
			return nil, err
		}
		if expires != nil {
			ent.Expires = expires.(time.Time)
		}
		ret = append(ret, *ent)
	}
	return ret, nil
}

func (db *DB) getPoliciesById(ids []int) ([]PolicyStatement, error) {

	var stmts []PolicyStatement
	for _, policyid := range ids {
		policies, err := db.listPolicy(&filter{polid: &policyid})
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, policies...)
	}

	return stmts, nil
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
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return db.readAttestations(rows)
}

func (db *DB) listPolicy(filter *filter) ([]PolicyStatement, error) {
	stmt := `SELECT id, namespace, resource, pset, indirections, permissions
			 FROM policies
			 `

	where, err := filter.toSQL()
	if err != nil {
		return nil, err
	}
	if len(where) > 0 {
		stmt += " WHERE " + where
	}

	rows, err := db.db.Query(stmt)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return db.readPolicies(rows)
}

func (db *DB) listEntity(filter *filter) ([]Entity, error) {
	stmt := `SELECT id, name, hash, expires
			 FROM entities
			 `

	where, err := filter.toSQL()
	if err != nil {
		return nil, err
	}
	if len(where) > 0 {
		stmt += " WHERE " + where
	}

	rows, err := db.db.Query(stmt)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return db.readEntities(rows)
}

type filter struct {
	polid           *int
	attid           *int
	hash            *string
	policy          *int
	namespace       *string
	resource        *string
	pset            *string
	expiring_before *time.Time
	permissions     []string
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
	if f.expiring_before != nil {
		today := time.Now().Format("2006-01-02 15:04:05")
		d := (*f.expiring_before).Format("2006-01-02 15:04:05")
		filters = append(filters, fmt.Sprintf("attestations.expires BETWEEN '%s' and '%s' ", today, d))
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
