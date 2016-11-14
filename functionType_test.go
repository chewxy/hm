package hm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fnCloneTests []*FunctionType

func init() {
	fnCloneTests = []*FunctionType{
		NewFnType(electron, list{photon}),
		NewFnType(list{photon}, electron),
		NewFnType(electron, NewTypeVar("a")),
		NewFnType(NewTypeVar("a"), electron),
	}
}

func TestFnBasics(t *testing.T) {
	var t0 *FunctionType

	// proton → electron
	name := "proton → electron"
	t0 = NewFnType(proton, electron)
	if t0.Name() != "→" {
		t.Error("Expected the name of a FunctionType to be \"→\"")
	}

	if fmt.Sprintf("%v", t0) != name {
		t.Errorf("Basic Format error. Got %q", fmt.Sprintf("%v", t0))
	}

	if t0.String() != name {
		t.Errorf("Basic String error: Got %q", t0.String())
	}

	if t0.Contains(NewTypeVar("a")) {
		t.Error("A Function Type that has no type variables shouldn't contain type variables")
	}

	// a → b → electron
	t0 = NewFnType(NewTypeVar("a"), NewTypeVar("b"), electron)
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
		NewFnType(NewTypeVar("b"), electron),
	}
	assert.EqualValues(t, correct, t0.Types())

	// equalities
	t1 := new(FunctionType)
	t1.ts[0] = correct[0]
	t1.ts[1] = correct[1]

	if !t0.Eq(t1) {
		t.Error("Expected them to be the same")
	}

	t1.ts[1] = electron
	if t0.Eq(t1) {
		t.Error("%v should not be equal to %v", t0, t1)
	}

	// set type
	// var top TypeOp
	// top = t0.SetTypes(proton, electron)
	// if t0 != top.(*FunctionType) {
	// 	t.Error("The return pointers should be the same")
	// }

	// bad shit
	// f := func() {
	// 	t0.SetTypes(proton, electron, neutron, photon)
	// }
	// assert.Panics(t, f)

	f := func() {
		NewFnType(proton)
	}
	assert.Panics(t, f)
}

func TestFnTypeReplace(t *testing.T) {
	assert := assert.New(t)
	var t0, correct *FunctionType
	var tv0, tv1 TypeVariable

	tv0 = NewTypeVar("a")
	t0 = NewFnType(tv0, tv0)
	t0.Replace(tv0, proton)
	correct = NewFnType(proton, proton)
	assert.Equal(correct, t0)

	tv1 = NewTypeVar("b")
	t0 = NewFnType(tv0, tv1)
	t0.Replace(tv0, proton)
	correct = NewFnType(proton, tv1)
	assert.Equal(correct, t0)

	tv1 = NewTypeVar("a")
	t0 = NewFnType(tv0, tv1)
	t0.Replace(tv0, proton)
	correct = NewFnType(proton, proton)
	assert.Equal(correct, t0)
}

func TestFnTypeSpecials(t *testing.T) {
	assert := assert.New(t)
	var t0 *FunctionType

	// electron anti-excitation example
	// (the slightly more correct representation is `electron → (electron, photon)`
	// given you know, the electron has simply moved to a lower state of energy)

	// electron → photon
	t0 = NewFnType(electron, photon)
	assert.Equal(Types{electron, photon}, t0.TypesRec())
	assert.Equal(photon, t0.ReturnType())

	// annihilation time!
	// annihilation :: positron → electron → photon → (quark, antiquark, gluon)
	// but for the purpose of this test we won't have antiquarks, gluons and tuple that represents them

	// positron → electron → photon → quark
	t0 = NewFnType(positron, electron, photon, quark)
	assert.Equal(Types{positron, electron, photon, quark}, t0.TypesRec())
	assert.Equal(quark, t0.ReturnType())
}

func TestFnTypeClone(t *testing.T) {
	assert := assert.New(t)
	for _, ft := range fnCloneTests {
		if ft == ft.Clone() {
			t.Error("Cloning of *FunctionType should not yield the same pointer")
			continue
		}
		assert.Equal(ft, ft.Clone())
	}
}
