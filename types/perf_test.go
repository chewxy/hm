package hmtypes

import "testing"

func TestFnTypePool(t *testing.T) {
	f := borrowFnType()
	f.A = NewFnType(proton, electron)
	f.B = NewFnType(proton, neutron)

	ReturnFnType(f)
	f = borrowFnType()
	if f.A != nil {
		t.Error("FunctionType not cleaned up: a is not nil")
	}
	if f.B != nil {
		t.Error("FunctionType not cleaned up: b is not nil")
	}

}
