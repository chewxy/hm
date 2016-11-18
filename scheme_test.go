package hm

import "testing"

func TestSchemeBasics(t *testing.T) {
	s := new(Scheme)
	s.tvs = TypeVarSet{'a', 'b'}
	s.t = NewFnType(TypeVariable('c'), proton)

	sub := mSubs{
		'a': proton,
		'b': neutron,
		'c': electron,
	}

	s2 := s.Apply(sub).(*Scheme)
	if s2 != s {
		t.Errorf("Different pointers")
	}

	if !s.tvs.Equals(TypeVarSet{'a', 'b'}) {
		t.Error("TypeVarSet mutated")
	}

	if s.t != NewFnType(electron, proton) {
		t.Error("Application failed")
	}

	s = new(Scheme)
	s.tvs = TypeVarSet{'a', 'b'}
	s.t = NewFnType(TypeVariable('c'), proton)

	ftv := s.FreeTypeVar()

	if !ftv.Equals(TypeVarSet{'c'}) {
		t.Errorf("Expected ftv: {'c'}. Got %v instead", ftv)
	}
}

func TestSchemeNormalize(t *testing.T) {
	s := new(Scheme)
	s.tvs = TypeVarSet{'c', 'z', 'd'}
	s.t = NewFnType(TypeVariable('a'), TypeVariable('c'))

	err := s.normalize()
	if err != nil {
		t.Error(err)
	}

	if !s.tvs.Equals(TypeVarSet{'a', 'b'}) {
		t.Errorf("Expected: TypeVarSet{'a','b'}. Got: %v", s.tvs)
	}
}
