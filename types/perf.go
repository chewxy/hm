package hmtypes

import (
	"sync"
	"unsafe"
)

var pairPool = &sync.Pool{
	New: func() interface{} { return new(Pair) },
}

func borrowPair() *Pair {
	return pairPool.Get().(*Pair)
}

func borrowFn() *Function {
	got := pairPool.Get().(*Pair)
	return (*Function)(unsafe.Pointer(got))
}

// ReturnFn returns a *FunctionType to the pool. NewFnType automatically borrows from the pool. USE WITH CAUTION
func ReturnFn(fnt *Function) {
	if a, ok := fnt.A.(*Function); ok {
		ReturnFn(a)
	}

	if b, ok := fnt.B.(*Function); ok {
		ReturnFn(b)
	}

	fnt.A = nil
	fnt.B = nil
	p := (*Pair)(unsafe.Pointer(fnt))
	pairPool.Put(p)
}
