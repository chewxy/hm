package hm

type Subs interface {
	Get(TypeVariable) (Type, bool)
	Add(TypeVariable, Type) Subs
	Remove(TypeVariable) Subs

	Iter() <-chan Substitution
	Size() int
	Clone() Subs
}

type Substitution struct {
	Tv TypeVariable
	T  Type
}

type sSubs []Substitution

func newSliceSubs(maybeSize ...int) sSubs {
	var size int
	if len(maybeSize) > 0 && maybeSize[0] > 0 {
		size = maybeSize[0]
	}
	retVal := make(sSubs, size)
	retVal = retVal[:0]
	return retVal
}

func (s sSubs) Get(tv TypeVariable) (Type, bool) {
	if i := s.index(tv); i >= 0 {
		return s[i].T, true
	}
	return nil, false
}

func (s sSubs) Add(tv TypeVariable, t Type) Subs {
	if i := s.index(tv); i >= 0 {
		s[i].T = t
		return s
	}
	s = append(s, Substitution{tv, t})
	return s
}

func (s sSubs) Remove(tv TypeVariable) Subs {
	if i := s.index(tv); i >= 0 {
		// for now we keep the order
		copy(s[i:], s[i+1:])
		s[len(s)-1].T = nil
		s = s[:len(s)-1]
	}

	return s
}

func (s sSubs) Iter() <-chan Substitution {
	ch := make(chan Substitution)

	go func() {
		for _, v := range s {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

func (s sSubs) Size() int { return len(s) }
func (s sSubs) Clone() Subs {
	retVal := make(sSubs, len(s))
	for i, v := range s {
		retVal[i] = v
	}
	return retVal
}

func (s sSubs) index(tv TypeVariable) int {
	for i, sub := range s {
		if sub.Tv == tv {
			return i
		}
	}
	return -1
}

type mSubs map[TypeVariable]Type

func (s mSubs) Get(tv TypeVariable) (Type, bool) { retVal, ok := s[tv]; return retVal, ok }
func (s mSubs) Add(tv TypeVariable, t Type) Subs { s[tv] = t; return s }
func (s mSubs) Remove(tv TypeVariable) Subs      { delete(s, tv); return s }

func (s mSubs) Iter() <-chan Substitution {
	ch := make(chan Substitution)
	go func() {
		for k, v := range s {
			ch <- Substitution{k, v}
		}
		close(ch)
	}()
	return ch
}

func (s mSubs) Size() int { return len(s) }
func (s mSubs) Clone() Subs {
	retVal := make(mSubs)
	for k, v := range s {
		retVal[k] = v
	}
	return retVal
}

func compose(a, b Subs) (retVal Subs) {
	if b == nil {
		return a
	}

	retVal = b.Clone()

	if a == nil {
		return
	}

	for v := range a.Iter() {
		retVal = retVal.Add(v.Tv, v.T)
	}

	logf("retVal: %v", retVal)
	enterLoggingContext()
	defer leaveLoggingContext()

	for v := range retVal.Iter() {
		retVal = retVal.Add(v.Tv, v.T.Apply(a).(Type))
	}
	logf("eh.. returning %v", retVal)
	return retVal
}
