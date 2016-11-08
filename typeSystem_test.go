package hm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrune(t *testing.T) {
	assert := assert.New(t)

	var tv0 TypeVariable
	var t0 Type
	var pruned Type

	t.Log("Empty Type Variable")
	tv0 = NewTypeVar("a")
	pruned = Prune(tv0)
	assert.Equal(tv0, pruned)

	t.Log("Type Var with instance")
	t0 = nathaniel
	tv0 = NewTypeVar("a", WithInstance(t0))
	pruned = Prune(tv0)
	assert.Equal(t0, pruned)
}

var unifyTests = []struct {
	name string
	a    Type
	b    Type

	retA Type
	retB Type
	e    bool // does it error?
}{
	{"a ~ empty", NewTypeVar("a"), TypeVariable{}, NewTypeVar("a"), TypeVariable{}, false},
	{"empty ~ a", TypeVariable{}, NewTypeVar("a"), nil, nil, true},
	{"a ~ a (recursive unification)", NewTypeVar("a"), NewTypeVar("a"), nil, nil, true},
	{"a ~ b", NewTypeVar("a"), NewTypeVar("b"), NewTypeVar("a", WithInstance(NewTypeVar("b"))), NewTypeVar("b"), false},
	{"a ~ nathaniel", NewTypeVar("a"), nathaniel, NewTypeVar("a", WithInstance(nathaniel)), nathaniel, false},
	{"nathaniel ~ a", nathaniel, NewTypeVar("a"), nathaniel, NewTypeVar("a", WithInstance(nathaniel)), false},

	// type op ~ type op
	{"nathaniel ~ nathaniel", nathaniel, nathaniel, nathaniel, nathaniel, false},
	{"List a ~ List nathaniel", list{NewTypeVar("a")}, list{nathaniel}, list{nathaniel}, list{nathaniel}, false},
	{"List a ~ GoateeList nathaniel", list{NewTypeVar("a")}, mirrorUniverseList{list{nathaniel}}, nil, nil, true},

	{"malformed ~ a", malformed{}, NewTypeVar("a"), nil, nil, true},
	{"nathaniel ~ malformed{}", nathaniel, malformed{}, nil, nil, true},

	// unsure of what the correct answer should be...
	// {"a ~ malformed", NewTypeVar("a"), malformed{}, nil, nil, true},
}

func TestUnify(t *testing.T) {
	assert := assert.New(t)
	var t0, t1 Type
	var u0, u1 Type // unified
	var err error

	for _, uts := range unifyTests {
		t0 = uts.a
		t1 = uts.b
		u0, u1, err = Unify(t0, t1)
		switch {
		case err == nil && (uts.retA == nil && uts.retB == nil):
			t.Errorf("Test %q - Expected an error: %v | u0: %#v, u1: %#v", uts.name, err, u0, u1)
			continue
		case err != nil && (uts.retA != nil && uts.retB != nil):
			t.Errorf("Test %q errored: %v ", uts.name, err)
			continue
		case err != nil && (uts.retA == nil && uts.retB == nil):
			// all good, an error was expected, and an error was returned
			continue
		}

		if uts.retA == nil {
			assert.Nil(u0)
		} else {
			assert.Equal(uts.retA, u0, "Test: %q\nWant:\n%# v  \nGot:\n%# v", uts.name, uts.retA, u0)
		}

		if uts.retB == nil {
			assert.Nil(u1)
		} else {
			assert.Equal(uts.retB, u1, "Test: %q\nWant:\n%# v  \nGot:\n%# v", uts.name, uts.retB, u1)
		}
	}

}
