package hmtypes

import (
	"sync"
	"unsafe"
)

var pairPool = &sync.Pool{
	New: func() interface{} { return new(Pair) },
}

func borrowFnType() *FunctionType {
	got := pairPool.Get().(*Pair)
	return (*FunctionType)(unsafe.Pointer(got))
}

// ReturnFnType returns a *FunctionType to the pool. NewFnType automatically borrows from the pool. USE WITH CAUTION
func ReturnFnType(fnt *FunctionType) {
	if a, ok := fnt.A.(*FunctionType); ok {
		ReturnFnType(a)
	}

	if b, ok := fnt.B.(*FunctionType); ok {
		ReturnFnType(b)
	}

	fnt.A = nil
	fnt.B = nil
	p := (*Pair)(unsafe.Pointer(fnt))
	pairPool.Put(p)
}
