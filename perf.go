package hm

import "sync"

var (
	fnPool = true
	tsPool = true
	tvPool = true

	fnPoolGuard = new(sync.Mutex)
	tsPoolGuard = new(sync.Mutex)
	tvPoolGuard = new(sync.Mutex)
)

func init() {

}

// DontUseFnPool ensures that the pool won't be used
func DontUseFnPool() { fnPoolGuard.Lock(); fnPool = false; fnPoolGuard.Unlock() }

// UseFnPool ensures that the pool for *FunctionType will be used
func UseFnPool() { fnPoolGuard.Lock(); fnPool = true; fnPoolGuard.Unlock() }

// IsUsingFnPool returns whether the *FunctionType pool is being used
func IsUsingFnPool() bool { return fnPool }

// DontUseTsPool ensures that the Types pool for sizes 1 won't be used
func DontUseTsPool() { tsPoolGuard.Lock(); tsPool = false; tsPoolGuard.Unlock() }

// UseFnPool ensures that the pool for Types with size 1 will be used
func UseTsPool() { tsPoolGuard.Lock(); tsPool = true; tsPoolGuard.Unlock() }

// IsUsingTsPool returns whether the pool for Types with size 1 is being used
func IsUsingTsPool() bool { return tsPool }

// DontUseTvPool ensures that the pool for *TypeVariable won't be used
func DontUseTvPool() { tvPoolGuard.Lock(); tvPool = false; tvPoolGuard.Unlock() }

// UseFnPool ensures that the pool for *TypeVariable will be used
func UseTvPool() { tvPoolGuard.Lock(); tvPool = true; tvPoolGuard.Unlock() }

// IsUsingTvPool returns whether the pool for *TypeVariable is being used
func IsUsingTvPool() bool { return tvPool }

var fntypePool = &sync.Pool{
	New: func() interface{} { return new(FunctionType) },
}

func borrowFnType() *FunctionType {
	return fntypePool.Get().(*FunctionType)
}

func ReturnFnType(t *FunctionType) {
	logf("returning FnType")
	enterLoggingContext()
	defer leaveLoggingContext()

	switch t0t := t.ts[0].(type) {
	case *FunctionType:
		ReturnFnType(t0t)
	case *TypeVariable:
		logf("Going to return t0t %p", t0t)
		ReturnTypeVar(t0t)
	}

	switch t1t := t.ts[1].(type) {
	case *FunctionType:
		ReturnFnType(t1t)
	case *TypeVariable:
		logf("Going to return t1t %p", t1t)
		ReturnTypeVar(t1t)
	}

	t.ts[0] = nil
	t.ts[1] = nil

	fntypePool.Put(t)
}

// pool for Types with size of 1
var types1Pool = &sync.Pool{
	New: func() interface{} { return make(Types, 1, 1) },
}

func BorrowTypes1() Types {
	return types1Pool.Get().(Types)
}

func ReturnTypes1(ts Types) {
	ts[0] = nil
	types1Pool.Put(ts)
}

// pool for typevar
// we also keep track of the used TypeVariables

var typeVarPool = &sync.Pool{
	New: func() interface{} { return new(TypeVariable) },
}
var typeVarLock = new(sync.Mutex)
var usedTypeVars = make(map[*TypeVariable]struct{})

func borrowTypeVar() *TypeVariable {
	typeVarLock.Lock()
	tv := typeVarPool.Get().(*TypeVariable)
	usedTypeVars[tv] = struct{}{}
	logf("borrowing tv %p %v", tv, tv)
	typeVarLock.Unlock()
	return tv
}

func ReturnTypeVar(tv *TypeVariable) {
	logf("returning tv %p %v", tv, tv)
	enterLoggingContext()
	defer leaveLoggingContext()

	typeVarLock.Lock()

	if _, ok := usedTypeVars[tv]; !ok {
		typeVarLock.Unlock()
		return
	}
	delete(usedTypeVars, tv)
	typeVarLock.Unlock()

	switch tit := tv.instance.(type) {
	case *TypeVariable:
		ReturnTypeVar(tit)
	case *FunctionType:
		ReturnFnType(tit)
	}
	tv.name = ""
	tv.instance = nil
	typeVarPool.Put(tv)

}

// handles Returning of Values
