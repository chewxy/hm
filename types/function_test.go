package hmtypes

import (
	"testing"

	"github.com/chewxy/hm"
	"github.com/stretchr/testify/assert"
)

func TestFunctionTypeBasics(t *testing.T) {
	fnType := NewFnType(hm.TypeVariable('a'), hm.TypeVariable('a'), hm.TypeVariable('a'))
	if fnType.Name() != "→" {
		t.Errorf("FunctionType should have \"→\" as a name. Got %q instead", fnType.Name())
	}

	if fnType.String() != "a → a → a" {
		t.Errorf("Expected \"a → a → a\". Got %q instead", fnType.String())
	}

	if !fnType.Arg().Eq(hm.TypeVariable('a')) {
		t.Error("Expected arg of function to be 'a'")
	}

	if !fnType.Ret(false).Eq(NewFnType(hm.TypeVariable('a'), hm.TypeVariable('a'))) {
		t.Error("Expected ret(false) to be a → a")
	}

	if !fnType.Ret(true).Eq(hm.TypeVariable('a')) {
		t.Error("Expected final return type to be 'a'")
	}

	// a very simple fn
	fnType = NewFnType(hm.TypeVariable('a'), hm.TypeVariable('a'))
	if !fnType.Ret(true).Eq(hm.TypeVariable('a')) {
		t.Error("Expected final return type to be 'a'")
	}

	ftv := fnType.FreeTypeVar()
	if len(ftv) != 1 {
		t.Errorf("Expected only one free type var")
	}

	for _, fas := range fnApplyTests {
		fn := fas.fn.Apply(fas.sub).(*FunctionType)
		if !fn.Eq(fas.expected) {
			t.Errorf("Expected %v. Got %v instead", fas.expected, fn)
		}
	}

	// bad shit
	f := func() {
		NewFnType(hm.TypeVariable('a'))
	}
	assert.Panics(t, f)
}

var fnApplyTests = []struct {
	fn  *FunctionType
	sub hm.Subs

	expected *FunctionType
}{
	{NewFnType(hm.TypeVariable('a'), hm.TypeVariable('a')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, proton)},
	{NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, neutron)},
	{NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b')), mSubs{'c': proton, 'd': neutron}, NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b'))},
	{NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b')), mSubs{'a': proton, 'c': neutron}, NewFnType(proton, hm.TypeVariable('b'))},
	{NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b')), mSubs{'c': proton, 'b': neutron}, NewFnType(hm.TypeVariable('a'), neutron)},
	{NewFnType(electron, proton), mSubs{'a': proton, 'b': neutron}, NewFnType(electron, proton)},

	// a -> (b -> c)
	{NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b'), hm.TypeVariable('a')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, neutron, proton)},
	{NewFnType(hm.TypeVariable('a'), hm.TypeVariable('a'), hm.TypeVariable('b')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, proton, neutron)},
	{NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b'), hm.TypeVariable('c')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, neutron, hm.TypeVariable('c'))},
	{NewFnType(hm.TypeVariable('a'), hm.TypeVariable('c'), hm.TypeVariable('b')), mSubs{'a': proton, 'b': neutron}, NewFnType(proton, hm.TypeVariable('c'), neutron)},

	// (a -> b) -> c
	{NewFnType(NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b')), hm.TypeVariable('a')), mSubs{'a': proton, 'b': neutron}, NewFnType(NewFnType(proton, neutron), proton)},
}

func TestFunctionType_FlatTypes(t *testing.T) {
	fnType := NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b'), hm.TypeVariable('c'))
	ts := fnType.FlatTypes()
	correct := hm.Types{hm.TypeVariable('a'), hm.TypeVariable('b'), hm.TypeVariable('c')}
	assert.Equal(t, ts, correct)

	fnType2 := NewFnType(fnType, hm.TypeVariable('d'))
	correct = append(correct, hm.TypeVariable('d'))
	ts = fnType2.FlatTypes()
	assert.Equal(t, ts, correct)
}

func TestFunctionType_Clone(t *testing.T) {
	fnType := NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b'), hm.TypeVariable('c'))
	assert.Equal(t, fnType.Clone(), fnType)

	rec := NewRecordType("", hm.TypeVariable('a'), NewFnType(hm.TypeVariable('a'), hm.TypeVariable('b')), hm.TypeVariable('c'))
	fnType = NewFnType(rec, rec)
	assert.Equal(t, fnType.Clone(), fnType)
}
