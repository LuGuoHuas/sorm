package sorm

import (
	"reflect"
	"sort"
	"strings"
	"unsafe"
)

type sorm interface {
	instant(object interface{})
	getTag(field interface{}) string
	setTableName()
	getTableName() string
	getFieldIndex() []int
	getField(i int) StructField
	getFieldValue(i int) interface{}
	getValue() []interface{}
}

type Model struct {
	Object     interface{}
	fields     map[int]StructField
	fieldIndex []int
	tableName  string
}

func Make(model sorm) interface{} {
	model.instant(model)
	model.setTableName()
	return model
}

func (m *Model) instant(object interface{}) {
	m.Object = object
	if m.fields == nil {
		m.fields = make(map[int]StructField)
	}
	var e = reflect.TypeOf(object).Elem()
	m.fieldIndex = make([]int, 0, 10)
	for i := 0; i < e.NumField(); i++ {
		if e.Field(i).Name == "Model" {
			continue
		}
		m.fields[int(e.Field(i).Offset)] = StructField{
			Name:    e.Field(i).Name,
			Tag:     analyseTag(e.Field(i).Tag.Get("sorm")),
			Pointer: unsafe.Pointer(reflect.ValueOf(object).Pointer() + e.Field(i).Offset),
			Type:    e.Field(i).Type.Kind(),
		}
		m.fieldIndex = append(m.fieldIndex, int(e.Field(i).Offset))
	}
	sort.Ints(m.fieldIndex)
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

func (m *Model) setTableName() {
	m.tableName = reflect.TypeOf(m.Object).Elem().Name()
}

func (m *Model) getTableName() string {
	return m.tableName
}

func (m *Model) getFieldIndex() []int {
	return m.fieldIndex
}

func (m *Model) getField(i int) StructField {
	return m.fields[i]
}

func (m *Model) getFieldValue(i int) interface{} {
	switch m.fields[i].Type {
	case reflect.String:
		return *(*string)(m.fields[i].Pointer)
	case reflect.Int:
		return *(*int)(m.fields[i].Pointer)
	case reflect.Bool:
		return *(*bool)(m.fields[i].Pointer)
	}
	return *(*interface{})(m.fields[i].Pointer)
}

func (m *Model) getValue() []interface{} {
	var result = make([]interface{}, 0, 10)
	for _, i := range m.fieldIndex {
		result = append(result, m.getFieldValue(i))
	}
	return result
}

func analyseTag(tag string) (result map[string]string) {
	result = make(map[string]string)
	for _, t := range strings.Split(tag, ";") {
		s := strings.Split(t, ":")
		if len(s) < 2 {
			continue
		}
		result[s[0]] = strings.Join(s[1:], ":")
	}
	return result
}
