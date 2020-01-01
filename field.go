package sorm

import (
	"reflect"
	"unsafe"
)

type StructField struct {
	Name    string            `json:"name"`
	Tag     map[string]string `json:"tag"`
	Pointer unsafe.Pointer
	Type    reflect.Kind
	get     func(pointer unsafe.Pointer) interface{}
}

func getString(pointer unsafe.Pointer) interface{} {
	return (*string)(pointer)
}

func getInt(pointer unsafe.Pointer) interface{} {
	return (*int)(pointer)
}

func getBool(pointer unsafe.Pointer) interface{} {
	return (*bool)(pointer)
}

func newField(In interface{}, pointer unsafe.Pointer, field reflect.StructField) *StructField {
	var newField = &StructField{
		Name:    field.Name,
		Tag:     analyseTag(field.Tag.Get("sorm")),
		Pointer: pointer,
		Type:    field.Type.Kind(),
	}

	switch v := In.(type) {
	case string:
		newField.get = getString
	case int:
		newField.get = getInt
	case bool:
		newField.get = getBool
	default:
		panic(v)
	}

	return newField
}

func getKey(field *StructField, keyMap map[string]*StructField) map[string]*StructField {
	for k, v := range field.Tag {
		switch k {
		case "PrimaryKey":
			if v == "true" {
				keyMap["PrimaryKey"] = field
			}
		}
	}
	return keyMap
}
