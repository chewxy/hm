package hm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpts(t *testing.T) {
	assert := assert.New(t)

	DontUseFnPool()
	assert.False(IsUsingFnPool())

	UseFnPool()
	assert.True(IsUsingFnPool())

	DontUseTsPool()
	assert.False(IsUsingTsPool())

	UseTsPool()
	assert.True(IsUsingTsPool())

	DontUseTvPool()
	assert.False(IsUsingTvPool())

	UseTvPool()
	assert.True(IsUsingTvPool())
}

func TestFnPool(t *testing.T) {
	assert := assert.New(t)

	fn := NewFnType(proton, proton, NewTypeVar("a"))
	ReturnFnType(fn)

	fn = borrowFnType()
	assert.Nil(fn.ts[0])
	assert.Nil(fn.ts[1])

	fn.ts[0] = NewFnType(NewTypeVar("a"), proton)
	fn.ts[1] = proton
	ReturnFnType(fn)

	fn = borrowFnType()
	assert.Nil(fn.ts[0])
	assert.Nil(fn.ts[1])
}

func TestTypes1Pool(t *testing.T) {
	assert := assert.New(t)

	ts := BorrowTypes1()
	ts[0] = proton

	ReturnTypes1(ts)
	ts = BorrowTypes1()
	assert.Nil(ts[0])
}

func TestTypeVarPool(t *testing.T) {
	assert := assert.New(t)

	tv := borrowTypeVar()
	tv.name = "hello"
	tv.instance = NewFnType(NewTypeVar("a"), NewTypeVar("a", WithInstance(NewTypeVar("b"))))
	ReturnTypeVar(tv)

	tv = borrowTypeVar()
	assert.Equal("", tv.name)
	assert.Nil(tv.instance)
	assert.Nil(tv.constraints)
}
