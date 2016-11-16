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
func Unify(t1, t2 Type) (retVal1, retVal2 Type, err error) {
	logf("Unifying %v and %v", t1, t2)
	enterLoggingContext()
	defer leaveLoggingContext()

	a := Prune(t1)
	b := Prune(t2)

	switch at := a.(type) {
	case *TypeVariable:
		if err = UnifyVar(at, b); err != nil {
			return
		}
		retVal1 = at
		retVal2 = b

	case TypeOp:
		switch bt := b.(type) {
		case *TypeVariable:
			// note the order change
			if err = UnifyVar(bt, at); err != nil {
				return
			}
			retVal1 = at
			retVal2 = bt
		case TypeOp:
			atypes := at.Types()
			btypes := bt.Types()
			if at.Name() != bt.Name() || len(atypes) != len(btypes) {
				err = errors.Errorf(typeMismatch, a, b)
				return
			}

			if len(atypes) == 1 {
				defer ReturnTypes1(atypes)
			}
			if len(btypes) == 1 {
				defer ReturnTypes1(btypes)
			}

			for i, att := range atypes {
				logf("att: %#v btt: %#v", att, btypes[i])
				att = att.Prune()
				btt := btypes[i].Prune()
				logf("PRUNED att: %#v btt: %#v", att, btypes[i])
				if att, btt, err = Unify(att, btt); err != nil {
					return
				}

				logf("i: %v att %#v, btt %#v", i, att, btt)
				logf("aty %v bty %v", atypes, btypes)
				atypes[i] = att.Prune()
				btypes[i] = btt.Prune()
				logf("PRUNED2: %v %v", atypes[i], btypes[i])
				logf("aty %v bty %v", atypes, btypes)

			}

			retVal1 = at.New(atypes...)
			retVal2 = bt.New(btypes...)
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
func UnifyVar(tv *TypeVariable, t Type) (err error) {
	logf("Unifying var %v and %v", tv, t)
	enterLoggingContext()
	defer leaveLoggingContext()

	if tv.IsEmpty() {
		return errors.Errorf(undefinedTV)
	}

	if ttv, ok := t.(*TypeVariable); ok {
		if ttv.IsEmpty() {
			return nil
		}

		if !tv.Eq(ttv) {
			unioned := tv.constraints.Union(ttv.constraints)
			ttv.constraints = unioned
			tv.constraints = unioned
			goto setinstance
		} else {
			return errors.Errorf(recursiveUnification, tv, ttv)
		}
	}

	if t.Contains(tv) {
		return errors.Errorf(recursiveUnification, tv, t)
	}

setinstance:
	tv.instance = t
	return nil
}

// Prune returns the defining instance of T
func Prune(t Type) Type {
	if tv, ok := t.(*TypeVariable); ok {
		logf("Is TV")
		if tv.instance != nil {
			tv.instance = Prune(tv.instance)
			return tv.instance
		}
	}
	return t
}
