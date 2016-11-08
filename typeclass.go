package hm

// TypeClass is like an interface{} in Go.
type TypeClass interface {
	AddInstance(Type)
}

type SimpleTypeClass struct {
	instances Types
}

func (tc *SimpleTypeClass) AddInstance(t Type) {
	tc.instances = tc.instances.Add(t)
}
