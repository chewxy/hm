package hm

import (
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

func TestPrune(t *testing.T) {
	assert := assert.New(t)

	var tv0 *TypeVariable
	var t0 Type
	var pruned Type

	t.Log("Empty Type Variable")
	tv0 = NewTypeVar("a")
	pruned = Prune(tv0)
	assert.Equal(tv0, pruned)

	t.Log("Type Var with instance")
	t0 = proton
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
	{"a ~ empty", NewTypeVar("a"), &TypeVariable{}, NewTypeVar("a"), &TypeVariable{}, false},
	{"empty ~ a", &TypeVariable{}, NewTypeVar("a"), nil, nil, true},
	{"empty ~ a", nil, NewTypeVar("a"), nil, nil, true},
	{"a ~ a (recursive unification)", NewTypeVar("a"), NewTypeVar("a"), nil, nil, true},
	{"a ~ b", NewTypeVar("a"), NewTypeVar("b"), NewTypeVar("a", WithInstance(NewTypeVar("b"))), NewTypeVar("b"), false},
	{"a ~ proton", NewTypeVar("a"), proton, NewTypeVar("a", WithInstance(proton)), proton, false},
	{"proton ~ a", proton, NewTypeVar("a"), proton, NewTypeVar("a", WithInstance(proton)), false},

	// type op ~ type op
	{"proton ~ proton", proton, proton, proton, proton, false},
	{"List a ~ List proton", list{NewTypeVar("a")}, list{proton}, list{proton}, list{proton}, false},
	{"List a ~ GoateeList proton", list{NewTypeVar("a")}, mirrorUniverseList{proton}, nil, nil, true},

	// function types (without references)
	{"List a → List a ~ List proton → List proton", NewFnType(list{NewTypeVar("a")}, list{NewTypeVar("a")}), NewFnType(list{proton}, list{proton}), NewFnType(list{proton}, list{proton}), NewFnType(list{proton}, list{proton}), false},
	{"List a → a ~ List proton → proton", NewFnType(list{NewTypeVar("a")}, NewTypeVar("a")), NewFnType(list{proton}, proton), NewFnType(list{proton}, proton), NewFnType(list{proton}, proton), false},
	{"malformed ~ a", malformed{}, NewTypeVar("a"), nil, nil, true},
	{"proton ~ malformed{}", proton, malformed{}, nil, nil, true},

	// unsure of what the correct answer should be...
	// {"a ~ malformed", NewTypeVar("a"), malformed{}, nil, nil, true},

	// {"a ~ empty(nil)", NewTypeVar("a"), nil, NewTypeVar("a"), nil, false},

}

func TestUnify(t *testing.T) {
	assert := assert.New(t)
	var t0, t1 Type
	var u0, u1 Type
	var err error

	for _, uts := range unifyTests {
		logf("unifying %v", uts.name)
		t0 = uts.a
		t1 = uts.b
		u0, u1, err = Unify(t0, t1)
		switch {
		case err == nil && uts.e:
			t.Errorf("Test %q - Expected an error: %v | u0: %#v, u1: %#v", uts.name, err, u0, u1)
		case err != nil && !uts.e:
			t.Errorf("Test %q errored: %v ", uts.name, err)
		}

		if uts.e {
			continue
		}

		assert.Equal(uts.retA, u0, "Test: %q (t0)\nWant:\n%# v  \nGot:\n%# v", uts.name, uts.retA, u0)
		assert.Equal(uts.retB, u1, "Test: %q (t1)\nWant:\n%# v  \nGot:\n% #v", uts.name, uts.retB, u1)
	}

	// for cases where a reference is required
	var name string
	var a, b *TypeVariable
	var e0, e1 Type

	// List a → a ~ List proton → b
	name = "List a → a ~ List proton → b"
	a = NewTypeVar("a")
	b = NewTypeVar("b")
	t0 = NewFnType(list{a}, a)
	t1 = NewFnType(list{proton}, b)
	e0 = NewFnType(list{proton}, proton)
	e1 = e0
	if u0, u1, err = Unify(t0, t1); err != nil {
		t.Errorf("Test %q error: %v", name, err)
	}
	assert.Equal(e0, u0, "Test %q (t0): Want:\n %#v\nGot: \n%# v", name, e0, u0)
	assert.Equal(e0, u0, "Test %q (t1): Want:\n %#v\nGot: \n%# v", name, e1, u1)

	// List a → a → List a ~ List proton → proton → b
	name = "List a → a → List a ~ List proton → proton → b"
	a = NewTypeVar("a")
	b = NewTypeVar("b")
	t0 = NewFnType(list{a}, a, list{a})
	t1 = NewFnType(list{proton}, proton, b)
	e0 = NewFnType(list{proton}, proton, list{proton})
	e1 = e0

	if u0, u1, err = Unify(t0, t1); err != nil {
		t.Errorf("Test %q error: %v", name, err)
	}
	assert.Equal(e0, u0, "Test %q (t0): Want:\n %# v\nGot: \n%# v", name, pretty.Formatter(e0), pretty.Formatter(u0))
	assert.Equal(e0, u0, "Test %q (t1): Want:\n %# v\nGot: \n%# v", name, e1, u1)

	// List proton → proton → b ~ List a → a → List a
	name = "List proton → proton → b ~ List a → a → List a"
	a = NewTypeVar("a")
	b = NewTypeVar("b")
	t0 = NewFnType(list{proton}, proton, b)
	t1 = NewFnType(list{a}, a, list{a})

	if u0, u1, err = Unify(t0, t1); err != nil {
		t.Errorf("Test %q error: %v", name, err)
	}
	assert.Equal(e0, u0, "Test %q (t0): Want:\n %# v\nGot: \n%# v", name, pretty.Formatter(e0), pretty.Formatter(u0))
	assert.Equal(e0, u0, "Test %q (t1): Want:\n %# v\nGot: \n%# v", name, e1, u1)

	// List proton → proton → b ~ List a → a → GoateeList a"
	name = "List proton → proton → b ~ List a → a → GoateeList a"
	a = NewTypeVar("a")
	b = NewTypeVar("b")
	t0 = NewFnType(list{proton}, proton, b)
	t1 = NewFnType(list{a}, a, mirrorUniverseList{a})
	e0 = NewFnType(list{proton}, proton, mirrorUniverseList{proton})
	e1 = e0

	if u0, u1, err = Unify(t0, t1); err != nil {
		t.Errorf("Test %q error: %v", name, err)
	}
	assert.Equal(e0, u0, "Test %q (t0): Want:\n %# v\nGot: \n%# v", name, pretty.Formatter(e0), pretty.Formatter(u0))
	assert.Equal(e0, u0, "Test %q (t1): Want:\n %# v\nGot: \n%# v", name, e1, u1)
}
