package hm

import "fmt"

type FunctionType struct {
	a, b Type
}

func NewFnType(ts ...Type) *FunctionType {
	if len(ts) < 2 {
		panic("Expected at least 2 input types")
	}

	retVal := BorrowFnType()
	retVal.a = ts[0]

	if len(ts) > 2 {
		retVal.b = NewFnType(ts[1:]...)
	} else {
		retVal.b = ts[1]
	}
	return retVal
}

func (t *FunctionType) Name() string { return "→" }
func (t *FunctionType) Apply(sub Subs) Substitutable {
	t.a = t.a.Apply(sub).(Type)
	t.b = t.b.Apply(sub).(Type)
	return t
}

func (t *FunctionType) FreeTypeVar() TypeVarSet    { return t.a.FreeTypeVar().Union(t.b.FreeTypeVar()) }
func (t *FunctionType) Format(s fmt.State, c rune) { fmt.Fprintf(s, "%v → %v", t.a, t.b) }
func (t *FunctionType) String() string             { return fmt.Sprintf("%v", t) }
func (t *FunctionType) Normalize(k, v TypeVarSet) (Type, error) {
	var a, b Type
	var err error
	if a, err = t.a.Normalize(k, v); err != nil {
		return nil, err
	}

	if b, err = t.b.Normalize(k, v); err != nil {
		return nil, err
	}

	return NewFnType(a, b), nil
}
func (t *FunctionType) Types() Types {
	retVal := BorrowTypes(2)
	retVal[0] = t.a
	retVal[1] = t.b
	return retVal
}

func (t *FunctionType) Eq(other Type) bool {
	if ot, ok := other.(*FunctionType); ok {
		return ot.a.Eq(t.a) && ot.b.Eq(t.b)
	}
	return false
}
