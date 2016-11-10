package hm

import "fmt"

type TyperNode interface {
	Node
	Typer
}

type λ struct {
	name string
	body Node
}

func (n λ) Name() string   { return n.name }
func (n λ) Body() Node     { return n.body }
func (n λ) IsLambda() bool { return true }

type lit string

func (n lit) Name() string   { return string(n) }
func (n lit) Body() Node     { return n }
func (n lit) Type() Type     { return nil }
func (n lit) IsLit() bool    { return true }
func (n lit) IsLambda() bool { return true }

type app struct {
	f   Lambda
	arg Node
}

func (n app) Fn() Lambda { return n.f }
func (n app) Body() Node { return n.arg }
func (n app) Arg() Node  { return n.arg }

type letrec struct {
	name string
	def  Lambda
	in   Apply
}

func (n letrec) Name() string     { return n.name }
func (n letrec) Def() Node        { return n.def }
func (n letrec) Body() Node       { return n.in }
func (n letrec) Children() []Node { return []Node{n.def, n.in} }
func (n letrec) IsLetRec() bool   { return true }

type prim byte

const (
	Float prim = iota
	Bool
)

func (t prim) Name() string                  { return t.String() }
func (t prim) Contains(tv TypeVariable) bool { return false }
func (t prim) Eq(other Type) bool {
	ot, ok := other.(prim)
	if !ok {
		return false
	}
	return ot == t
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
func (t prim) IsConst() bool { return true }

//Phillip Greenspun's tenth law says:
//		"Any sufficiently complicated C or Fortran program contains an ad hoc, informally-specified, bug-ridden, slow implementation of half of Common Lisp."
//
// So let's implement a half-arsed lisp (Or rather, an AST that can optionally be executed upon if you write the correct interpreter)!
func Example() {
	// haskell envy in a greenspun's tenth law example function!
	//
	// We'll assume the following is the "input" code
	// 		let fac n = if n == 0 then 1 else n * fac (n - 1) in fac 5
	// and what we have is the AST

	fmt.Println("HELLO")

	fac := letrec{
		"fac",
		λ{
			"n",
			app{
				λ{
					"λ0",
					app{
						λ{
							"λ1",
							app{
								lit("if"),
								app{lit("iszero"), lit("n")},
							},
						},
						lit("1"),
					},
				},
				app{
					λ{
						"λ2",
						app{lit("mul"), lit("n")},
					},
					app{
						lit("fac"),
						app{lit("--"), lit("n")},
					},
				},
			},
		},
		app{lit("fac"), lit("5")},
	}

	predef := map[string]Type{
		"--": NewFnType(Float, Float),
	}

	env := NewSimpleEnv(WithDict(predef))
	t, err := Infer(fac, env)
	fmt.Printf("Type: %v | err: %v", t, err)

	// Ouput:
	// Type: Float | err: <nil>

}
