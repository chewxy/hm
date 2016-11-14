package hm

import "fmt"

type TypeVariable struct {
	name     string
	instance Type

	constraints *TypeClassSet
}

type TypeVarConsOpt func(tv *TypeVariable)

func WithInstance(t Type) TypeVarConsOpt {
	f := func(tv *TypeVariable) {
		tv.instance = t
	}
	return f
}

func WithConstraints(cs *TypeClassSet) TypeVarConsOpt {
	f := func(tv *TypeVariable) {
		tv.constraints = cs
	}
	return f
}

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

// func (t TypeVariable) In(t0 Type) bool {
// 	pruned := Prune(t0)

// 	if ptv, ok := pruned.(TypeVariable); ok && t.Eq(ptv) {
// 		return true
// 	}

// 	if op, ok := pruned.(TypeOp); ok {
// 		ts := op.Types()
// 		if len(ts) == 1 {

// 		}
// 		return t.InTypes(ts)
// 	}
// 	return false
// }

// func (t TypeVariable) InTypes(ts Types) bool {
// 	for _, typ := range ts {
// 		if t.In(typ) {
// 			return true
// 		}
// 	}
// 	return false
// }

func (t TypeVariable) IsEmpty() bool {
	return t.name == "" && t.instance == nil && (t.constraints == nil || len(t.constraints.s) == 0)
}
