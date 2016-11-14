package hm

import "fmt"

// TypeVariable represents a type variable. It allows polymorphic types
type TypeVariable struct {
	name     string
	instance Type

	constraints *TypeClassSet
}

type TypeVarConsOpt func(tv *TypeVariable)

// WithInstance is an option that creates a TypeVariable with an instance already
func WithInstance(t Type) TypeVarConsOpt {
	f := func(tv *TypeVariable) {
		tv.instance = t
	}
	return f
}

// WithConstraints is an option that creates a TypeVariable with a type class constraints
func WithConstraints(cs ...TypeClass) TypeVarConsOpt {
	f := func(tv *TypeVariable) {
		constraints := NewTypeClassSet(cs...)
		tv.constraints = constraints
	}
	return f
}

// NewTypeVar creates a new TypeVariable
func NewTypeVar(name string, opts ...TypeVarConsOpt) TypeVariable {
	retVal := TypeVariable{
		name: name,
	}

	for _, opt := range opts {
		opt(&retVal)
	}

	return retVal
}

func (t TypeVariable) Name() string { return t.name }

func (t TypeVariable) Contains(tv TypeVariable) bool {
	if t.Eq(tv) {
		return true
	}

	return false
}

func (t TypeVariable) Eq(other Type) bool {
	var tv TypeVariable
	var ok bool
	if tv, ok = other.(TypeVariable); !ok {
		return false
	}

	if t.name != tv.name {
		return false
	}

	switch {
	case t.instance != nil && tv.instance != nil:
		if !t.instance.Eq(tv.instance) {
			panic(fmt.Sprintf(tvinstance, t, tv, t.name))
		}
		return true
	case t.instance == nil && tv.instance == nil:
		return true
	default:
		return false
	}

}

func (t TypeVariable) Format(state fmt.State, c rune) {
	if t.instance == nil {
		name := "''"
		if t.name != "" {
			name = t.name
		}

		if state.Flag('#') {
			fmt.Fprintf(state, "%s:%#v", name, t.instance)
		} else {
			fmt.Fprintf(state, "%v", name)
		}

	} else {
		if state.Flag('#') {
			fmt.Fprintf(state, "%s:%#v", t.name, t.instance)
		} else {
			fmt.Fprintf(state, "%v", t.instance)
		}
	}
}

func (t TypeVariable) String() string {
	if t.instance != nil {
		return t.instance.String()
	}
	if t.name == "" {
		return "''"
	}
	return t.name
}

// IsEmpty returns true if it's a dummy/empty type variable - defined as a TypeVariable with no name, and no constraints nor instances
func (t TypeVariable) IsEmpty() bool {
	return t.name == "" && t.instance == nil && (t.constraints == nil || len(t.constraints.s) == 0)
}

// Instance returns the instance that defines the TypeVariable
func (t TypeVariable) Instance() Type { return t.instance }
