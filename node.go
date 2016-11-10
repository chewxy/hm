package hm

// Expr represents an expression.
type Node interface {
	Body() Node
}

type UntypedNode interface {
	// ???
}

type Namer interface {
	Name() string
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

// Lit is a node that represents an identifier or a literal
type Lit interface {
	Var
	IsLit() bool
}

// Var is a node that represents `var x int`
type Var interface {
	Node
	Namer
	Typer
}

// Lambda is a node that represents a function definition
type Lambda interface {
	Node
	Namer
	IsLambda() bool
}

// Apply is a node that represents a function call/application
type Apply interface {
	Node
	Fn() Lambda
}

// Let is a node that represents a Haskell-like `let`: let x = blah in blahbody
type Let interface {
	Node
	Namer

	Def() Node
}

// LetRec is a recursive version of the above. It's implemented here because it's a useful thing to have
type LetRec interface {
	Let
	IsLetRec() bool
}
