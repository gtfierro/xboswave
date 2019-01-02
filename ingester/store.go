package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gtfierro/xboswave/ingester/types"
	"github.com/immesys/wavemq/mqpb"
	_ "github.com/mattn/go-sqlite3"
)

// these are addressable true/false values used internally for the RequestFilter
var _FALSE = false
var _TRUE = true

type ArchiveRequest struct {
	// message schema this request applies to
	// (this is a field in the WAVEMQ wrapper)
	Schema string
	// path to the plugin to use to extract
	// timeseries data (e.g. plugins/dent.so)
	// Path should end in .so
	Plugin string
	// the URI we subscribe to
	URI types.SubscriptionURI

	// reported values

	// time this archive request was created
	Inserted time.Time
	// the text of the last error this archive request experienced
	LastError string
	// what time that error occured
	ErrorTimestamp time.Time
	// whether or not this archive request is active
	Enabled bool
	// unique identifier
	Id int
}

type subscription struct {
	S    mqpb.WAVEMQ_SubscribeClient
	stop chan struct{}
	uri  types.SubscriptionURI
}

type ConfigManager struct {
	db *sql.DB
}

func NewCfgManager(cfg *Config) (*ConfigManager, error) {
	var err error
	cfgmgr := &ConfigManager{}
	cfgmgr.db, err = sql.Open("sqlite3", cfg.Store.Path)
	if err != nil {
		return nil, err
	}

	// set up tables
	_, err = cfgmgr.db.Exec(`CREATE TABLE IF NOT EXISTS requests (
        id  INTEGER PRIMARY KEY,
        schema  TEXT NOT NULL,
        plugin  TEXT NOT NULL,
        namespace TEXT NOT NULL,
        resource TEXT NOT NULL,
        inserted DATETIME DEFAULT CURRENT_TIMESTAMP,
        lastError TEXT,
        enabled BOOLEAN,
        errorTimestamp DATETIME DEFAULT NULL
    );`)

	return cfgmgr, err
}

type RequestFilter struct {
	Schema    *string
	Plugin    *string
	Namespace *string
	Resource  *string
	HasError  *bool
	Enabled   *bool
	Id        *int
}

func (cfgmgr *ConfigManager) Add(req ArchiveRequest) error {
	filter := &RequestFilter{
		Schema:    &req.Schema,
		Plugin:    &req.Plugin,
		Namespace: &req.URI.Namespace,
		Resource:  &req.URI.Resource,
	}
	results, err := cfgmgr.List(filter)
	if err != nil {
		return err
	} else if len(results) == 0 {
		stmt := "INSERT INTO requests(schema, plugin, namespace, resource, enabled) VALUES (?, ?, ?, ?, 1);"
		_, err := cfgmgr.db.Exec(stmt, req.Schema, req.Plugin, req.URI.Namespace, req.URI.Resource)
		return err
	}

	return nil
}

// deletes the request from the config manager; returns true if anyone else uses the same subscription URI
func (cfgmgr *ConfigManager) Delete(req ArchiveRequest) (bool, error) {

	stmt := "DELETE FROM requests WHERE schema=? AND plugin=? AND namespace=? AND resource=?;"
	_, err := cfgmgr.db.Exec(stmt, req.Schema, req.Plugin, req.URI.Namespace, req.URI.Resource)
	if err != nil {
		return true, err
	}

	filter := &RequestFilter{
		Namespace: &req.URI.Namespace,
		Resource:  &req.URI.Resource,
	}
	existingSubs, err := cfgmgr.List(filter)
	return len(existingSubs) > 0, err
}

func (cfgmgr *ConfigManager) Disable(req ArchiveRequest) (bool, error) {

	stmt := "UPDATE requests SET enabled = 0 WHERE schema=? AND plugin=? AND namespace=? AND resource=?;"
	_, err := cfgmgr.db.Exec(stmt, req.Schema, req.Plugin, req.URI.Namespace, req.URI.Resource)
	if err != nil {
		return true, err
	}

	filter := &RequestFilter{
		Namespace: &req.URI.Namespace,
		Resource:  &req.URI.Resource,
		Enabled:   &_TRUE,
	}
	existingSubs, err := cfgmgr.List(filter)
	return len(existingSubs) > 0, err
}

func (cfgmgr *ConfigManager) Enable(req ArchiveRequest) error {

	stmt := "UPDATE requests SET enabled = 1 WHERE schema=? AND plugin=? AND namespace=? AND resource=?;"
	_, err := cfgmgr.db.Exec(stmt, req.Schema, req.Plugin, req.URI.Namespace, req.URI.Resource)
	return err
}

func (cfgmgr *ConfigManager) MarkErrorURI(uri types.SubscriptionURI, subErr string) error {
	stmt := "UPDATE requests SET lastError = ?, errorTimestamp = ? WHERE namespace = ? AND resource = ?"
	_, err := cfgmgr.db.Exec(stmt, subErr, time.Now(), uri.Namespace, uri.Resource)
	return err
}

func (cfgmgr *ConfigManager) ClearErrorURI(uri types.SubscriptionURI) error {
	stmt := "UPDATE requests SET lastError = '', errorTimestamp = '' WHERE namespace = ? AND resource = ?"
	_, err := cfgmgr.db.Exec(stmt, uri.Namespace, uri.Resource)
	return err
}

func (cfgmgr *ConfigManager) List(filter *RequestFilter) ([]ArchiveRequest, error) {
	stmt := "SELECT schema, plugin, namespace, resource, inserted, coalesce(lastError, ''), errorTimestamp, enabled, id FROM requests"

	if filter == nil {
		stmt += ";"
	} else {
		var filters []string
		if filter.Schema != nil {
			filters = append(filters, fmt.Sprintf("schema='%s' ", *filter.Schema))
		}
		if filter.Plugin != nil {
			filters = append(filters, fmt.Sprintf("plugin='%s' ", *filter.Plugin))
		}
		if filter.Namespace != nil {
			filters = append(filters, fmt.Sprintf("namespace='%s' ", *filter.Namespace))
		}
		if filter.Resource != nil {
			filters = append(filters, fmt.Sprintf("resource='%s' ", *filter.Resource))
		}
		if filter.Id != nil {
			filters = append(filters, fmt.Sprintf("id='%d' ", *filter.Id))
		}
		if filter.HasError != nil {
			if *filter.HasError {
				filters = append(filters, fmt.Sprint("lastError!='' "))
			} else {
				filters = append(filters, fmt.Sprint("lastError='' "))
			}
		}
		if filter.Enabled != nil {
			if *filter.Enabled {
				filters = append(filters, fmt.Sprint("enabled=1 "))
			} else {
				filters = append(filters, fmt.Sprint("enabled=0 "))
			}
		}
		stmt = fmt.Sprintf("%s WHERE %s;", stmt, strings.Join(filters, " AND "))
	}

	var results []ArchiveRequest
	rows, err := cfgmgr.db.Query(stmt)
	if err != nil {
		return results, err
	}
	defer rows.Close()
	for rows.Next() {
		row := &ArchiveRequest{
			URI: types.SubscriptionURI{},
		}
		var et interface{}
		if err := rows.Scan(&row.Schema, &row.Plugin, &row.URI.Namespace, &row.URI.Resource, &row.Inserted, &row.LastError, &et, &row.Enabled, &row.Id); err != nil {
			return results, err
		}
		if et != nil {
			row.ErrorTimestamp = et.(time.Time)
		} else {
			row.ErrorTimestamp = time.Unix(0, 0)
		}

		results = append(results, *row)
	}
	return results, nil
}
