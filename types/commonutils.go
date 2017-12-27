package hmtypes

import "github.com/chewxy/hm"

// Pair is a convenient structural abstraction for types that are composed of two types.
// Depending on use cases, it may be useful to embed Pair, or define a new type base on *Pair.
//
// Pair partially implements hm.Type, as the intention is merely for syntactic abstraction
//
// It has very specific semantics -
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

// Types returns all the types of the Pair's constituents
func (t *Pair) Types() hm.Types {
	retVal := hm.BorrowTypes(2)
	retVal[0] = t.A
	retVal[1] = t.B
	return retVal
}

// FreeTypeVar returns a set of free (unbound) type variables.
func (t *Pair) FreeTypeVar() hm.TypeVarSet { return t.A.FreeTypeVar().Union(t.B.FreeTypeVar()) }

// Clone implements Cloner
func (t *Pair) Clone() *Pair {
	retVal := borrowPair()

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

// Monuple is a convenient structural abstraction for types that are composed of one type.
//
// Monuple implements hm.Substitutable, but with very specific semantics -
// It's useful for singly polymorphic types like arrays, linear types, reference types, etc
type Monuple struct {
	T hm.Type
}

// Apply applies a substitution to the monuple type.
func (t Monuple) Apply(subs hm.Subs) Monuple {
	t.T = t.T.Apply(subs).(hm.Type)
	return t
}

// FreeTypeVar returns the set of free type variables in the monuple.
func (t Monuple) FreeTypeVar() hm.TypeVarSet { return t.T.FreeTypeVar() }

// Normalize is the method to normalize all type variables
func (t Monuple) Normalize(k, v hm.TypeVarSet) (Monuple, error) {
	var t2 hm.Type
	var err error
	if t2, err = t.T.Normalize(k, v); err != nil {
		return nil, err
	}
	t.T = t2
	return t, nil
}
