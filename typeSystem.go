package hm

import "github.com/pkg/errors"

// Unify unifies the two types.
// These are the rules:
//
// type constants (atomic types) have no substitution
//		c ~ c : []
//
// type variables have no substitutions if there are no instances
// 		a ~ a : []
//
// if type variable 'a' is not in 'T', then unification is simple: replace all instances of 'a' with 'T'
// 		a ∉ T
//		----------
//		a ~ T : [a/T]
//
// The more complicated constructor unification and arrow unification isn't quite covered yet.
func Unify(t1, t2 Type) (retVal1, retVal2 Type, err error) {
	a := Prune(t1)
	b := Prune(t2)

	switch at := a.(type) {
	case TypeVariable:
		retVal1, retVal2, err = UnifyVar(at, b)
	case TypeOp:
		switch bt := b.(type) {
		case TypeVariable:
			retVal2, retVal1, err = UnifyVar(bt, at) // note the order change
			return
		case TypeOp:
			atypes := at.Types()
			btypes := bt.Types()
			if at.Name() != bt.Name() || len(atypes) != len(btypes) {
				err = errors.Errorf(typeMismatch, a, b)
				return
			}

			var t_a, t_b Type
			var i int
			for i, t_a = range atypes {
				if t_a, t_b, err = Unify(t_a, btypes[i]); err != nil {
					return
				}

				atypes[i] = Prune(t_a)
				btypes[i] = Prune(t_b)
			}

			retVal1 = at.SetTypes(atypes...)
			retVal2 = bt.SetTypes(btypes...)
			return
		default:
			err = errors.Errorf(nyi, "Unify of TypeOp ", b, b)
			return
		}

	default:
		err = errors.Errorf(nu, t1, t2)
		return
	}
	return
}

// UnifyVar unifies a TypeVariable and a Type.
func UnifyVar(tv TypeVariable, t Type) (ret1, ret2 Type, err error) {
	if tv.IsEmpty() {
		err = errors.Errorf(undefinedTV)
		return
	}
	ret1 = tv
	ret2 = t

	var unioned TypeClassSet
	if ttv, ok := t.(TypeVariable); ok {
		if ttv.IsEmpty() {
			return
		}

		if t.Eq(ttv) {
			unioned = tv.constraints.Union(ttv.constraints)

			tv.constraints = unioned
			ttv.constraints = unioned
			ret2 = ttv
		}

	}

	if ret2.Contains(tv) {
		err = errors.Errorf(recursiveUnification, tv, t)
		return
	}

	tv.instance = ret2
	ret1 = tv
	return
}

// Prune returns the defining instance of T
func Prune(t Type) Type {
	if tv, ok := t.(TypeVariable); ok {
		if tv.instance != nil {
			return Prune(tv.instance)
		}
	}
	return t
}