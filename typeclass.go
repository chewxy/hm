package hm

// TypeClass is like an interface{} in Go.
type TypeClass interface {
	Name() string
	AddInstance(Type)
}

type SimpleTypeClass struct {
	name      string
	instances Types
}

func NewSimpleTypeClass(name string) *SimpleTypeClass {
	return &SimpleTypeClass{
		name: name,
	}
}

func (tc *SimpleTypeClass) Name() string { return tc.name }

func (tc *SimpleTypeClass) AddInstance(t Type) {
	tc.instances = tc.instances.Add(t)
}

func (tc *SimpleTypeClass) String() string { return tc.name }
