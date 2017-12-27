package hmtypes

import (
	"fmt"

	"github.com/chewxy/hm"
)

// Slice is the type of a Slice/List
type Slice Monuple

func (t Slice) Name() string                        { return "List" }
func (t Slice) Apply(subs hm.Subs) hm.Substitutable { return Slice(Monuple(t).Apply(subs)) }
func (t Slice) FreeTypeVar() hm.TypeVarSet          { return Monuple(t).FreeTypeVar() }
func (t Slice) Format(s fmt.State, c rune)          { fmt.Fprintf(s, "[]%v", t.T) }
func (t Slice) String() string                      { return fmt.Sprintf("%v", t) }
func (t Slice) Types() hm.Types                     { return hm.Types{t.T} }
func (t Slice) Eq(other hm.Type) bool {
	if ot, ok := other.(Slice); ok {
		return ot.T.Eq(t.T)
	}
	return false
}

// Linear is a linear type (i.e types that can only appear once)
type Linear Monuple

func (t Linear) Name() string                        { return "Linear" }
func (t Linear) Apply(subs hm.Subs) hm.Substitutable { return Linear(Monuple(t).Apply(subs)) }
func (t Linear) FreeTypeVar() hm.TypeVarSet          { return Monuple(t).FreeTypeVar() }
func (t Linear) Format(s fmt.State, c rune)          { fmt.Fprintf(s, "Linear[%v]", t.T) }
func (t Linear) String() string                      { return fmt.Sprintf("%v", t) }
func (t Linear) Types() hm.Types                     { return hm.Types{t.T} }
func (t Linear) Eq(other hm.Type) bool {
	if ot, ok := other.(Linear); ok {
		return ot.T.Eq(t.T)
	}
	return false
}

// Ref is a reference type (think pointers)
type Ref Monuple

func (t Ref) Name() string                        { return "Ref" }
func (t Ref) Apply(subs hm.Subs) hm.Substitutable { return Ref(Monuple(t).Apply(subs)) }
func (t Ref) FreeTypeVar() hm.TypeVarSet          { return Monuple(t).FreeTypeVar() }
func (t Ref) Format(s fmt.State, c rune)          { fmt.Fprintf(s, "*%v", t.T) }
func (t Ref) String() string                      { return fmt.Sprintf("%v", t) }
func (t Ref) Types() hm.Types                     { return hm.Types{t.T} }
func (t Ref) Eq(other hm.Type) bool {
	if ot, ok := other.(Ref); ok {
		return ot.T.Eq(t.T)
	}
	return false
}