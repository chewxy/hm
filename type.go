package hm

import "fmt"

// Type represents all the possible type constructors.
type Type interface {
	Substitutable
	Name() string                                   // Name is the name of the constructor
	Normalize(TypeVarSet, TypeVarSet) (Type, error) // Normalize normalizes all the type variable names in the type
	Types() Types                                   // If the type is made up of smaller types, then it will return them
	Eq(Type) bool                                   // equality operation

	fmt.Formatter
	fmt.Stringer
}

// Substitutable is any type that can have a set of substitutions applied on it, as well as being able to know what its free type variables are
type Substitutable interface {
	Apply(Subs) Substitutable
	FreeTypeVar() TypeVarSet
}

// TypeConst are the default implementation of a constant type. Feel free to implement your own
type TypeConst string

func (t TypeConst) Name() string                            { return string(t) }
func (t TypeConst) Apply(Subs) Substitutable                { return t }
func (t TypeConst) FreeTypeVar() TypeVarSet                 { return nil }
func (t TypeConst) Normalize(k, v TypeVarSet) (Type, error) { return t, nil }
func (t TypeConst) Types() Types                            { return nil }
func (t TypeConst) String() string                          { return string(t) }
func (t TypeConst) Format(s fmt.State, c rune)              { fmt.Fprintf(s, "%s", string(t)) }
func (t TypeConst) Eq(other Type) bool                      { return other == t }
