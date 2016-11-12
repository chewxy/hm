package hm

import "fmt"

// FunctionType is a type constructor that builds function types.
type FunctionType struct {
	ts [2]Type // from → to
}

func NewFnType(params ...Type) *FunctionType {
	if len(params) < 2 {
		panic(fmt.Sprintf("Needs more than 2 params to make a function. Got %v", params))
	}

	t := borrowFnType()
	t.ts[0] = params[0]
	if len(params) == 2 {
		t.ts[1] = params[1]
		return t
	}

	t.ts[1] = NewFnType(params[1:]...)
	return t
}

/* Type interface fulfilment */

func (t *FunctionType) Name() string { return "→" }

func (t *FunctionType) Contains(tv TypeVariable) bool {
	for _, ty := range t.ts {
		if ty.Contains(tv) {
			return true
		}
	}
	return false
}

func (t *FunctionType) Eq(other Type) bool {
	oft, ok := other.(*FunctionType)
	if !ok {
		return false
	}

	for i, tt := range t.ts {
		if !tt.Eq(oft.ts[i]) {
			return false
		}
	}
	return true
}

func (t *FunctionType) Format(state fmt.State, c rune) {
	fmt.Fprintf(state, "%s → %s", t.ts[0], t.ts[1])
}

func (t *FunctionType) String() string { return fmt.Sprintf("%v", t) }

/* TypeOp interface fulfilment */

func (t *FunctionType) Types() Types { return Types(t.ts[:]) }

func (t *FunctionType) Replace(tv TypeVariable, with Type) TypeOp {
	switch tt := t.ts[0].(type) {
	case TypeVariable:
		t.ts[0] = with
	case TypeConst:
		// do nothing
	case TypeOp:
		tt = tt.Replace(tv, with)
		t.ts[0] = tt
	default:
		panic("WTF?")
	}

	switch tt := t.ts[1].(type) {
	case TypeVariable:
		t.ts[1] = with
	case TypeConst:
		// do nothing
	case TypeOp:
		t.ts[1] = tt.Replace(tv, with)
	default:
		panic("WTF?")
	}

	return t
}

func (t *FunctionType) Clone() TypeOp {
	retVal := new(FunctionType)

	switch tt := t.ts[0].(type) {
	case TypeVariable:
		retVal.ts[0] = tt
	case TypeConst:
		retVal.ts[0] = tt
	case TypeOp:
		retVal.ts[0] = tt.Clone()
	default:
		panic("What")
	}

	switch tt := t.ts[1].(type) {
	case TypeVariable:
		retVal.ts[1] = tt
	case TypeConst:
		retVal.ts[1] = tt
	case TypeOp:
		retVal.ts[1] = tt.Clone()
	default:
		panic("What")
	}

	return retVal
}

/* Useful methods */

// TypesRec is like Types(), but recursively gets more Types. This is often used to get the return type of a function
func (t *FunctionType) TypesRec() (retVal Types) {
	for _, tt := range t.ts {
		if fn, ok := tt.(*FunctionType); ok {
			retVal = append(retVal, fn.TypesRec()...)
			continue
		}
		retVal = append(retVal, tt)
	}
	return
}

// ReturnType is the specialization of TypesRec(), specialized for finding return types
func (t *FunctionType) ReturnType() Type {
	if fn, ok := t.ts[1].(*FunctionType); ok {
		return fn.ReturnType()
	}
	return t.ts[1]
}
