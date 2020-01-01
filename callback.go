package sorm

type Callback struct {
	logger     logger
	creates    []*func(scope *scope)
	updates    []*func(scope *scope)
	deletes    []*func(scope *scope)
	queries    []*func(scope *scope)
	rowQueries []*func(scope *scope)
	processors []*CallbackProcessor
}

type CallbackProcessor struct {
	logger    logger
	name      string
	before    string
	after     string
	replace   bool
	remove    bool
	kind      string
	processor *func(scope *scope)
	parent    *Callback
}
