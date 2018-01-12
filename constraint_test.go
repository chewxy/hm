package hm

import "testing"

func TestConstraint(t *testing.T) {
	c := Constraint{
		A: TypeVariable('a'),
		B: NewFnType(TypeVariable('b'), TypeVariable('c')),
	}

	ftv := c.FreeTypeVar()
	if !ftv.Equals(TypeVarSet{TypeVariable('a'), TypeVariable('b'), TypeVariable('c')}) {
		t.Error("the free type variables of a Constraint is not as expected")
	}

	subs := mSubs{
		'a': NewFnType(proton, proton),
		'b': proton,
		'c': neutron,
	}

	c = c.Apply(subs).(Constraint)
	if !c.A.Eq(NewFnType(proton, proton)) {
		t.Errorf("c.a: %v", c)
	}

	if !c.B.Eq(NewFnType(proton, neutron)) {
		t.Errorf("c.b: %v", c)
	}
}
