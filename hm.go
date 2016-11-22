package hm

import "github.com/pkg/errors"

type Inferer struct {
	e   Expression
	env Env
	cs  Constraints
	err error
	t   Type

	count int
}

func NewInferer(env Env) *Inferer {
	return &Inferer{
		env: env,
	}
}

func (infer *Inferer) fresh() TypeVariable {
	retVal := letters[infer.count]
	infer.count++
	return TypeVariable(retVal)
}

func (infer *Inferer) lookup(name string) {
	s, ok := infer.env.SchemeOf(name)
	if !ok {
		infer.err = errors.Errorf("Undefined %v", name)
		return
	}
	infer = instantiate(infer, s)
}

func (infer *Inferer) Infer(expr Expression) *Inferer {
	logf("Infering type of %v", expr)
	enterLoggingContext()
	defer leaveLoggingContext()

	if infer.err != nil {
		return infer
	}

	if et, ok := expr.(Typer); ok {
		if infer.t = et.Type(); infer.t != nil {
			return infer
		}
	}

	switch et := expr.(type) {
	case Literal:
		infer.lookup(et.Name())
	case Var:
		logf("Is Var")
		infer.lookup(et.Name())
		if infer.err != nil {
			infer.env.Add(et.Name(), &Scheme{t: et.Type()})
		}
	case Lambda:
		logf("%v Is Lambda", et)
		tv := infer.fresh()
		env := infer.env // backup

		logf("Cloning env")
		infer.env = infer.env.Clone()
		infer.env.Remove(et.Name())
		sc := new(Scheme)
		sc.t = tv
		infer.env.Add(et.Name(), sc)
		logf("cloned env : %v", infer.env)

		infer.Infer(et.Body())
		if infer.err != nil {
			return infer
		}
		infer.t = NewFnType(tv, infer.t)
		infer.env = env // restore backup
		logf("infer.t: %v", infer.t)
	case Apply:
		logf("%v Is Apply", et)
		infer.Infer(et.Fn())
		if infer.err != nil {
			return infer
		}

		fnType, fnCs := infer.t, infer.cs

		logf("fnType is %v", fnType)
		logf("fnCs %v", fnCs)
		logf("env %v", infer.env)

		infer.Infer(et.Body())
		bodyType, bodyCs := infer.t, infer.cs

		logf("Body Type is %v", bodyType)
		logf("bodyCs %v", bodyCs)

		tv := infer.fresh()
		cs := append(fnCs, bodyCs...)
		cs = append(cs, Constraint{fnType, NewFnType(bodyType, tv)})
		logf("cs %v", cs)

		infer.t = tv
		infer.cs = cs
	case LetRec:
		logf("Is LetRec")
		tv := infer.fresh()
		// env := infer.env // backup

		logf("Setting up env")
		infer.env = infer.env.Clone()
		infer.env.Remove(et.Name())
		infer.env.Add(et.Name(), &Scheme{tvs: TypeVarSet{tv}, t: tv})

		logf("Inferring def")
		infer.Infer(et.Def())
		if infer.err != nil {
			return infer
		}
		defType, defCs := infer.t, infer.cs

		logf("Solving def")
		s := newSolver()
		s.solve(defCs)
		if s.err != nil {
			infer.err = s.err
			return infer
		}

		logf("Generalizing")
		sc := generalize(infer.env.Apply(s.sub).(Env), defType.Apply(s.sub).(Type))

		infer.env.Remove(et.Name())
		infer.env.Add(et.Name(), sc)

		logf("Inferring body")
		infer.Infer(et.Body())
		logf("Done!")

		if infer.err != nil {
			return infer
		}

		logf("Applying sub %v to %v", s.sub, infer.t)
		infer.t = infer.t.Apply(s.sub).(Type)
		logf("Applying sub to constraints: %v", infer.cs)
		infer.cs = infer.cs.Apply(s.sub).(Constraints)
		logf("Putting together them constraints")
		infer.cs = append(infer.cs, defCs...)

	case Let:
		logf("Is Let")
		env := infer.env

		logf("Inferring def")
		infer.Infer(et.Def())
		defType, defCs := infer.t, infer.cs
		logf("defType %v", defType)

		s := newSolver()
		s.solve(defCs)
		if s.err != nil {
			infer.err = s.err
			return infer
		}

		logf("Generalizing %v within %v", defType, env)
		sc := generalize(env.Apply(s.sub).(Env), defType.Apply(s.sub).(Type))
		logf("Generalized: %v", sc)
		infer.env = infer.env.Clone()
		infer.env.Remove(et.Name())
		infer.env.Add(et.Name(), sc)

		logf("env %v", infer.env)
		logf("Inferring body")
		infer.Infer(et.Body())
		if infer.err != nil {
			return infer
		}

		infer.t = infer.t.Apply(s.sub).(Type)
		infer.cs = infer.cs.Apply(s.sub).(Constraints)
		infer.cs = append(infer.cs, defCs...)

	}

	return infer
}

func instantiate(infer *Inferer, s *Scheme) *Inferer {
	l := len(s.tvs)
	tvs := make(TypeVarSet, l)

	var sub Subs
	if l > 30 {
		sub = make(mSubs)
	} else {
		sub = newSliceSubs(l)
	}

	for i, tv := range s.tvs {
		f := infer.fresh()
		tvs[i] = f
		sub = sub.Add(tv, f)
	}

	infer.t = s.t.Apply(sub).(Type)
	return infer
}

func generalize(env Env, t Type) *Scheme {
	logf("generalizing %v over %v", t, env)
	enterLoggingContext()
	defer leaveLoggingContext()
	var envFree, tFree, diff TypeVarSet

	if env != nil {
		envFree = env.FreeTypeVar()
	}

	tFree = t.FreeTypeVar()

	switch {
	case envFree == nil && tFree == nil:
		goto ret
	case len(envFree) > 0 && len(tFree) > 0:
		defer ReturnTypeVarSet(envFree)
		defer ReturnTypeVarSet(tFree)
	case len(envFree) > 0 && len(tFree) == 0:
		// cannot return envFree because envFree will just be sorted and set
	case len(envFree) == 0 && len(tFree) > 0:
		// return ?
	}
	logf("tFree: %v, envFree %v", tFree, envFree)

	diff = tFree.Difference(envFree)

ret:
	return &Scheme{
		tvs: diff,
		t:   t,
	}
}

func Infer(env Env, expr Expression) (*Scheme, error) {
	logf("Infer")
	enterLoggingContext()
	defer leaveLoggingContext()
	logf("Infering with env")
	infer := NewInferer(env)
	infer.Infer(expr)
	if infer.err != nil {
		return nil, infer.err
	}
	logf("Solving...")
	logf("%v", infer.cs)
	s := newSolver()
	s.solve(infer.cs)

	if s.err != nil {
		return nil, s.err
	}
	logf("infer.t %v", infer.t)
	logf("infer.cs: %v", infer.cs)
	logf("s.sub %v", s.sub)
	t := infer.t.Apply(s.sub).(Type)
	return closeOver(t)
}

func Unify(a, b Type) (sub Subs, err error) {
	logf("%v ~ %v", a, b)
	enterLoggingContext()
	defer leaveLoggingContext()

	switch at := a.(type) {
	case TypeVariable:
		return bind(at, b)
	default:
		if a.Eq(b) {
			return nil, nil
		}

		if btv, ok := b.(TypeVariable); ok {
			return bind(btv, a)
		}
		atypes := a.Types()
		btypes := b.Types()
		defer ReturnTypes(atypes)
		defer ReturnTypes(btypes)

		if len(atypes) == 0 && len(btypes) == 0 {
			goto e
		}

		return unifyMany(atypes, btypes)

	e:
	}
	err = errors.Errorf("Unification Fail: %v ~ %v cannot be unified", a, b)
	return
}

func unifyMany(a, b Types) (sub Subs, err error) {
	logf("UnifyMany %v %v", a, b)
	enterLoggingContext()
	defer leaveLoggingContext()

	if len(a) != len(b) {
		return nil, errors.Errorf("Unequal length. a: %v b %v", a, b)
	}

	for i, at := range a {
		bt := b[i]

		if sub != nil {
			at = at.Apply(sub).(Type)
			bt = bt.Apply(sub).(Type)
		}

		var s2 Subs
		if s2, err = Unify(at, bt); err != nil {
			return nil, err
		}

		if sub == nil {
			sub = s2
		} else {
			sub2 := compose(sub, s2)
			defer ReturnSubs(s2)
			if sub2 != sub {
				defer ReturnSubs(sub)
			}
			sub = sub2
		}
	}
	return
}

func bind(tv TypeVariable, t Type) (sub Subs, err error) {
	logf("Binding %v to %v", tv, t)
	switch {
	// case tv == t:
	case occurs(tv, t):
		err = errors.Errorf("recursive unification")
	default:
		ssub := BorrowSSubs(1)
		ssub.s[0] = Substitution{tv, t}
		sub = ssub
	}
	logf("Sub %v", sub)
	return
}

func occurs(tv TypeVariable, s Substitutable) bool {
	ftv := s.FreeTypeVar()
	defer ReturnTypeVarSet(ftv)

	return ftv.Contains(tv)
}

func closeOver(t Type) (sch *Scheme, err error) {
	sch = generalize(nil, t)
	err = sch.normalize()
	logf("closeoversch: %v", sch)
	return
}
