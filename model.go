package sorm

import (
	"reflect"
)

type sorm interface {
	instant(object interface{})
	getTag(field interface{}) string
}

type Model struct {
	Object interface{}
}

func Make(model sorm) {
	model.instant(model)
}

func (m *Model) instant(object interface{}) {
	m.Object = object
}

func (m *Model) getTag(field interface{}) (tag string) {
	var te = reflect.TypeOf(m.Object).Elem()
	for i := 0; i < te.NumField(); i++ {
		if te.Field(i).Offset == reflect.ValueOf(field).Pointer()-reflect.ValueOf(m.Object).Pointer() {
			return te.Field(i).Name
		}
	}
	return ""
}
