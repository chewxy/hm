package hm

import "github.com/pkg/errors"

// Unify unifies the two types.
// These are the rules:
//
// Type Constants and Type Constants
//
// Type constants (atomic types) have no substitution
//		c ~ c : []
//
// Type Variables and Type Variables
//
// Type variables have no substitutions if there are no instances:
// 		a ~ a : []
//
// Default Unification
//
// if type variable 'a' is not in 'T', then unification is simple: replace all instances of 'a' with 'T'
// 		     a ∉ T
//		---------------
//		 a ~ T : [a/T]
//
// The more complicated constructor unification and arrow unification isn't quite covered yet.
func Unify(t1, t2 Type) (retVal1, retVal2 Type, replacements map[TypeVariable]Type, err error) {
	a := Prune(t1)
	b := Prune(t2)

	switch at := a.(type) {
	case TypeVariable:
		if retVal1, retVal2, err = UnifyVar(at, b); err != nil {
			return
		}
		if replacements == nil {
			replacements = make(map[TypeVariable]Type)
		}

		replacements[at] = retVal1
	case TypeOp:
		switch bt := b.(type) {
		case TypeVariable:
			// note the order change
			if retVal2, retVal1, err = UnifyVar(bt, at); err != nil {
				return
			}

			if replacements == nil {
				replacements = make(map[TypeVariable]Type)
			}

			replacements[bt] = retVal2
		case TypeOp:
			atypes := at.Types()
			btypes := bt.Types()
			if at.Name() != bt.Name() || len(atypes) != len(btypes) {
				err = errors.Errorf(typeMismatch, a, b)
				return
			}

			var t_a, t_b Type
			for i := 0; i < len(atypes); i++ {
				t_a = atypes[i]
				t_b = btypes[i]

				var t_a2, t_b2 Type
				var r2 map[TypeVariable]Type
				if t_a2, t_b2, r2, err = Unify(t_a, t_b); err != nil {
					return
				}

				if replacements == nil {
					replacements = r2
				} else {
					for k, v := range r2 {
						replacements[k] = v
					}
				}

				pt_a2 := Prune(t_a2)
				pt_b2 := Prune(t_b2)

				at = at.Replace(t_a, pt_a2)
				bt = bt.Replace(t_b, pt_b2)

				for k, v := range replacements {
					at = at.Replace(k, v)
					bt = bt.Replace(k, v)
				}

				if tv, ok := t_a.(TypeVariable); ok {
					replacements[tv] = pt_a2
				}

				if tv, ok := t_b.(TypeVariable); ok {
					replacements[tv] = pt_b2
				}

				atypes = at.Types()
				btypes = bt.Types()
			}

			retVal1 = at
			retVal2 = bt
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
	var unioned *TypeClassSet
	if ttv, ok := t.(TypeVariable); ok {
		if ttv.IsEmpty() {
			return
		}

		if t.Eq(ttv) {
			if tv.constraints == nil {
			}
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
