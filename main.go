package sorm

type model interface {
	Instant(obj interface{})
}

type Model struct {
	Object interface{}
}

func Make(mod model) (object interface{}) {
	mod.Instant(mod)
	return mod
}
