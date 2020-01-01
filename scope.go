package sorm

import "bytes"

type scope struct {
	table     string
	object    sorm
	parent    *scope
	subScope  *scope
	SQL       bytes.Buffer
	SQLValues []interface{}
	Values    interface{}
	Search    *search
}

func (s *scope) generate() *scope {
	var newScope = &scope{
		table:  s.table,
		object: s.object,
		parent: s,
	}
	s.subScope = newScope
	return newScope
}
