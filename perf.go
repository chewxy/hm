package hm

import "sync"

var fntypePool = new(sync.Pool)

func borrowFnType() *FunctionType {
	return fntypePool.Get().(*FunctionType)
}

func returnFnType(t *FunctionType) {
	switch t0t := t.ts[0].(type) {
	case *FunctionType:
		returnFnType(t0t)
	}

	switch t1t := t.ts[1].(type) {
	case *FunctionType:
		returnFnType(t1t)
	}

	t.ts[0] = nil
	t.ts[1] = nil

	fntypePool.Put(t)
}

func init() {
	fntypePool.New = func() interface{} { return new(FunctionType) }
}
