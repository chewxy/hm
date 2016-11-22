package hm

import (
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

const digits = "0123456789"

type TyperExpression interface {
	Expression
	Typer
}

type λ struct {
	name string
	body Expression
}

func (n λ) Name() string     { return n.name }
func (n λ) Body() Expression { return n.body }
func (n λ) IsLambda() bool   { return true }

type lit string

func (n lit) Name() string     { return string(n) }
func (n lit) Body() Expression { return n }
func (n lit) Type() Type {
	switch {
	case strings.ContainsAny(digits, string(n)) && strings.ContainsAny(digits, string(n[0])):
		return Float
	case string(n) == "true" || string(n) == "false":
		return Bool
	default:
		return nil
	}
}
func (n lit) IsLit() bool    { return true }
func (n lit) IsLambda() bool { return true }

type app struct {
	f   Expression
	arg Expression
}

func (n app) Fn() Expression   { return n.f }
func (n app) Body() Expression { return n.arg }
func (n app) Arg() Expression  { return n.arg }

type let struct {
	name string
	def  Expression
	in   Expression
}

func (n let) Name() string     { return n.name }
func (n let) Def() Expression  { return n.def }
func (n let) Body() Expression { return n.in }

type letrec struct {
	name string
	def  Expression
	in   Expression
}

func (n letrec) Name() string           { return n.name }
func (n letrec) Def() Expression        { return n.def }
func (n letrec) Body() Expression       { return n.in }
func (n letrec) Children() []Expression { return []Expression{n.def, n.in} }
func (n letrec) IsRecursive() bool      { return true }

type prim byte

const (
	Float prim = iota
	Bool
)

// implement Type
func (t prim) Name() string                                   { return t.String() }
func (t prim) Apply(Subs) Substitutable                       { return t }
func (t prim) FreeTypeVar() TypeVarSet                        { return nil }
func (t prim) Normalize(TypeVarSet, TypeVarSet) (Type, error) { return t, nil }
func (t prim) Types() Types                                   { return nil }
func (t prim) Eq(other Type) bool {
	if ot, ok := other.(prim); ok {
		return ot == t
	}
	return false
}

func (t prim) Format(s fmt.State, c rune) { fmt.Fprintf(s, t.String()) }
func (t prim) String() string {
	switch t {
	case Float:
		return "Float"
	case Bool:
		return "Bool"
	}
	return "HELP"
}

//Phillip Greenspun's tenth law says:
//		"Any sufficiently complicated C or Fortran program contains an ad hoc, informally-specified, bug-ridden, slow implementation of half of Common Lisp."
//
// So let's implement a half-arsed lisp (Or rather, an AST that can optionally be executed upon if you write the correct interpreter)!
func Example_greenspun() {
	// haskell envy in a greenspun's tenth law example function!
	//
	// We'll assume the following is the "input" code
	// 		let fac n = if n == 0 then 1 else n * fac (n - 1) in fac 5
	// and what we have is the AST

	fac := letrec{
		"fac",
		λ{
			"n",
			app{
				app{
					app{
						lit("if"),
						app{
							lit("isZero"),
							lit("n"),
						},
					},
					lit("1"),
				},
				app{
					app{lit("mul"), lit("n")},
					app{lit("fac"), app{lit("--"), lit("n")}},
				},
			},
		},
		app{lit("fac"), lit("5")},
	}

	// but first, let's start with something simple:
	// let x = 3 in x+5
	simple := let{
		"x",
		lit("3"),
		app{
			app{
				lit("+"),
				lit("5"),
			},
			lit("x"),
		},
	}

	env := SimpleEnv{
		"--":     &Scheme{tvs: TypeVarSet{'a'}, t: NewFnType(TypeVariable('a'), TypeVariable('a'))},
		"if":     &Scheme{tvs: TypeVarSet{'a'}, t: NewFnType(Bool, TypeVariable('a'), TypeVariable('a'), TypeVariable('a'))},
		"isZero": &Scheme{t: NewFnType(Float, Bool)},
		"mul":    &Scheme{t: NewFnType(Float, Float, Float)},
		"+":      &Scheme{tvs: TypeVarSet{'a'}, t: NewFnType(TypeVariable('a'), TypeVariable('a'), TypeVariable('a'))},
	}

	var scheme *Scheme
	var err error
	scheme, err = Infer(env, simple)
	if err != nil {
		log.Printf("%+v", errors.Cause(err))
	}
	simpleType, ok := scheme.Type()
	fmt.Printf("simple Type: %v | isMonoType: %v | err: %v\n", simpleType, ok, err)

	scheme, err = Infer(env, fac)
	if err != nil {
		log.Printf("%+v", errors.Cause(err))
	}

	facType, ok := scheme.Type()
	fmt.Printf("fac Type: %v | isMonoType: %v | err: %v", facType, ok, err)

	// Output:
	// simple Type: Float | isMonoType: true | err: <nil>
	// fac Type: Float | isMonoType: true | err: <nil>

}
