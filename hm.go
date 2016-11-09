package hm

import "github.com/pkg/errors"

type Env interface {
	TypeOf(id string) (Type, error)
	Add(id string, t Type) Env

	Clone() Env
}

type SimpleEnv map[string]Type

func (env SimpleEnv) TypeOf(id string) (Type, error) {
	if t, ok := env[id]; ok {
		return t, nil
	}
	return nil, errors.Errorf("Identifier %q not defined", id)
}

func (env SimpleEnv) Add(id string, t Type) Env { env[id] = t; return env }
func (env SimpleEnv) Clone() Env {
	retVal := make(SimpleEnv)
	for k, v := range env {
		env[k] = v
	}
	return retVal
}

// Expr represents an expression.
type Node interface {
	Children() []Node
}

type UntypedNode interface {
	// ???
}

// Typer is any type that can report its own Type
type Typer interface {
	Type() Type
}

// Value is a node that represents a value
type Value interface {
	Node
	Typer
}

// Ident is a node that represents an identifier
type Ident interface {
	Node
	Name() string
}

// Var is a node that represents `var x int`
type Var interface {
	Node
	Typer
	Name() string
}

// Lambda is a node that represents a function definition
type Lambda interface {
	Node
	Arg() Node
	Body() Node
}

// Apply is a node that represents a function call/application
type Apply interface {
	Node
	Typer
	ArgTypes() Types
}

type Let interface {
	Node
}

type LetRec interface {
	Let
	IsLetRec() bool
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
	case Ident:
		// if the node knows its own type...
		var typer Typer
		if typer, ok = n.(Typer); ok {
			return typer.Type(), nil
		}

		return env.TypeOf(n.Name())
	case Lambda:
	case Apply:
	case LetRec:
	case Let:
	case UntypedNode:

	}
	return nil
}
