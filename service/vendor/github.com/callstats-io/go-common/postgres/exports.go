package postgres

import (
	"github.com/go-pg/pg"
)

// Exported constants from pg
var (
	ErrNoRows    = pg.ErrNoRows
	ErrMultiRows = pg.ErrMultiRows
)

// DB aliased from mgo.DB
type DB struct {
	*pg.DB
}

// Status return an error if the status query failed, nil otherwise
func (d *DB) Status() error {
	if _, err := d.Exec("Select 1 as ONE"); err != nil {
		return err
	}
	return nil
}

// struct alias helpers
func asAliasedDB(db *pg.DB) *DB {
	return &DB{
		DB: db,
	}
}
