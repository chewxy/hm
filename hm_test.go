package hm

import "testing"

var unifyTests = []struct {
	name string
	a    Type
	b    Type

	subs Subs
	err  bool // does it error?
}{
	{"a ~ a (recursive unification)", TypeVariable('a'), TypeVariable('a'), nil, true},
	{"a ~ b", TypeVariable('a'), TypeVariable('b'), mSubs{'a': TypeVariable('b')}, false},
	{"a ~ proton", TypeVariable('a'), proton, mSubs{'a': proton}, false},
	{"proton ~ a", proton, TypeVariable('a'), mSubs{'a': proton}, false},

	// typeconst ~ typeconst
	{"proton ~ proton", proton, proton, nil, false},
	{"proton ~ neutron", proton, neutron, nil, true},
	{"List a ~ List proton", list{TypeVariable('a')}, list{proton}, mSubs{'a': proton}, false},

	// function types
	{"List a → List a ~ List proton → List proton",
		NewFnType(list{TypeVariable('a')}, list{TypeVariable('a')}),
		NewFnType(list{proton}, list{proton}),
		mSubs{'a': proton}, false},
	{"List proton → List proton ~ List a → List a",
		NewFnType(list{proton}, list{proton}),
		NewFnType(list{TypeVariable('a')}, list{TypeVariable('a')}),
		mSubs{'a': proton}, false},
	{"List a → a ~ List proton → proton",
		NewFnType(list{TypeVariable('a')}, TypeVariable('a')),
		NewFnType(list{proton}, proton),
		mSubs{'a': proton}, false},
	{"List proton → proton ~ List a → a ",
		NewFnType(list{proton}, proton),
		NewFnType(list{TypeVariable('a')}, TypeVariable('a')),
		mSubs{'a': proton}, false},
	{"List a → a → List a ~ List proton → proton → b",
		NewFnType(list{TypeVariable('a')}, TypeVariable('a'), list{TypeVariable('a')}),
		NewFnType(list{proton}, proton, TypeVariable('b')),
		mSubs{'a': proton, 'b': list{proton}}, false},
	{"(a, a, b) ~ (proton, proton, neutron)",
		NewRecordType("", TypeVariable('a'), TypeVariable('a'), TypeVariable('b')),
		NewRecordType("", proton, proton, neutron),
		mSubs{'a': proton, 'b': neutron}, false},
}

func TestUnify(t *testing.T) {
	// assert := assert.New(t)
	var t0, t1 Type
	var u0, u1 Type
	var sub Subs
	var err error

	for _, uts := range unifyTests {
		// logf("unifying %v", uts.name)
		t0 = uts.a
		t1 = uts.b
		sub, err = Unify(t0, t1)

		switch {
		case err == nil && uts.err:
			t.Errorf("Test %q - Expected an error: %v | u0: %#v, u1: %#v", uts.name, err, u0, u1)
		case err != nil && !uts.err:
			t.Errorf("Test %q errored: %v ", uts.name, err)
		}

		if uts.err {
			continue
		}

		if uts.subs == nil {
			if sub != nil {
				t.Errorf("Test: %q Expected no substitution. Got %v instead", uts.name, sub)
			}
			continue
		}

		for _, s := range uts.subs.Iter() {
			if T, ok := sub.Get(s.Tv); !ok {
				t.Errorf("Test: %q TypeVariable %v expected in result", uts.name, s.Tv)
			} else if T != s.T {
				t.Errorf("Test: %q Expected TypeVariable %v to be substituted by %v. Got %v instead", uts.name, s.Tv, s.T, T)
			}
		}

		if uts.subs.Size() != sub.Size() {
			t.Errorf("Test: %q Expected subs to be the same size", uts.name)
		}

		sub = nil
	}
}

var inferTests = []struct {
	name string

	expr       Expression
	correct    Type
	correctTVS TypeVarSet
	err        bool
}{
	{"Lit", lit("1"), Float, nil, false},
	{"Undefined Lit", lit("a"), nil, nil, true},
	{"App", app{lit("+"), lit("1")}, NewFnType(Float, Float), nil, false},

	{"Lambda", λ{"n", app{lit("+"), lit("1")}}, NewFnType(TypeVariable('a'), Float, Float), TypeVarSet{'a'}, false},
	{"Lambda (+1)", λ{"a", app{lit("+1"), lit("a")}}, NewFnType(TypeVariable('a'), TypeVariable('a')), TypeVarSet{'a'}, false},

	{"Var - found", variable("x"), proton, nil, false},
	{"Var - notfound", variable("y"), nil, nil, true},

	{"Self Infer - no err", selfInferer(true), proton, nil, false},
	{"Self Infer - err", selfInferer(false), nil, nil, true},

	{"nil expr", nil, nil, nil, true},
}

func TestInfer(t *testing.T) {
	env := SimpleEnv{
		"+":  &Scheme{tvs: TypeVarSet{'a'}, t: NewFnType(TypeVariable('a'), TypeVariable('a'), TypeVariable('a'))},
		"+1": &Scheme{tvs: TypeVarSet{'a'}, t: NewFnType(TypeVariable('a'), TypeVariable('a'))},
		"x":  NewScheme(nil, proton),
	}

	for _, its := range inferTests {
		sch, err := Infer(env, its.expr)

		if its.err {
			if err == nil {
				t.Errorf("Test %q : Expected error. %v", its.name, sch)
			}
			continue
		} else {
			if err != nil {
				t.Errorf("Test %q Error: %v", its.name, err)
			}
		}

		if !sch.t.Eq(its.correct) {
			t.Errorf("Test %q: Expected %v. Got %v", its.name, its.correct, sch.t)
		}

		for _, tv := range its.correctTVS {
			if !sch.tvs.Contains(tv) {
				t.Errorf("Test %q: Expected %v to be in the scheme.", its.name, tv)
				break
			}
		}

		if len(its.correctTVS) != len(sch.tvs) {
			t.Errorf("Test %q: Expected scheme to have %v. Got %v instead", its.name, its.correctTVS, sch.tvs)
		}
	}

	// test without env
	its := inferTests[0]
	sch, err := Infer(nil, its.expr)
	if err != nil {
		t.Errorf("Testing a nil Env. Shouldn't have errored. Got err:  %v", err)
	}
	if !sch.t.Eq(its.correct) {
		t.Errorf("Testing nil Env. Expected %v to be in the scheme. Got scheme %v instead", its.correct, sch)
	}

}
