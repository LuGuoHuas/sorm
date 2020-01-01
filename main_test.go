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
		initTable(db)
	}
}

func initTable(db *DB) {
	if result, err := db.GetRawDB().Exec(
		`
			CREATE TABLE IF NOT EXISTS table1 (
				field_1 varchar(32) NOT NULL,
				field_2 integer NOT NULL,
				field_3 bool NOT NULL ,
				CONSTRAINT table1_pk PRIMARY KEY (field_1)
			);`); err != nil {
		panic(err)
	} else {
		fmt.Println(result)
	}
}

type Table1 struct {
	Field1 string `json:"field_1" sorm:"column:field_1"`
	Field2 int    `json:"field_2" sorm:"column:field_2"`
	Field3 bool   `json:"field_3" sorm:"column:field_3"`
	Model
}

func (t Table1) getTableName() string {
	return "table1"
}

func TestMake(t *testing.T) {
	var table = Table1{
		Field1: "1234",
		Field2: 1234,
		Field3: true,
	}

	type args struct {
		model sorm
	}
	var tests = []struct {
		name       string
		args       args
		wantObject interface{}
	}{
		{
			"normal",
			args{
				model: &table,
			},
			&table,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotObject := Make(tt.args.model).(*Table1); !reflect.DeepEqual(gotObject, tt.wantObject) {
				t.Errorf("Make() = %v, want %v", tt.args.model, tt.wantObject)
			} else if gotObject.Field1 != table.Field1 ||
				gotObject.Field2 != table.Field2 ||
				gotObject.Field3 != table.Field3 ||
				gotObject.Model.Object == nil {
				t.Errorf("Make() = %v, want %v", tt.args.model, tt.wantObject)
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
	fmt.Println(filepath.Join(os.TempDir(), "sorm.db"))
	var source = os.Getenv("SOURCE")
	switch os.Getenv("DRIVER") {
	case "mysql":
		if len(source) < 1 {
			source = "sorm:T4SR46t4Iyev1qhTHQ8u5yuD@tcp(localhost:50001)/sorm?charset=utf8&parseTime=True"
		}
		return Open("mysql", source)
	case "postgres":
		if len(source) < 1 {
			source = "user=sorm password=T4SR46t4Iyev1qhTHQ8u5yuD db=sorm port=50002 sslmode=disable"
		}
		return Open("postgres", source)
	case "mssql":
		if len(source) < 1 {
			source = "sqlserver://sorm:T4SR46t4Iyev1qhTHQ8u5yuD@localhost:50003?database=sorm"
		}
		return Open("mssql", source)
	default:
		return Open("sqlite3", filepath.Join(os.TempDir(), "sorm.db"))
	}
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
			wantErr: false,
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

func TestDB_Create(t *testing.T) {
	var d *DB
	if i, ok := TestData.Load(TestConnectKey); !ok || i == nil {
		panic("failed to load database connection")
	} else {
		d = i.(*DB)
	}
	type args struct {
		value sorm
	}
	tests := []struct {
		name   string
		fields *DB
		args   args
		want   *DB
	}{
		{
			"normal",
			d,
			args{
				value: nil,
			},
			d,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.value = Make(&Table1{
				//Field1: "test",
				Field2: 1024,
				Field3: true,
			}).(*Table1)
			if got := d.Create(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_Find(t *testing.T) {
	var d *DB
	if i, ok := TestData.Load(TestConnectKey); !ok || i == nil {
		panic("failed to load database connection")
	} else {
		d = i.(*DB)
	}
	type args struct {
		obj sorm
	}
	tests := []struct {
		name   string
		fields *DB
		args   args
	}{
		{
			name:   "normal",
			fields: d,
			args:   args{obj: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.obj = Make(&Table1{
				Field1: "test",
				Field2: 1024,
				Field3: true,
			}).(*Table1)
			d.Table(tt.args.obj).Find(tt.args.obj)
			fmt.Println(tt.args.obj)
		})
	}
}
