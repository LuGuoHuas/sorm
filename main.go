package sorm

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"sync"
)

type DB struct {
	sync.RWMutex
	db    *sql.DB
	Error error
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

func (d *DB) Create(value sorm) *DB {
	var s bytes.Buffer
	s.WriteString("INSERT INTO ")
	s.WriteString(value.getTableName())
	s.WriteString(" (")
	for _, i := range value.getFieldIndex() {
		s.WriteString(value.getField(i).Tag["column"])
		s.WriteString(",")
		if value.getField(i).Type == reflect.String {
			fmt.Println(*(*string)(value.getField(i).Pointer))
		}
	}
	s.Truncate(s.Len() - 1)
	s.WriteString(") VALUES (")
	for _, _ = range value.getFieldIndex() {
		s.WriteString("?,")
	}
	s.Truncate(s.Len() - 1)
	s.WriteString(")")
	fmt.Println(s.String())
	if _, err := d.db.Exec(s.String(), value.getValue()...); err != nil {
		panic(err)
	}
	return d
}
