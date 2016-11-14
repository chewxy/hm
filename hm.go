package hm

import "github.com/pkg/errors"

type Env interface {
	// TypeOf returns the type of the identifier
	TypeOf(id string) (Type, error)

	// Add adds the identifier and type
	Add(id string, t Type) Env

	// AddSpecified adds a TypeVariable to the set of specified type variables
	AddConcreteVar(tv TypeVariable) Env

	// AddReplacements adds a replacement
	AddReplacement(map[TypeVariable]Type) Env

	// Replacements returns the set of replacements
	Replacements() map[TypeVariable]Type

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

func WithConcreteVars(s Types) SimpleEnvConsOpt {
	f := func(env *SimpleEnv) {
		env.s = env.s.Union(s)
	}
	return f
}

type SimpleEnv struct {
	m map[string]Type
	s Types

	r map[TypeVariable]Type
}

func NewSimpleEnv(opts ...SimpleEnvConsOpt) *SimpleEnv {
	env := &SimpleEnv{
		m: make(map[string]Type),
		r: make(map[TypeVariable]Type),
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

func (env *SimpleEnv) AddReplacement(r map[TypeVariable]Type) Env {
	for k, v := range r {
		env.r[k] = v
	}
	return env
}

func (env *SimpleEnv) Replacements() map[TypeVariable]Type {
	return env.r
}

func (env *SimpleEnv) ConcreteVars() Types {
	return env.s
}

func (env *SimpleEnv) Clone() Env {
	m := make(map[string]Type)
	for k, v := range env.m {
		m[k] = v
	}

	r := make(map[TypeVariable]Type)
	for k, v := range env.r {
		r[k] = v
	}

	s := make(Types, len(env.s))
	copy(s, env.s)

	return &SimpleEnv{
		m: m,
		s: s,
		r: r,
	}
}

func (env *SimpleEnv) Fresh(t Type) Type {
	enterLoggingContext()
	defer leaveLoggingContext()
	// since TypeVariable cannot be a map key, we'll not use a map and use two slices to keep track of mapping instead
	retVal := env.fresh(t)
	return retVal
}

// recursively creates a fresh type
func (env *SimpleEnv) fresh(t Type) (freshType Type) {
	enterLoggingContext()
	defer leaveLoggingContext()

	switch p := Prune(t).(type) {
	case TypeVariable:
		if env.s.Contains(p) {
			return p
		}

		if tv, ok := env.r[p]; ok {
			return tv
		}
		tv := NewTypeVar(randomStr(5))
		env.r[p] = tv
		env.r[tv] = tv
		return tv
	case TypeConst:
		return p.Clone()
	case TypeOp:
		pts := p.Types()

		// ts := make(Types, len(pts))
		// for i, tt := range pts {
		// 	ts[i], k, v = env.fresh(tt, k, v)
		// }
		enterLoggingContext()
		for i := 0; i < len(pts); i++ {
			tt := pts[i]

			var tt2 Type
			tt2 = env.fresh(tt)

			if tv, ok := tt.(TypeVariable); ok {
				p = p.Replace(tv, tt2)
				pts = p.Types()
			}

		}
		leaveLoggingContext()

		// top := p.Clone()
		// top = top.SetTypes(ts...)
		// return top, k, v

		return p
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

	// if the node knows its own type...
	var typer Typer
	if typer, ok = node.(Typer); ok {
		// and if the type isn't nil...
		if retVal = typer.Type(); retVal != nil {
			return
		}
	}

	switch n := node.(type) {
	case Lit:
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

		if replacement, ok := scope.Replacements()[argType]; ok {
			retVal = NewFnType(replacement, retVal)
		} else {
			retVal = NewFnType(argType, retVal)
		}

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
		var r map[TypeVariable]Type
		if t0, _, r, err = Unify(fn, fnType); err != nil {
			return
		}

		env = env.AddReplacement(r)
		fn = t0.(*FunctionType)
		for k, v := range r {
			fn = fn.Replace(k, v).(*FunctionType)
		}

		retVal = fn.ts[1]
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

		var r map[TypeVariable]Type
		if tmp, _, r, err = Unify(tmp, def); err != nil {
			return
		}

		env = env.AddReplacement(r)
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
