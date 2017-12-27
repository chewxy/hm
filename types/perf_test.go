package hmtypes

import "testing"

func TestFnTypePool(t *testing.T) {
	f := borrowFn()
	f.A = NewFunction(proton, electron)
	f.B = NewFunction(proton, neutron)

	ReturnFn(f)
	f = borrowFn()
	if f.A != nil {
		t.Error("FunctionType not cleaned up: a is not nil")
	}
	if f.B != nil {
		t.Error("FunctionType not cleaned up: b is not nil")
	}

}
