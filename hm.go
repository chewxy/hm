package hm

import "github.com/pkg/errors"

const digits = "0123456789"

type Env interface {
	// TypeOf returns the type of the identifier
	TypeOf(id string) (Type, error)

	// Add adds the identifier and type
	Add(id string, t Type) Env

	// AddSpecified adds a TypeVariable to the set of specified type variables
	AddConcreteVar(tv TypeVariable) Env

	// Specified is a set of TypeVariables that have been specified
	ConcreteVars() Types

	// Clone clones the Env
	Clone() Env
}

type SimpleEnvConsOpt func(*SimpleEnv)

func WithDict(m map[string]Type) SimpleEnvConsOpt {
	f := func(env *SimpleEnv) {
		env.m = m
	}
	return f
}

type SimpleEnv struct {
	m map[string]Type
	s Types
}

func NewSimpleEnv(opts ...SimpleEnvConsOpt) *SimpleEnv {
	env := &SimpleEnv{
		m: make(map[string]Type),
	}

	for _, opt := range opts {
		opt(env)
	}

	return env
}

func (env *SimpleEnv) TypeOf(id string) (Type, error) {
	if t, ok := env.m[id]; ok {
		return env.Fresh(t), nil
	}

	return nil, errors.Errorf("Identifier %q not defined", id)
}

func (env *SimpleEnv) Add(id string, t Type) Env {
	env.m[id] = t
	return env
}

func (env *SimpleEnv) AddConcreteVar(tv TypeVariable) Env {
	env.s = env.s.Add(tv)
	return env
}

func (env *SimpleEnv) ConcreteVars() Types {
	return env.s
}

func (env *SimpleEnv) Clone() Env {
	m := make(map[string]Type)
	for k, v := range env.m {
		m[k] = v
	}

	s := make(Types, len(env.s))
	copy(s, env.s)

	return &SimpleEnv{
		m: m,
		s: s,
	}
}

func (env *SimpleEnv) Fresh(t Type) Type {
	// since TypeVariable cannot be a map key, we'll not use a map and use two slices to keep track of mapping instead
	var k, v Types
	retVal, _, _ := env.fresh(t, k, v)
	return retVal
}

// recursively creates a fresh type
func (env *SimpleEnv) fresh(t Type, k, v Types) (freshType Type, keys Types, values Types) {
	switch p := Prune(t).(type) {
	case TypeVariable:
		if env.s.Contains(p) {
			return p, k, v
		}

		var i int
		if i = k.Index(p); i > -1 {
			return v[i], k, v
		}

		tv := NewTypeVar(randomStr(5))
		k = append(k, p)
		v = append(v, tv)
		return tv, k, v
	case TypeConst:
		return p.Clone(), k, v
	case TypeOp:
		ts := make(Types, len(p.Types()))
		for i, tt := range ts {
			ts[i], k, v = env.fresh(tt, k, v)
		}
		top := p.Clone()
		top = top.SetTypes(ts...)
		return top, k, v
	default:
		panic("Not implemented yet")
	}
}

// The Infer function is the core of the HM type inference system. This is a reference implementation and is completely servicable, but not quite performant.
// You should use this as a reference and write your own infer function.
//
// Very briefly, these rules are implemented:
//
// Var
//
// If x is of type T, in a collection of statements Γ, then we can infer that x has type T when we come to a new instance of x
//		 x: T ∈ Γ
//		-----------
//		 Γ ⊢ x: T
//
// Apply
//
// If f is a function that takes T1 and returns T2; and if x is of type T1;
// then we can infer that the result of applying f on x will yield a result has type T2
//		 Γ ⊢ f: T1→T2  Γ ⊢ x: T1
//		-------------------------
//		     Γ ⊢ f(x): T2
//
//
// Lambda Abstraction
//
// If we assume x has type T1, and because of that we were able to infer e has type T2
// then we can infer that the lambda abstraction of e with respect to the variable x,  λx.e,
// will be a function with type T1→T2
//		  Γ, x: T1 ⊢ e: T2
//		-------------------
//		  Γ ⊢ λx.e: T1→T2
//
// Let
//
// If we can infer that e1 has type T1 and if we take x to have type T1 such that we could infer that e2 has type T2,
// then we can infer that the result of letting x = e1 and substituting it into e2 has type T2
//		  Γ, e1: T1  Γ, x: T1 ⊢ e2: T2
//		--------------------------------
//		     Γ ⊢ let x = e1 in e2: T2
//
// Instantiation
//
// If ...
// 		  Γ ⊢ e: T1  T1 ⊑ T
//		----------------------
//		       Γ ⊢ e: T
//
// Generalization
//
// If ...
//		  Γ ⊢ e: T1  T1 ∉ free(Γ)
//		---------------------------
//		   Γ ⊢ e: ∀ α.T1
func Infer(node Node, env Env) (retVal Type, err error) {
	var ok bool
	switch n := node.(type) {
	case Lit:
		// if the node knows its own type...
		var typer Typer
		if typer, ok = n.(Typer); ok {
			return typer.Type(), nil
		}

		return env.TypeOf(n.Name())
	case Var:
		if retVal, err = env.TypeOf(n.Name()); err != nil {
			// add to env
			env.Add(n.Name(), n.Type())
			return n.Type(), nil
		}
		return
	case Lambda:
		argType := NewTypeVar(randomStr(5))
		scope := env.Clone()
		scope = scope.Add(n.Name(), argType)
		scope = scope.AddConcreteVar(argType)

		if retVal, err = Infer(n.Body(), scope); err != nil {
			return
		}
		retVal = NewFnType(argType, retVal)
		return
	case Apply:
		var fnType Type
		if fnType, err = Infer(n.Fn(), env); err != nil {
			return
		}

		var arg Type
		if arg, err = Infer(n.Body(), env); err != nil {
			return
		}
		retType := NewTypeVar(randomStr(5))

		fn := NewFnType(arg, retType)
		var t0 Type
		if t0, _, err = Unify(fn, fnType); err != nil {
			return
		}

		fn = t0.(*FunctionType)
		retVal = fn.ReturnType()
	case LetRec:
		var tmp Type
		tmp = NewTypeVar(randomStr(5))
		scope := env.Clone()
		scope = scope.Add(n.Name(), tmp)
		scope = scope.AddConcreteVar(tmp.(TypeVariable))

		var def Type
		if def, err = Infer(n.Def(), scope); err != nil {
			return
		}

		if tmp, _, err = Unify(tmp, def); err != nil {
			return
		}

		scope = scope.Add(n.Name(), tmp)
		return Infer(n.Body(), scope)

	case Let:
		var def Type
		if def, err = Infer(n.Def(), env); err != nil {
			return
		}

		scope := env.Clone()
		scope = scope.Add(n.Name(), def)

		return Infer(n.Body(), scope)
	case UntypedNode:

	}
	return
}
