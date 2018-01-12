package hmtypes

import (
	"unsafe"

	"github.com/chewxy/hm"
)

func borrowFn() *Function {
	got := hm.BorrowPair()
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
	p := (*hm.Pair)(unsafe.Pointer(fnt))
	hm.ReturnPair(p)
}
