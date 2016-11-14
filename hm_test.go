package hm

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSimpleEnv_Fresh(t *testing.T) {
	assert := assert.New(t)
	var fresh Type
	var t0 Type
	env := NewSimpleEnv()

	fresh = env.Fresh(proton)
	assert.Equal(proton, fresh)

	t0 = NewTypeVar("a")
	fresh = env.Fresh(t0)
	assert.NotEqual(t0.Name(), fresh.Name())

	m := map[string]Type{
		"quarks":   list{quark},
		"electron": electron,
	}
	concrete := Types{
		NewTypeVar("a"),
	}
	env = NewSimpleEnv(WithDict(m), WithConcreteVars(concrete))

	t0 = NewTypeVar("a")
	assert.Equal(t0, env.Fresh(t0))
}

// to test infer, we'll just borrow the definitions from the greenspun example
var infer1 []struct {
	name string

	node    Node
	correct Type
	err     bool
}

func init() {
	infer1 = []struct {
		name string

		node    Node
		correct Type
		err     bool
	}{
		{"Lit", lit("1"), Float, false},
		{"Undefined Lit", lit("a"), nil, true},
		{"App", app{lit("+"), lit("1")}, NewFnType(Float, Float), false},

		// have to write helper functions to test these:
		{"Lambda", λ{"a", app{lit("+"), lit("1")}}, NewFnType(NewTypeVar("∀"), Float, Float), false},
		{"Lambda (+1)", λ{"a", app{lit("+1"), lit("a")}}, NewFnType(NewTypeVar("∀"), NewTypeVar("∀")), false},
	}
}

func TestInfer(t *testing.T) {
	assert := assert.New(t)
	var t0 Type
	var err error

	m := map[string]Type{
		"+":  NewFnType(NewTypeVar("a"), NewFnType(NewTypeVar("a"), NewTypeVar("a"))),
		"+1": NewFnType(NewTypeVar("a"), NewTypeVar("a")),
	}
	for _, its := range infer1 {
		env := NewSimpleEnv(WithDict(m))
		if t0, err = Infer(its.node, env); (its.err && err == nil) || (!its.err && err != nil) {
			if its.err {
				t.Errorf("Test %q: Expected an error", its.name)
			} else {
				t.Errorf("Test %q Err: %v", its.name, errors.Cause(err))
			}
			continue
		}

		if its.err {
			continue
		}

		assert.True(typeEqAnyVar(its.correct, t0), "Test : %v Correct: %#v | Got %#v", its.name, its.correct, t0)

		// assert.True(its.correct.Eq(t0), "Test : %v Correct: %#v | Got %#v", its.name, its.correct, t0)
	}
}
