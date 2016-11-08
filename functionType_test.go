package hm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFnBasics(t *testing.T) {
	var t0 *FunctionType

	// nathaniel → adam
	t0 = NewFnType(nathaniel, adam)
	if t0.Name() != "→" {
		t.Error("Expected the name of a FunctionType to be \"→\"")
	}

	if fmt.Sprintf("%v", t0) != "atom(0) → atom(1)" {
		t.Errorf("Basic Format error. Got %q", fmt.Sprintf("%v", t0))
	}

	if t0.String() != "atom(0) → atom(1)" {
		t.Errorf("Basic String error: Got %q", t0.String())
	}

	if t0.Contains(NewTypeVar("a")) {
		t.Error("A Function Type that has no type variables shouldn't contain type variables")
	}

	// a → b → adam
	t0 = NewFnType(NewTypeVar("a"), NewTypeVar("b"), adam)
	if !t0.Contains(NewTypeVar("a")) {
		t.Errorf("Expected %v to contain Type Var `a`", t0)
	}

	if !t0.Contains(NewTypeVar("b")) {
		t.Errorf("Expected %v to contain Type Var `b`", t0)
	}

	if t0.Contains(NewTypeVar("c")) {
		t.Errorf("%v shouldn't contain Type Var `c`", t0)
	}

	correct := Types{
		NewTypeVar("a"),
		NewFnType(NewTypeVar("b"), adam),
	}
	assert.EqualValues(t, correct, t0.Types())

	// equalities
	t1 := new(FunctionType)
	t1.ts[0] = correct[0]
	t1.ts[1] = correct[1]

	if !t0.Eq(t1) {
		t.Error("Expected them to be the same")
	}

	t1.ts[1] = adam
	if t0.Eq(t1) {
		t.Error("%v should not be equal to %v", t0, t1)
	}

	// set type
	var top TypeOp
	top = t0.SetTypes(nathaniel, adam)
	if t0 != top.(*FunctionType) {
		t.Error("The return pointers should be the same")
	}

	// bad shit
	f := func() {
		t0.SetTypes(nathaniel, adam, jonathan, osterman)
	}
	assert.Panics(t, f)

	f = func() {
		NewFnType(nathaniel)
	}
	assert.Panics(t, f)
}
