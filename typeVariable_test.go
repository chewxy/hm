package hm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tvEq = []struct {
	name string
	a    TypeVariable
	b    Type

	equal bool
}{
	{"empty == empty", TypeVariable{}, TypeVariable{}, true},
	{"a == a", NewTypeVar("a"), NewTypeVar("a"), true},
	{"a == b", NewTypeVar("a"), NewTypeVar("b"), false},
	{"a == proton", NewTypeVar("a"), proton, false},
	{"a:proton == proton", NewTypeVar("a", WithInstance(proton)), proton, false},
	{"a:proton == a:proton", NewTypeVar("a", WithInstance(proton)), NewTypeVar("a", WithInstance(proton)), true},
	{"a:proton == b:proton", NewTypeVar("a", WithInstance(proton)), NewTypeVar("b", WithInstance(proton)), false},
	{"a:b:<nil> == a:b:<nil>", NewTypeVar("a", WithInstance(NewTypeVar("b"))), NewTypeVar("a", WithInstance(NewTypeVar("b"))), true},
}

var tvContains = []struct {
	name string
	a, b TypeVariable

	contains bool
}{
	{"empty <: empty", TypeVariable{}, TypeVariable{}, true},
	{"empty <: a", TypeVariable{}, NewTypeVar("a"), false},
	{"a <: a", NewTypeVar("a"), NewTypeVar("a"), true},
	{"a:proton <: a:proton", NewTypeVar("a", WithInstance(proton)), NewTypeVar("a", WithInstance(proton)), true},
	{"a <: a:proton", NewTypeVar("a"), NewTypeVar("a", WithInstance(proton)), false},
}

var tvStrings = []struct {
	name string
	a    TypeVariable
	s    string
	v    string
}{
	{"empty", TypeVariable{}, "''", "'':<nil>"},
	{"a", NewTypeVar("a"), "a", "a:<nil>"},
	{"a:proton", NewTypeVar("a", WithInstance(proton)), "proton", "a:proton"},
	{"a:b:proton", NewTypeVar("a", WithInstance(NewTypeVar("b", WithInstance(proton)))), "proton", "a:b:proton"},
}

func TestTypeVariableBasics(t *testing.T) {
	assert := assert.New(t)

	var tv0, tv1 TypeVariable
	var t0, t1 Type
	t.Log("Empty Type Variable")
	if ok := tv0.IsEmpty(); !ok {
		t.Error("Expected empty type variable")
	}

	t.Log("Equality tests")
	for _, tves := range tvEq {
		tv0 = tves.a
		t0 = tves.b

		if tv0.Eq(t0) != tves.equal {
			t.Errorf("Test %q error", tves.name)
		}
	}

	t.Log("Equality - same name but different instances: Panic expected")
	t1 = electron
	tv0 = NewTypeVar("a", WithInstance(t0))
	tv1 = NewTypeVar("a", WithInstance(t1))
	fail := func() {
		tv0.Eq(tv1)
	}
	assert.Panics(fail)

	t.Log("Contains")
	for _, tvcs := range tvContains {
		tv0 = tvcs.a
		tv1 = tvcs.b
		if tv0.Contains(tv1) != tvcs.contains {
			t.Errorf("Test %q error", tvcs.name)
		}
	}

	t.Log("String (for completeness sake")
	for _, tvss := range tvStrings {
		if tvss.a.String() != tvss.s {
			t.Errorf("Test %q error: Got %q", tvss.name, tvss.a.String())
		}
	}

	t.Log("TypeVar Format (for completeness sake)")
	for _, tvss := range tvStrings {
		if fmt.Sprintf("%v", tvss.a) != tvss.s {
			t.Errorf("Format(%%v) error. Got %q", fmt.Sprintf("%v", tv0))
		}
		if fmt.Sprintf("%#v", tvss.a) != tvss.v {
			t.Errorf("Format(%%#v) error. Got %q", fmt.Sprintf("%v", tv0))
		}
	}

	t.Log("TypeVar Name (for completeness sake)")
	if tv0.Name() != "a" {
		t.Error("Expected \"a\" to be the name")
	}

	t.Log("TypeVar Instance (for completeness sake")
	t0 = NewTypeVar("a", WithInstance(proton))
	assert.Equal(proton, t0.(TypeVariable).Instance())
}

func TestTVConsOpt(t *testing.T) {
	constraints := &TypeClassSet{
		s: []TypeClass{&SimpleTypeClass{}},
	}

	tv0 := NewTypeVar("a", WithConstraints(constraints))
	assert.Equal(t, constraints, tv0.constraints)
}
