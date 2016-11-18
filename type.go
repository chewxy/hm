package hm

import "fmt"

type Type interface {
	Substitutable
	Name() string // Name is the name of the constructor
	Normalize(TypeVarSet, TypeVarSet) (Type, error)
	Types() Types

	fmt.Formatter
	fmt.Stringer
}

type Substitutable interface {
	Apply(Subs) Substitutable
	FreeTypeVar() TypeVarSet
}

type TypeConst string

func (t TypeConst) Name() string                            { return string(t) }
func (t TypeConst) Apply(Subs) Substitutable                { return t }
func (t TypeConst) FreeTypeVar() TypeVarSet                 { return nil }
func (t TypeConst) Normalize(k, v TypeVarSet) (Type, error) { return t, nil }
func (t TypeConst) Types() Types                            { return nil }
func (t TypeConst) String() string                          { return string(t) }
func (t TypeConst) Format(s fmt.State, c rune)              { fmt.Fprintf(s, "%s", string(t)) }
