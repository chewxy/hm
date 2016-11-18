package hm

import "fmt"

type Constraints []Constraint

func (cs Constraints) Apply(sub Subs) Substitutable {
	logf("Constraints: %d", len(cs))
	for i, c := range cs {
		cs[i] = c.Apply(sub).(Constraint)
	}
	return cs
}

func (cs Constraints) FreeTypeVar() TypeVarSet {
	var retVal TypeVarSet
	for _, v := range cs {
		retVal = v.FreeTypeVar().Union(retVal)
	}
	return retVal
}

func (cs Constraints) Format(state fmt.State, c rune) {
	state.Write([]byte("Constraints["))
	for i, c := range cs {
		if i < len(cs)-1 {
			fmt.Fprintf(state, "%v, ", c)
		} else {
			fmt.Fprintf(state, "%v", c)
		}
	}
	state.Write([]byte{']'})
}

type Types []Type

// func (ts Types) Apply(sub Subs) Substitutable {
// 	for i, t := range ts {
// 		ts[i] = t.Apply(sub).(Type)
// 	}
// 	return ts
// }

// func (ts Types) FreeTypeVar() TypeVarSet {
// 	var retVal TypeVarSet
// 	for _, v := range ts {
// 		retVal = v.FreeTypeVar().Union(retVal)
// 	}
// 	return retVal
// }
