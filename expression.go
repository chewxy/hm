package hm

type Namer interface {
	Name() string
}

type Typer interface {
	Type() Type
}

type Expression interface {
	Body() Expression
}

type Var interface {
	Expression
	Namer
	Typer
}

type Literal interface {
	Var
	IsLit() bool
}

type Apply interface {
	Expression
	Fn() Expression
}

type LetRec interface {
	Let
	IsRecursive() bool
}

type Let interface {
	// let name = def in body
	Expression
	Namer
	Def() Expression
}

type Lambda interface {
	Expression
	Namer
	IsLambda() bool
}
