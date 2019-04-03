package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/immesys/wave/eapi/pb"
	"github.com/pkg/errors"
)

type DB struct {
	db          *sql.DB
	wave        pb.WAVEClient
	perspective *pb.Perspective
}

func NewDB(cfg *Config) (*DB, error) {
	var err error

	db := &DB{
		wave:        getConn(cfg.Agent),
		perspective: getPerspective(cfg.Perspective, "", "missing perspective"),
	}
	db.db, err = sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, err
	}

	// set up tables
	_, err = db.db.Exec(`CREATE TABLE IF NOT EXISTS attestations (
        id          INTEGER PRIMARY KEY,
        hash        TEXT UNIQUE NOT NULL,
        inserted    DATETIME DEFAULT CURRENT_TIMESTAMP,
        expires     DATETIME NOT NULL,
        policy      INTEGER NOT NULL
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

	return db, nil
}

func (db *DB) watch(dir string) {
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
				if event.Op == fsnotify.Create {
					if strings.HasSuffix(event.Name, ".pem") {
						log.Println("Loading new .pem file ", event.Name)
						log.Println("load detected file: ", db.LoadAttestationFile(event.Name))
					}
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

func (db *DB) insertAttestation(att Attestation) error {

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
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
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
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
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
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}
		ids = append(ids, id)
	}
	fmt.Println(ids)

	// insert attestations
	stmt = "INSERT OR IGNORE INTO attestations(hash, expires, policy) VALUES (?, ?, ?)"
	for _, id := range ids {
		_, err = tx.Exec(stmt, att.Hash, att.ValidUntil, id)
		if err != nil {
			log.Println(errors.Wrap(err, "insert att"))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
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

func getAttestations(*sql.Rows) ([]Attestation, error) {
	return nil, nil
}
