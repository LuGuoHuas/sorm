package sorm

import (
	"math"
	"testing"
)

func TestModel_GetTag(t *testing.T) {
	type args struct {
		Field1 string      `json:"field_1"`
		Field2 int         `json:"field_2"`
		Field3 bool        `json:"field_3"`
		Field4 float64     `json:"field_4"`
		Field5 interface{} `json:"field_5"`
		Model
	}
	tests := []struct {
		name    string
		args    *args
		wantTag string
	}{
		{
			name: "normal",
			args: &args{
				Field1: "string field",
				Field2: 1024,
				Field3: true,
				Field4: math.Pi,
				Field5: "string field",
			},
			wantTag: "Field5",
		},
		{
			name: "normal",
			args: &args{
				Field1: "string field",
				Field2: 1024,
				Field3: true,
				Field4: math.Pi,
				Field5: 1024,
			},
			wantTag: "Field5",
		},
		{
			name: "normal",
			args: &args{
				Field1: "string field",
				Field2: 1024,
				Field3: true,
				Field4: math.Pi,
				Field5: true,
			},
			wantTag: "Field5",
		},
		{
			name: "normal",
			args: &args{
				Field1: "string field",
				Field2: 1024,
				Field3: true,
				Field4: math.Pi,
				Field5: make(map[string]interface{}),
			},
			wantTag: "Field5",
		},
		{
			name: "normal",
			args: &args{
				Field1: "string field",
				Field2: 1024,
				Field3: true,
				Field4: math.Pi,
				Field5: struct {
					field string
				}{field: "field"},
			},
			wantTag: "Field5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Make(tt.args)
			if name := tt.args.getTag(&tt.args.Field5); name != tt.wantTag {
				t.Errorf("getTag() name = %v, want %v", name, tt.wantTag)
			}
		})
	}
}
