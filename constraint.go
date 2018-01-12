package hm

import "fmt"

// A Constraint is well.. a constraint that says a must equal to b. It's used mainly in the constraint generation process.
type Constraint Pair

func (c Constraint) Apply(sub Subs) Substitutable   { (*Pair)(&c).Apply(sub); return c }
func (c Constraint) FreeTypeVar() TypeVarSet        { return Pair(c).FreeTypeVar() }
func (c Constraint) Format(state fmt.State, r rune) { fmt.Fprintf(state, "{%v = %v}", c.A, c.B) }
