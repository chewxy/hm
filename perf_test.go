package hm

import "testing"

func TestSubsPool(t *testing.T) {
	var def TypeVariable
	for i := 0; i < poolSize; i++ {
		s := BorrowSSubs(i + 1)
		if cap(s) != i+1 {
			t.Errorf("Expected s to have cap of %d", i+1)
			goto mSubTest
		}

		s[0] = Substitution{TypeVariable('a'), electron}
		ReturnSubs(s)
		s = BorrowSSubs(i + 1)

		for _, subst := range s {
			if subst.T != nil {
				t.Errorf("sSubsPool %d error: not clean: %v", i, subst)
				break
			}

			if subst.Tv != def {
				t.Errorf("sSubsPool %d error: not clean: %v", i, subst)
				break
			}
		}

	mSubTest:
		m := BorrowMSubs()
		if len(m) != 0 {
			t.Errorf("Expected borrowed mSubs to have 0 length")
		}

		m['a'] = electron
		ReturnSubs(m)

		m = BorrowMSubs()
		if len(m) != 0 {
			t.Errorf("Expected borrowed mSubs to have 0 length")
		}

	}

	// oob tests
	s := BorrowSSubs(10)
	if cap(s) != 10 {
		t.Error("Expected a cap of 10")
	}
	ReturnSubs(s)
}

func TestTypesPool(t *testing.T) {
	for i := 0; i < poolSize; i++ {
		ts := BorrowTypes(i + 1)
		if cap(ts) != i+1 {
			t.Errorf("Expected ts to have a cap of %v", i+1)
		}

		ts[0] = proton
		ReturnTypes(ts)
		ts = BorrowTypes(i + 1)
		for _, v := range ts {
			if v != nil {
				t.Errorf("Expected reshly borrowed Types to be nil")
			}
		}
	}

	// oob
	ts := BorrowTypes(10)
	if cap(ts) != 10 {
		t.Errorf("Expected a cap to 10")
	}
}

func TestTypeVarSetPool(t *testing.T) {
	var def TypeVariable
	for i := 0; i < poolSize; i++ {
		ts := BorrowTypeVarSet(i + 1)
		if cap(ts) != i+1 {
			t.Errorf("Expected ts to have a cap of %v", i+1)
		}

		ts[0] = 'z'
		ReturnTypeVarSet(ts)
		ts = BorrowTypeVarSet(i + 1)
		for _, v := range ts {
			if v != def {
				t.Errorf("Expected reshly borrowed Types to be def")
			}
		}
	}

	// oob
	tvs := BorrowTypeVarSet(10)
	if cap(tvs) != 10 {
		t.Error("Expected a cap of 10")
	}
}
