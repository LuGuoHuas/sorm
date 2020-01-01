package sorm

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

type Scope struct {
	table  string
	object sorm
}

type DB struct {
	sync.RWMutex
	db    *sql.DB
	Error error
	scope *Scope
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

func (d *DB) Create(obj sorm) *DB {
	var s bytes.Buffer
	s.WriteString("INSERT INTO ")
	s.WriteString(obj.getTableName())
	s.WriteString(" (")
	for _, i := range obj.getFieldIndex() {
		s.WriteString(obj.getField(i).Tag["column"])
		s.WriteString(",")
		if obj.getField(i).Type == reflect.String {
			fmt.Println(*(*string)(obj.getField(i).Pointer))
		}
	}
	s.Truncate(s.Len() - 1)
	s.WriteString(") VALUES (")
	for range obj.getFieldIndex() {
		s.WriteString("?,")
	}
	s.Truncate(s.Len() - 1)
	s.WriteString(")")
	fmt.Println(s.String())
	if _, err := d.db.Exec(s.String(), obj.getValue()...); err != nil {
		panic(err)
	}
	return d
}

func (d *DB) Table(obj sorm) *DB {
	d.scope = &Scope{table: obj.getTableName(), object: obj}
	return d
}

func (d *DB) Find(obj sorm) {
	var s = bytes.NewBufferString("SELECT * FROM ")
	s.WriteString(obj.getTableName())
	s.WriteString(" WHERE ")
	for _, i := range obj.getFieldIndex() {
		s.WriteString(obj.getField(i).Tag["column"])
		s.WriteString("=? AND ")
	}
	s.Truncate(s.Len() - 4)
	fmt.Println(s.String())
	var rows, err = d.db.Query(s.String(), obj.getValue()...)
	d.Error = err
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		if err = rows.Scan(obj.getValue()...); err != nil {
			d.Error = err
		}
	}
}

func (d *DB) Update(field, value interface{}) {
	var err error
	var f = d.scope.object.getFieldByPointer(unsafe.Pointer(reflect.ValueOf(field).Pointer()))
	var s = bytes.NewBufferString("UPDATE ")
	s.WriteString(d.scope.table)
	s.WriteString(" SET ")
	s.WriteString(f.Tag["column"])
	s.WriteString("=? WHERE ")
	s.WriteString(d.scope.object.getField(0).Tag["column"])
	s.WriteString("=?")

	fmt.Println(s.String())
	if _, err = d.db.Exec(s.String(), value, d.scope.object.getField(0).get(d.scope.object.getField(0).Pointer)); err != nil {
		d.Error = err
	}
}

func (d *DB) Delete() {
	var err error
	var s = bytes.NewBufferString("DELETE FROM ")
	s.WriteString(d.scope.table)
	s.WriteString(" WHERE ")
	s.WriteString(d.scope.object.getField(0).Tag["column"])
	s.WriteString("=?")

	fmt.Println(s.String())
	if _, err = d.db.Exec(s.String(), d.scope.object.getField(0).get(d.scope.object.getField(0).Pointer)); err != nil {
		d.Error = err
	}
}
