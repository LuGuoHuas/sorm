package sorm

import "fmt"

const BeginTransaction int = -1400
const BeforeCreate int = -1300
const SaveBeforeAssociations int = -1200
const UpdateTimeStamp int = -1100
const create int = 1000
const ForceReloadAfterCreate int = 1100
const SaveAfterAssociations int = 1200
const AfterCreate int = 1300
const CommitOrRollbackTransaction int = 1400

const createCallbackKind = "create"
const updateCallbackKind = "update"
const deleteCallbackKind = "delete"
const queryCallbackKind = "query"
const rowQueryCallbackKind = "row_query"

var DefaultCallback = new(Callback)

type Callback struct {
	logger     logger
	creates    map[int]*CallbackProcessor
	updates    map[int]*CallbackProcessor
	deletes    map[int]*CallbackProcessor
	queries    map[int]*CallbackProcessor
	rowQueries map[int]*CallbackProcessor
}

type CallbackProcessor struct {
	logger   logger
	kind     string
	priority int
	callback *func(Scope *Scope)
	parent   *Callback
}

func (c *Callback) Create() (processor *CallbackProcessor) {
	processor = &CallbackProcessor{logger: c.logger, kind: createCallbackKind, parent: c}
	return processor
}

// Update could be used to register callbacks for updating object, refer `Create` for usage
func (c *Callback) Update() (processor *CallbackProcessor) {
	processor = &CallbackProcessor{logger: c.logger, kind: updateCallbackKind, parent: c}
	return processor
}

// Delete could be used to register callbacks for deleting object, refer `Create` for usage
func (c *Callback) Delete() (processor *CallbackProcessor) {
	processor = &CallbackProcessor{logger: c.logger, kind: deleteCallbackKind, parent: c}
	return processor
}

// Query could be used to register callbacks for querying objects with query methods like `Find`, `First`, `Related`, `Association`...
// Refer `Create` for usage
func (c *Callback) Query() (processor *CallbackProcessor) {
	processor = &CallbackProcessor{logger: c.logger, kind: queryCallbackKind, parent: c}
	return processor
}

// RowQuery could be used to register callbacks for querying objects with `Row`, `Rows`, refer `Create` for usage
func (c *Callback) RowQuery() (processor *CallbackProcessor) {
	processor = &CallbackProcessor{logger: c.logger, kind: rowQueryCallbackKind, parent: c}
	return processor
}

// Register a new callback, refer `Callbacks.Create`
func (processor *CallbackProcessor) Register(priority int, callback func(scope *Scope)) {
	processor.logger.Print("info", fmt.Sprintf("[info] registering callback `%v` from %v", priority, fileWithLineNum()))
	processor.callback = &callback
	processor.priority = priority

	switch processor.kind {
	case createCallbackKind:
		processor.parent.creates[priority] = processor
	case updateCallbackKind:
		processor.parent.updates[priority] = processor
	case deleteCallbackKind:
		processor.parent.deletes[priority] = processor
	case queryCallbackKind:
		processor.parent.queries[priority] = processor
	case rowQueryCallbackKind:
		processor.parent.rowQueries[priority] = processor
	}
	return
}

// Remove a registered callback
//     db.Callback().Create().Remove(1)
func (processor *CallbackProcessor) Remove(priority int) {
	processor.logger.Print("info", fmt.Sprintf("[info] removing callback `%v` from %v", priority, fileWithLineNum()))

	switch processor.kind {
	case createCallbackKind:
		delete(processor.parent.creates, priority)
	case updateCallbackKind:
		delete(processor.parent.updates, priority)
	case deleteCallbackKind:
		delete(processor.parent.deletes, priority)
	case queryCallbackKind:
		delete(processor.parent.queries, priority)
	case rowQueryCallbackKind:
		delete(processor.parent.rowQueries, priority)
	}
	return
}

// Get registered callback
//    db.Callback().Create().Get(0)
func (processor *CallbackProcessor) Get(priority int) (callback func(scope *Scope)) {
	switch processor.kind {
	case createCallbackKind:
		callback = *processor.parent.creates[priority].callback
	case updateCallbackKind:
		callback = *processor.parent.updates[priority].callback
	case deleteCallbackKind:
		callback = *processor.parent.deletes[priority].callback
	case queryCallbackKind:
		callback = *processor.parent.queries[priority].callback
	case rowQueryCallbackKind:
		callback = *processor.parent.rowQueries[priority].callback
	}
	return
}
