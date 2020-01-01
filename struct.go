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
}
