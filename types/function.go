package hmtypes

import (
	"fmt"

	"github.com/chewxy/hm"
)

// Function is a type constructor that builds function types.
type Function hm.Pair

// NewFunction creates a new FunctionType. Functions are by default right associative. This:
//		NewFunction(a, a, a)
// is short hand for this:
// 		NewFunction(a, NewFunction(a, a))
func NewFunction(ts ...hm.Type) *Function {
	if len(ts) < 2 {
		panic("Expected at least 2 input types")
	}

	retVal := borrowFn()
	retVal.A = ts[0]

	if len(ts) > 2 {
		retVal.B = NewFunction(ts[1:]...)
	} else {
		retVal.B = ts[1]
	}
	return retVal
}

func (t *Function) Name() string                       { return "→" }
func (t *Function) Apply(sub hm.Subs) hm.Substitutable { return (*Function)((*hm.Pair)(t).Apply(sub)) }
func (t *Function) FreeTypeVar() hm.TypeVarSet         { return ((*hm.Pair)(t)).FreeTypeVar() }
func (t *Function) Format(s fmt.State, c rune)         { fmt.Fprintf(s, "%v → %v", t.A, t.B) }
func (t *Function) String() string                     { return fmt.Sprintf("%v", t) }
func (t *Function) Normalize(k, v hm.TypeVarSet) (hm.Type, error) {
	var a, b hm.Type
	var err error
	if a, err = t.A.Normalize(k, v); err != nil {
		return nil, err
	}

	if b, err = t.B.Normalize(k, v); err != nil {
		return nil, err
	}

	return NewFunction(a, b), nil
}
func (t *Function) Types() hm.Types { return ((*hm.Pair)(t)).Types() }

func (t *Function) Eq(other hm.Type) bool {
	if ot, ok := other.(*Function); ok {
		return ot.A.Eq(t.A) && ot.B.Eq(t.B)
	}
	return false
}

// Other methods (accessors mainly)

// Arg returns the type of the function argument
func (t *Function) Arg() hm.Type { return t.A }

// Ret returns the return type of a function. If recursive is true, it will get the final return type
func (t *Function) Ret(recursive bool) hm.Type {
	if !recursive {
		return t.B
	}

	if fnt, ok := t.B.(*Function); ok {
		return fnt.Ret(recursive)
	}

	return t.B
}

// FlatTypes returns the types in FunctionTypes as a flat slice of types. This allows for easier iteration in some applications
func (t *Function) FlatTypes() hm.Types {
	retVal := hm.BorrowTypes(8) // start with 8. Can always grow
	retVal = retVal[:0]

	if a, ok := t.A.(*Function); ok {
		ft := a.FlatTypes()
		retVal = append(retVal, ft...)
		hm.ReturnTypes(ft)
	} else {
		retVal = append(retVal, t.A)
	}

	if b, ok := t.B.(*Function); ok {
		ft := b.FlatTypes()
		retVal = append(retVal, ft...)
		hm.ReturnTypes(ft)
	} else {
		retVal = append(retVal, t.B)
	}
	return retVal
}

// Clone implenents cloner
func (t *Function) Clone() interface{} {
	p := (*hm.Pair)(t)
	cloned := p.Clone()
	return (*Function)(cloned)
}
