package hm

import (
	"fmt"
	"testing"
)

func TestConstraints(t *testing.T) {
	cs := Constraints{
		{TypeVariable('a'), proton},
		{TypeVariable('b'), proton},
	}
	correct := TypeVarSet{'a', 'b'}

	ftv := cs.FreeTypeVar()
	for _, v := range correct {
		if !ftv.Contains(v) {
			t.Errorf("Expected free type vars to contain %v", v)
			break
		}
	}

	sub := mSubs{
		'a': neutron,
	}

	cs = cs.Apply(sub).(Constraints)
	if cs[0].a != neutron {
		t.Error("Expected neutron")
	}
	if cs[0].b != proton {
		t.Error("Expected proton")
	}

	if cs[1].a != TypeVariable('b') {
		t.Error("There was nothing to substitute b with")
	}
	if cs[1].b != proton {
		t.Error("Expected proton")
	}

	if fmt.Sprintf("%v", cs) != "Constraints[{neutron = proton}, {b = proton}]" {
		t.Errorf("Error in formatting cs")
	}

}

func TestTypes_Contains(t *testing.T) {
	ts := Types{TypeVariable('a'), proton}

	if !ts.Contains(TypeVariable('a')) {
		t.Error("Expected ts to contain 'a'")
	}

	if !ts.Contains(proton) {
		t.Error("Expected ts to contain proton")
	}

	if ts.Contains(neutron) {
		t.Error("ts shouldn't contain neutron")
	}
}
