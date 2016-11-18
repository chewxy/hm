package hm

import (
	"fmt"
	"testing"
)

func TestTypeVariableBasics(t *testing.T) {
	tv := TypeVariable('a')
	if name := tv.Name(); name != "a" {
		t.Errorf("Expected name to be \"a\". Got %q instead", name)
	}

	if str := tv.String(); str != "a" {
		t.Errorf("Expected String() of 'a'. Got %q instead", str)
	}

	if tv.Types() != nil {
		t.Errorf("Expected Types() of TypeVariable to be nil")
	}

	ftv := tv.FreeTypeVar()
	if len(ftv) != 1 {
		t.Errorf("Expected a type variable to be free when FreeTypeVar() is called")
	}

	if ftv[0] != tv {
		t.Errorf("Expected ...")
	}

	sub := mSubs{
		'a': proton,
	}

	if tv.Apply(sub) != proton {
		t.Error("Expected proton")
	}

	sub = mSubs{
		'b': proton,
	}

	if tv.Apply(sub) != tv {
		t.Error("Expected unchanged")
	}
}

func TestTypeVariableNormalize(t *testing.T) {
	original := TypeVarSet{'c', 'a', 'd'}
	normalized := TypeVarSet{'a', 'b', 'c'}

	tv := TypeVariable('a')
	norm, err := tv.Normalize(original, normalized)
	if err != nil {
		t.Error(err)
	}

	if norm != TypeVariable('b') {
		t.Errorf("Expected 'b'. Got %v", norm)
	}

	tv = TypeVariable('e')
	if _, err = tv.Normalize(original, normalized); err == nil {
		t.Error("Expected an error")
	}
}

func TestTypeConst(t *testing.T) {
	T := proton
	if T.Name() != "proton" {
		t.Error("Expected name to be proton")
	}

	if fmt.Sprintf("%v", T) != "proton" {
		t.Error("Expected name to be proton")
	}

	if T.String() != "proton" {
		t.Error("Expected name to be proton")
	}

	if T2, err := T.Normalize(nil, nil); err != nil {
		t.Error(err)
	} else if T2 != T {
		t.Error("Const types should return itself")
	}
}
