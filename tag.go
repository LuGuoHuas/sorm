package sorm

func (m *Model) Instant(obj interface{}) {
	(*m).Object = &obj
}
