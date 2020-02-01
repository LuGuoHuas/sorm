package sorm

import "bytes"

type Scope struct {
	table     string
	object    sorm
	parent    *Scope
	subScope  *Scope
	SQL       bytes.Buffer
	SQLValues []interface{}
	Values    interface{}
	Search    *search
}

func (s *Scope) generate() *Scope {
	var newScope = &Scope{
		table:  s.table,
		object: s.object,
		parent: s,
	}
	s.subScope = newScope
	return newScope
}
