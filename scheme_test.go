package hm

import (
	"fmt"
	"testing"
)

func TestSchemeBasics(t *testing.T) {
	s := new(Scheme)
	s.tvs = TypeVarSet{'a', 'b'}
	s.t = NewFnType(TypeVariable('c'), proton)

	sub := mSubs{
		'a': proton,
		'b': neutron,
		'c': electron,
	}

	s2 := s.Apply(nil).(*Scheme)
	if s2 != s {
		t.Errorf("Different pointers")
	}

	s2 = s.Apply(sub).(*Scheme)
	if s2 != s {
		t.Errorf("Different pointers")
	}

	if !s.tvs.Equals(TypeVarSet{'a', 'b'}) {
		t.Error("TypeVarSet mutated")
	}

	if !s.t.Eq(NewFnType(electron, proton)) {
		t.Error("Application failed")
	}

	s = new(Scheme)
	s.tvs = TypeVarSet{'a', 'b'}
	s.t = NewFnType(TypeVariable('c'), proton)

	ftv := s.FreeTypeVar()

	if !ftv.Equals(TypeVarSet{'c'}) {
		t.Errorf("Expected ftv: {'c'}. Got %v instead", ftv)
	}

	// format
	if fmt.Sprintf("%v", s) != "∀[a, b]: c → proton" {
		t.Errorf("Scheme format is wrong.: Got %q", fmt.Sprintf("%v", s))
	}

	// Polytype scheme.Type
	T, isMono := s.Type()
	if isMono {
		t.Errorf("%v is supposed to be a polytype. It shouldn't return true", s)
	}
	if !T.Eq(NewFnType(TypeVariable('c'), proton)) {
		t.Error("Wrong type returned by scheme")
	}
}

func TestSchemeNormalize(t *testing.T) {
	s := new(Scheme)
	s.tvs = TypeVarSet{'c', 'z', 'd'}
	s.t = NewFnType(TypeVariable('a'), TypeVariable('c'))

	err := s.Normalize()
	if err != nil {
		t.Error(err)
	}

	if !s.tvs.Equals(TypeVarSet{'a', 'b'}) {
		t.Errorf("Expected: TypeVarSet{'a','b'}. Got: %v", s.tvs)
	}
}
