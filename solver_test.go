package hm

import "testing"

var solverTest = []struct {
	cs Constraints

	expected Subs
	err      bool
}{
	{Constraints{{TypeVariable('a'), proton}}, mSubs{'a': proton}, false},
	{Constraints{{NewFnType(TypeVariable('a'), proton), neutron}}, nil, true},
	{Constraints{{NewFnType(TypeVariable('a'), proton), NewFnType(proton, proton)}}, mSubs{'a': proton}, false},

	{Constraints{
		{
			NewFnType(TypeVariable('a'), TypeVariable('a'), list{TypeVariable('a')}),
			NewFnType(proton, proton, TypeVariable('b')),
		},
	},
		mSubs{'a': proton, 'b': list{proton}}, false,
	},

	{
		Constraints{
			{TypeVariable('a'), TypeVariable('b')},
			{TypeVariable('a'), proton},
		},
		mSubs{'a': proton}, false,
	},

	{
		Constraints{
			{
				NewRecordType("", TypeVariable('a'), TypeVariable('a'), TypeVariable('b')),
				NewRecordType("", neutron, neutron, proton),
			},
		},
		mSubs{'a': neutron, 'b': proton}, false,
	},
}

func TestSolver(t *testing.T) {
	for i, sts := range solverTest {
		solver := newSolver()
		solver.solve(sts.cs)

		if sts.err {
			if solver.err == nil {
				t.Errorf("Test %d Expected an error", i)
			}
			continue
		} else if solver.err != nil {
			t.Error(solver.err)
		}

		for _, v := range sts.expected.Iter() {
			if T, ok := solver.sub.Get(v.Tv); !ok {
				t.Errorf("Test %d: Expected type variable %v in subs: %v", i, v.Tv, solver.sub)
				break
			} else if T != v.T {
				t.Errorf("Test %d: Expected replacement to be %v. Got %v instead", i, v.T, T)
			}
		}
	}
}
