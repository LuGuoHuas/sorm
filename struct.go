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
