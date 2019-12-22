package sorm

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

var TestData sync.Map

const TestConnectKey = "db"

func init() {
	if db, err := OpenTestConnection(); err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to test database, but got err=%+v", err))
	} else {
		TestData.Store(TestConnectKey, db)
	}
}

func TestMake(t *testing.T) {
	type Table struct {
		Field1 string `json:"field_1"`
		Field2 int    `json:"field_2"`
		Field3 bool   `json:"field_3"`
		Model
	}

	var st1 = Table{
		Field1: "1234",
		Field2: 1234,
		Field3: true,
	}

	type args struct {
		model model
	}
	var tests = []struct {
		name       string
		args       args
		wantObject interface{}
	}{
		{
			"normal",
			args{
				model: &st1,
			},
			&st1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotObject := Make(tt.args.model); !reflect.DeepEqual(gotObject, tt.wantObject) {
				t.Errorf("Make() = %v, want %v", gotObject, tt.wantObject)
			} else if gotObject.(*Table).Field1 != st1.Field1 ||
				gotObject.(*Table).Field2 != st1.Field2 ||
				gotObject.(*Table).Field3 != st1.Field3 ||
				gotObject.(*Table).Model.Object == nil {
				t.Errorf("Make() = %v, want %v", gotObject, tt.wantObject)
			}
		})
	}
}

func TestOpen(t *testing.T) {
	type args struct {
		driver string
		source string
	}
	tests := []struct {
		name    string
		args    args
		wantDb  bool
		wantErr bool
	}{
		{
			name: "normal: sqlite3",
			args: args{
				driver: "sqlite3",
				source: filepath.Join(os.TempDir(), "sorm.db"),
			},
			wantDb:  true,
			wantErr: false,
		},
		{
			name: "exception: nil driver",
			args: args{
				source: filepath.Join(os.TempDir(), "sorm.db"),
			},
			wantDb:  false,
			wantErr: true,
		},
		{
			name: "exception: error source",
			args: args{
				driver: "postgres",
				source: "1234",
			},
			wantDb:  false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDb, err := Open(tt.args.driver, tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotDb != nil) != tt.wantDb {
				t.Errorf("Open() gotDb = %v, want %v", gotDb, tt.wantDb)
			}
		})
	}
}

func OpenTestConnection() (db *DB, err error) {
	return Open("sqlite3", filepath.Join(os.TempDir(), "sorm.db"))
}

func TestDB_GetRawDB(t *testing.T) {
	var d *DB
	if i, ok := TestData.Load(TestConnectKey); !ok || i == nil {
		panic("failed to load database connection")
	} else {
		d = i.(*DB)
	}

	tests := []struct {
		name   string
		fields *DB
		wantDb *sql.DB
	}{
		{
			name:   "normal",
			fields: d,
			wantDb: d.db,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDb := tt.fields.GetRawDB(); !reflect.DeepEqual(gotDb, tt.wantDb) {
				t.Errorf("GetRawDB() = %v, want %v", gotDb, tt.wantDb)
			}
		})
	}
}

func TestDB_Close(t *testing.T) {
	var d *DB
	if i, ok := TestData.Load(TestConnectKey); !ok || i == nil {
		panic("failed to load database connection")
	} else {
		d = i.(*DB)
	}
	tests := []struct {
		name    string
		fields  *DB
		wantErr bool
	}{
		{
			name:    "normal",
			fields:  d,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
