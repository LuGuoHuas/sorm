package sorm

import (
	"database/sql"
)

type model interface {
	instant(obj interface{})
}

type Model struct {
	Object interface{}
}

func Make(mod model) (object interface{}) {
	mod.instant(mod)
	return mod
}

type DB struct {
	db *sql.DB
}

func Open(driver, source string) (db *DB, err error) {
	if len(driver) == 0 || len(source) == 0 {
		return nil, InvalidDatabaseSource
	}

	var sqlDB *sql.DB
	if sqlDB, err = sql.Open(driver, source); err != nil {
		return nil, err
	} else if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	db = &DB{
		db: sqlDB,
	}
	return db, nil
}

func (d *DB) GetRawDB() (db *sql.DB) {
	return d.db
}

func (d *DB) Close() (err error) {
	return d.db.Close()
}
