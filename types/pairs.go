package hmtypes

import (
	"fmt"

	"github.com/chewxy/hm"
)

var (
	_ hm.Type = &Choice{}
	_ hm.Type = &Super{}
	_ hm.Type = &Application{}
)

// pair types

// Choice is the type of choice of algorithm to use within a class method.
//
// Imagine how one would implement a class in an OOP language.
// Then imagine how one would implement method overloading for the class.
// The typical approach is name mangling followed by having a jump table.
//
// Now consider OOP classes and the ability to override methods, based on subclassing ability.
// The typical approach to this is to use a Vtable.
//
// Both overloading and overriding have a general notion: a jump table of sorts.
// How does one type such a table?
//
// By using  Choice.
//
// The first type is the key of either the vtable or the name mangled table.
// The second type is the value of the table.
type Choice hm.Pair

func (t *Choice) Name() string                       { return ":" }
func (t *Choice) Apply(sub hm.Subs) hm.Substitutable { ((*hm.Pair)(t)).Apply(sub); return t }
func (t *Choice) FreeTypeVar() hm.TypeVarSet         { return ((*hm.Pair)(t)).FreeTypeVar() }
func (t *Choice) Format(s fmt.State, c rune)         { fmt.Fprintf(s, "%v : %v", t.A, t.B) }
func (t *Choice) String() string                     { return fmt.Sprintf("%v", t) }

func (t *Choice) Normalize(k hm.TypeVarSet, v hm.TypeVarSet) (hm.Type, error) {
	panic("not implemented")
}

func (t *Choice) Types() hm.Types { return ((*hm.Pair)(t)).Types() }

func (t *Choice) Eq(other hm.Type) bool {
	if ot, ok := other.(*Choice); ok {
		return ot.A.Eq(t.A) && ot.B.Eq(t.B)
	}
	return false
}

func (t *Choice) Clone() interface{} { return (*Choice)((*hm.Pair)(t).Clone()) }

func (t *Choice) Pair() *hm.Pair { return (*hm.Pair)(t) }

// Super is the inverse of Choice. It allows for supertyping functions.
//
// Supertyping is typically  implemented as a adding an entry to the vtable/mangled table.
// But there needs to be a separate accounting structure to keep account of the types.
//
// This is where Super comes in.
type Super hm.Pair

func (t *Super) Name() string                       { return "§" }
func (t *Super) Apply(sub hm.Subs) hm.Substitutable { ((*hm.Pair)(t)).Apply(sub); return t }
func (t *Super) FreeTypeVar() hm.TypeVarSet         { return ((*hm.Pair)(t)).FreeTypeVar() }
func (t *Super) Format(s fmt.State, c rune)         { fmt.Fprintf(s, "%v §: %v", t.A, t.B) }
func (t *Super) String() string                     { return fmt.Sprintf("%v", t) }

func (t *Super) Normalize(k hm.TypeVarSet, v hm.TypeVarSet) (hm.Type, error) {
	panic("not implemented")
}

func (t *Super) Types() hm.Types { return ((*hm.Pair)(t)).Types() }

func (t *Super) Eq(other hm.Type) bool {
	if ot, ok := other.(*Super); ok {
		return ot.A.Eq(t.A) && ot.B.Eq(t.B)
	}
	return false
}

func (t *Super) Clone() interface{} { return (*Super)((*hm.Pair)(t).Clone()) }

func (t *Super) Pair() *hm.Pair { return (*hm.Pair)(t) }

// Application is the pre-unified type for a function application.
// In a simple HM system this would not be needed as the type of an
// application expression would be found during the unification phase of
// the expression.
//
// In advanced systems where unification may be done concurrently, this would
// be required, as a "thunk" of sorts for the type system.
type Application hm.Pair

func (t *Application) Name() string                       { return "•" }
func (t *Application) Apply(sub hm.Subs) hm.Substitutable { ((*hm.Pair)(t)).Apply(sub); return t }
func (t *Application) FreeTypeVar() hm.TypeVarSet         { return ((*hm.Pair)(t)).FreeTypeVar() }
func (t *Application) Format(s fmt.State, c rune)         { fmt.Fprintf(s, "%v • %v", t.A, t.B) }
func (t *Application) String() string                     { return fmt.Sprintf("%v", t) }

func (t *Application) Normalize(k hm.TypeVarSet, v hm.TypeVarSet) (hm.Type, error) {
	panic("not implemented")
}

func (t *Application) Types() hm.Types { return ((*hm.Pair)(t)).Types() }

func (t *Application) Eq(other hm.Type) bool {
	if ot, ok := other.(*Application); ok {
		return ot.A.Eq(t.A) && ot.B.Eq(t.B)
	}
	return false
}

func (t *Application) Clone() interface{} { return (*Application)((*hm.Pair)(t).Clone()) }

func (t *Application) Pair() *hm.Pair { return (*hm.Pair)(t) }
