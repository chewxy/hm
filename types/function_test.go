package hmtypes

import (
	"testing"

	"github.com/chewxy/hm"
	"github.com/stretchr/testify/assert"
)

func TestFunctionTypeBasics(t *testing.T) {
	fnType := NewFunction(hm.TypeVariable('a'), hm.TypeVariable('a'), hm.TypeVariable('a'))
	if fnType.Name() != "→" {
		t.Errorf("FunctionType should have \"→\" as a name. Got %q instead", fnType.Name())
	}

	if fnType.String() != "a → a → a" {
		t.Errorf("Expected \"a → a → a\". Got %q instead", fnType.String())
	}

	if !fnType.Arg().Eq(hm.TypeVariable('a')) {
		t.Error("Expected arg of function to be 'a'")
	}

	if !fnType.Ret(false).Eq(NewFunction(hm.TypeVariable('a'), hm.TypeVariable('a'))) {
		t.Error("Expected ret(false) to be a → a")
	}

	if !fnType.Ret(true).Eq(hm.TypeVariable('a')) {
		t.Error("Expected final return type to be 'a'")
	}

	// a very simple fn
	fnType = NewFunction(hm.TypeVariable('a'), hm.TypeVariable('a'))
	if !fnType.Ret(true).Eq(hm.TypeVariable('a')) {
		t.Error("Expected final return type to be 'a'")
	}

	ftv := fnType.FreeTypeVar()
	if len(ftv) != 1 {
		t.Errorf("Expected only one free type var")
	}

	for _, fas := range fnApplyTests {
		fn := fas.fn.Apply(fas.sub).(*Function)
		if !fn.Eq(fas.expected) {
			t.Errorf("Expected %v. Got %v instead", fas.expected, fn)
		}
	}

	// bad shit
	f := func() {
		NewFunction(hm.TypeVariable('a'))
	}
	assert.Panics(t, f)
}

var fnApplyTests = []struct {
	fn  *Function
	sub hm.Subs

	expected *Function
}{
	{NewFunction(hm.TypeVariable('a'), hm.TypeVariable('a')), mSubs{'a': proton, 'b': neutron}, NewFunction(proton, proton)},
	{NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b')), mSubs{'a': proton, 'b': neutron}, NewFunction(proton, neutron)},
	{NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b')), mSubs{'c': proton, 'd': neutron}, NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b'))},
	{NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b')), mSubs{'a': proton, 'c': neutron}, NewFunction(proton, hm.TypeVariable('b'))},
	{NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b')), mSubs{'c': proton, 'b': neutron}, NewFunction(hm.TypeVariable('a'), neutron)},
	{NewFunction(electron, proton), mSubs{'a': proton, 'b': neutron}, NewFunction(electron, proton)},

	// a -> (b -> c)
	{NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b'), hm.TypeVariable('a')), mSubs{'a': proton, 'b': neutron}, NewFunction(proton, neutron, proton)},
	{NewFunction(hm.TypeVariable('a'), hm.TypeVariable('a'), hm.TypeVariable('b')), mSubs{'a': proton, 'b': neutron}, NewFunction(proton, proton, neutron)},
	{NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b'), hm.TypeVariable('c')), mSubs{'a': proton, 'b': neutron}, NewFunction(proton, neutron, hm.TypeVariable('c'))},
	{NewFunction(hm.TypeVariable('a'), hm.TypeVariable('c'), hm.TypeVariable('b')), mSubs{'a': proton, 'b': neutron}, NewFunction(proton, hm.TypeVariable('c'), neutron)},

	// (a -> b) -> c
	{NewFunction(NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b')), hm.TypeVariable('a')), mSubs{'a': proton, 'b': neutron}, NewFunction(NewFunction(proton, neutron), proton)},
}

func TestFunctionType_FlatTypes(t *testing.T) {
	fnType := NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b'), hm.TypeVariable('c'))
	ts := fnType.FlatTypes()
	correct := hm.Types{hm.TypeVariable('a'), hm.TypeVariable('b'), hm.TypeVariable('c')}
	assert.Equal(t, ts, correct)

	fnType2 := NewFunction(fnType, hm.TypeVariable('d'))
	correct = append(correct, hm.TypeVariable('d'))
	ts = fnType2.FlatTypes()
	assert.Equal(t, ts, correct)
}

func TestFunctionType_Clone(t *testing.T) {
	fnType := NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b'), hm.TypeVariable('c'))
	assert.Equal(t, fnType.Clone(), fnType)

	rec := NewRecordType("", hm.TypeVariable('a'), NewFunction(hm.TypeVariable('a'), hm.TypeVariable('b')), hm.TypeVariable('c'))
	fnType = NewFunction(rec, rec)
	assert.Equal(t, fnType.Clone(), fnType)
}
