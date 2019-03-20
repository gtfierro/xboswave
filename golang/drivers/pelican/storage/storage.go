package storage

import (
	"database/sql"
	"encoding/json"
	"time"

	sqlTypes "github.com/jmoiron/sqlx/types"
	"github.com/john-b-yang/xboswave/golang/drivers/pelican/types"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const dbURL = "postgres://xbosreadonly@corbusier.cs.berkeley.edu:26257/xbos?sslmode=verify-full&sslrootcert=ca.crt&sslcert=client.xbosreadonly.crt&sslkey=client.xbosreadonly.key"
const dbKey = "pelicans"

func WritePelicans(pelicans []*types.Pelican, sitename string) error {
	bytes, err := json.Marshal(pelicans)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize Pelicans")
	}
	var json sqlTypes.JSONText
	if err := json.Scan(bytes); err != nil {
		return errors.Wrap(err, "Failed to convert to JSON Text for SQL")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return errors.Wrap(err, "Failed to connect to database")
	}
	_, err = db.Exec("INSERT INTO SETTINGS (sitename, inserted, key, object) VALUES ($1, $2, $3, $4)",
		sitename, time.Now(), dbKey, json)
	if err != nil {
		return errors.Wrap(err, "Failed to insert into database")
	}

	return nil
}

func ReadPelicans(username, password, sitename string) ([]*types.Pelican, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to database")
	}

	var ob sqlTypes.JSONText
	err = db.QueryRow("SELECT object FROM settings where sitename = $1 and key = $2 order by inserted desc limit 1;",
		sitename, dbKey).Scan(&ob)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read from database")
	}

	var pelicans []*types.Pelican
	if err := ob.Unmarshal(&pelicans); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize Pelican info")
	}

	// To properly regenerate internal fields
	for i := 0; i < len(pelicans); i++ {
		newPelican, err := types.NewPelican(&types.NewPelicanParams{
			Username:      username,
			Password:      password,
			Sitename:      sitename,
			Name:          pelicans[i].Name,
			HeatingStages: pelicans[i].HeatingStages,
			CoolingStages: pelicans[i].CoolingStages,
			Timezone:      pelicans[i].TimezoneName,
		})
		if err != nil {
			return nil, errors.Wrap(err, "Failed to instantiate Pelican")
		}
		pelicans[i] = newPelican
	}
	return pelicans, nil
}
