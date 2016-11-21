package hm

type solver struct {
	sub Subs
	err error
}

func newSolver() *solver {
	return new(solver)
}

func (s *solver) solve(cs Constraints) {
	logf("solving constraints")
	if s.err != nil {
		return
	}

	switch len(cs) {
	case 0:
		return
	case 1:
		logf("len1")
		c := cs[0]
		s.sub, s.err = Unify(c.a, c.b)
		logf("s.sub: %v", s.sub)
		logf("s.err %v", s.err)
	default:
		var sub Subs
		c := cs[0]
		s.sub, s.err = Unify(c.a, c.b)
		defer ReturnSubs(s.sub)

		s.sub = compose(sub, s.sub)
		s.solve(cs[1:])

	}

	return
}
