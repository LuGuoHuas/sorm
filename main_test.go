package sorm

import (
	"reflect"
	"testing"
)

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
