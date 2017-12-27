package hmtypes

import "github.com/chewxy/hm"

// Pair is a convenient structural abstraction for types that are composed of two types.
// Depending on use cases, it may be useful to embed Pair, or define a new type base on *Pair.
//
// Pair partially implements hm.Substitutable, but with very specific semantics -
// it's useful for a small subset of types like function types, or supertypes.
// See the documentation for Apply and FreeTypeVar.
type Pair struct {
	A, B hm.Type
}

// Apply applies a substitution on both the first and second types of the Pair.
func (t *Pair) Apply(sub hm.Subs) {
	t.A = t.A.Apply(sub).(hm.Type)
	t.B = t.B.Apply(sub).(hm.Type)
}

// FreeTypeVar returns a set of free (unbound) type variables.
func (t Pair) FreeTypeVar() hm.TypeVarSet { return t.A.FreeTypeVar().Union(t.B.FreeTypeVar()) }

// Monuple is a convenient structural abstraction for types that are composed of one type.
//
// Monuple implements hm.Substitutable, but with very specific semantics -
// It's useful for singly polymorphic types like arrays, linear types, reference types, etc
type Monuple struct {
	A hm.Type
}

// Apply applies a substitution to the monuple type.
func (t Monuple) Apply(subs hm.Subs) hm.Substitutable { t.A = t.A.Apply(subs).(hm.Type); return t }

// FreeTypeVar returns the set of free type variables in the monuple.
func (t Monuple) FreeTypeVar() hm.TypeVarSet { return t.A.FreeTypeVar() }
