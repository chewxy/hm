package hm

import "testing"

var tvSetTests = []struct {
	op   string
	tvs0 TypeVarSet
	tvs1 TypeVarSet

	expected TypeVarSet
	ind      int
	eq       bool
}{
	{"set", TypeVarSet{'a', 'a', 'a'}, nil, TypeVarSet{'a'}, 0, false},
	{"set", TypeVarSet{'c', 'b', 'a'}, nil, TypeVarSet{'a', 'b', 'c'}, 0, false},
	{"intersect", TypeVarSet{'a', 'b', 'c'}, TypeVarSet{'d', 'e', 'f'}, TypeVarSet{}, -1, false},
	{"intersect", TypeVarSet{'a', 'b', 'c'}, TypeVarSet{'b', 'c', 'd'}, TypeVarSet{'b', 'c'}, -1, false},
	{"intersect", TypeVarSet{'a', 'b', 'c'}, nil, nil, -1, false},
	{"intersect", TypeVarSet{'a', 'b', 'c'}, TypeVarSet{'c', 'b', 'a'}, TypeVarSet{'a', 'b', 'c'}, 0, true},
	{"union", TypeVarSet{'a', 'b'}, TypeVarSet{'c', 'd'}, TypeVarSet{'a', 'b', 'c', 'd'}, 0, false},
	{"union", TypeVarSet{'a', 'c', 'b'}, TypeVarSet{'c', 'd'}, TypeVarSet{'a', 'b', 'c', 'd'}, 0, false},
	{"union", TypeVarSet{'a', 'b'}, nil, TypeVarSet{'a', 'b'}, 0, false},
	{"diff", TypeVarSet{'a', 'b', 'c'}, TypeVarSet{'d', 'e', 'c'}, TypeVarSet{'a', 'b'}, 0, false},
	{"diff", TypeVarSet{'a', 'b', 'c'}, TypeVarSet{'c', 'd', 'e'}, TypeVarSet{'a', 'b'}, 0, false},
	{"diff", TypeVarSet{'a', 'b', 'c'}, TypeVarSet{'d', 'e', 'f'}, TypeVarSet{'a', 'b', 'c'}, 0, false},
}

func TestTypeVarSet(t *testing.T) {
	for i, tst := range tvSetTests {
		var s TypeVarSet
		switch tst.op {
		case "set":
			s = tst.tvs0.Set()
			if !s.Equals(tst.expected) {
				t.Errorf("%s op (%d): expected: %v, got %v", tst.op, i, tst.expected, s)
			}
		case "intersect":
			s = tst.tvs0.Intersect(tst.tvs1)
			if !s.Equals(tst.expected) {
				t.Errorf("%s op (%d): expected: %v, got %v", tst.op, i, tst.expected, s)
			}
		case "union":
			s = tst.tvs0.Union(tst.tvs1)
			if !s.Equals(tst.expected) {
				t.Errorf("%s op (%d): expected: %v, got %v", tst.op, i, tst.expected, s)
			}
		case "diff":
			s = tst.tvs0.Difference(tst.tvs1)
			if !s.Equals(tst.expected) {
				t.Errorf("%s op (%d): expected: %v, got %v", tst.op, i, tst.expected, s)
			}
		}

		if ind := s.Index('a'); ind != tst.ind {
			t.Errorf("%s op %d index : expected %d got %v", tst.op, i, tst.ind, ind)
		}

		if eq := tst.tvs0.Equals(tst.tvs1); eq != tst.eq {
			t.Errorf("%s op %d eq: expected %t got %v", tst.op, i, tst.eq, eq)
		}
	}

	tvs := TypeVarSet{'a'}
	if !tvs.Equals(tvs) {
		t.Error("A set should be equal to itself")
	}

}
