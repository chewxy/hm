package hm

import "sync"

const poolSize = 4

var sSubPool = [poolSize]*sync.Pool{
	&sync.Pool{
		New: func() interface{} { return make(sSubs, 1) },
	},
	&sync.Pool{
		New: func() interface{} { return make(sSubs, 2) },
	},
	&sync.Pool{
		New: func() interface{} { return make(sSubs, 3) },
	},
	&sync.Pool{
		New: func() interface{} { return make(sSubs, 4) },
	},
}

var mSubPool = &sync.Pool{
	New: func() interface{} { return make(mSubs) },
}

func ReturnSubs(sub Subs) {
	switch s := sub.(type) {
	case mSubs:
		for k := range s {
			delete(s, k)
		}
		mSubPool.Put(sub)
	case sSubs:
		size := cap(s)
		s = s[:cap(s)]
		if size > 0 && size < poolSize+1 {
			// reset to empty
			for i := range s {
				s[i] = Substitution{}
			}

			sSubPool[size-1].Put(sub)
		}
	}
}

func BorrowMSubs() mSubs {
	return mSubPool.Get().(mSubs)
}

func BorrowSSubs(size int) sSubs {
	if size > 0 && size < 5 {
		retVal := sSubPool[size-1].Get().(sSubs)
		return retVal
	}
	return make(sSubs, size)
}

var typesPool = [poolSize]*sync.Pool{
	&sync.Pool{
		New: func() interface{} { return make(Types, 1) },
	},

	&sync.Pool{
		New: func() interface{} { return make(Types, 2) },
	},

	&sync.Pool{
		New: func() interface{} { return make(Types, 3) },
	},

	&sync.Pool{
		New: func() interface{} { return make(Types, 4) },
	},
}

func BorrowTypes(size int) Types {
	if size > 0 && size < poolSize+1 {
		return typesPool[size-1].Get().(Types)
	}
	return make(Types, size)
}

func ReturnTypes(ts Types) {
	if size := cap(ts); size > 0 && size < poolSize+1 {
		ts = ts[:cap(ts)]
		for i := range ts {
			ts[i] = nil
		}
		typesPool[size-1].Put(ts)
	}
}

var typeVarSetPool = [poolSize]*sync.Pool{
	&sync.Pool{
		New: func() interface{} { return make(TypeVarSet, 1) },
	},

	&sync.Pool{
		New: func() interface{} { return make(TypeVarSet, 2) },
	},

	&sync.Pool{
		New: func() interface{} { return make(TypeVarSet, 3) },
	},

	&sync.Pool{
		New: func() interface{} { return make(TypeVarSet, 4) },
	},
}

func BorrowTypeVarSet(size int) TypeVarSet {
	if size > 0 && size < poolSize+1 {
		return typeVarSetPool[size-1].Get().(TypeVarSet)
	}
	return make(TypeVarSet, size)
}

func ReturnTypeVarSet(ts TypeVarSet) {
	var def TypeVariable
	if size := cap(ts); size > 0 && size < poolSize+1 {
		ts = ts[:cap(ts)]
		for i := range ts {
			ts[i] = def
		}
		typeVarSetPool[size-1].Put(ts)
	}
}
