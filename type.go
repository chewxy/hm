package hm

import "fmt"

// Type represents the metatype that is required to build a type system.
//		- TypeOp
//		- TypeVariable
//		- TypeConst
type Type interface {
	Name() string
	Contains(tv TypeVariable) bool
	Eq(Type) bool

	fmt.Formatter
	fmt.Stringer
}

// TypeOp is a type constructor. It takes n Types, and creates a new one from it.
type TypeOp interface {
	Type
	Types() Types

	SetTypes(...Type) TypeOp
	Clone() TypeOp
}

// TypeConst is a constant type. SetTypes(...) will yield the same exact values. It's useful for implementing atomic types. Formerly called Atomic
type TypeConst interface {
	TypeOp
	IsConstant() bool
}
