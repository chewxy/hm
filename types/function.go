package hmtypes

import (
	"fmt"

	"github.com/chewxy/hm"
)

// FunctionType is a type constructor that builds function types.
type FunctionType Pair

// NewFnType creates a new FunctionType. Functions are by default right associative. This:
//		NewFnType(a, a, a)
// is short hand for this:
// 		NewFnType(a, NewFnType(a, a))
func NewFnType(ts ...hm.Type) *FunctionType {
	if len(ts) < 2 {
		panic("Expected at least 2 input types")
	}

	retVal := borrowFnType()
	retVal.A = ts[0]

	if len(ts) > 2 {
		retVal.B = NewFnType(ts[1:]...)
	} else {
		retVal.B = ts[1]
	}
	return retVal
}

func (t *FunctionType) Name() string                       { return "→" }
func (t *FunctionType) Apply(sub hm.Subs) hm.Substitutable { ((*Pair)(t)).Apply(sub); return t }
func (t *FunctionType) FreeTypeVar() hm.TypeVarSet         { return ((*Pair)(t)).FreeTypeVar() }
func (t *FunctionType) Format(s fmt.State, c rune)         { fmt.Fprintf(s, "%v → %v", t.A, t.B) }
func (t *FunctionType) String() string                     { return fmt.Sprintf("%v", t) }
func (t *FunctionType) Normalize(k, v hm.TypeVarSet) (hm.Type, error) {
	var a, b hm.Type
	var err error
	if a, err = t.A.Normalize(k, v); err != nil {
		return nil, err
	}

	if b, err = t.B.Normalize(k, v); err != nil {
		return nil, err
	}

	return NewFnType(a, b), nil
}
func (t *FunctionType) Types() hm.Types {
	retVal := hm.BorrowTypes(2)
	retVal[0] = t.A
	retVal[1] = t.B
	return retVal
}

func (t *FunctionType) Eq(other hm.Type) bool {
	if ot, ok := other.(*FunctionType); ok {
		return ot.A.Eq(t.A) && ot.B.Eq(t.B)
	}
	return false
}

// Other methods (accessors mainly)

// Arg returns the type of the function argument
func (t *FunctionType) Arg() hm.Type { return t.A }

// Ret returns the return type of a function. If recursive is true, it will get the final return type
func (t *FunctionType) Ret(recursive bool) hm.Type {
	if !recursive {
		return t.B
	}

	if fnt, ok := t.B.(*FunctionType); ok {
		return fnt.Ret(recursive)
	}

	return t.B
}

// FlatTypes returns the types in FunctionTypes as a flat slice of types. This allows for easier iteration in some applications
func (t *FunctionType) FlatTypes() hm.Types {
	retVal := hm.BorrowTypes(8) // start with 8. Can always grow
	retVal = retVal[:0]

	if a, ok := t.A.(*FunctionType); ok {
		ft := a.FlatTypes()
		retVal = append(retVal, ft...)
		hm.ReturnTypes(ft)
	} else {
		retVal = append(retVal, t.A)
	}

	if b, ok := t.B.(*FunctionType); ok {
		ft := b.FlatTypes()
		retVal = append(retVal, ft...)
		hm.ReturnTypes(ft)
	} else {
		retVal = append(retVal, t.B)
	}
	return retVal
}

// Clone implements Cloner
func (t *FunctionType) Clone() interface{} {
	retVal := new(FunctionType)

	if ac, ok := t.A.(Cloner); ok {
		retVal.A = ac.Clone().(hm.Type)
	} else {
		retVal.A = t.A
	}

	if bc, ok := t.B.(Cloner); ok {
		retVal.B = bc.Clone().(hm.Type)
	} else {
		retVal.B = t.B
	}
	return retVal
}
