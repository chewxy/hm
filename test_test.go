package hm

import "fmt"

// atom is a mock atomic type
type atom byte

const (
	nathaniel atom = iota
	adam
	jonathan
	osterman
)

func (t atom) Contains(tv TypeVariable) bool { return false }
func (t atom) Eq(other Type) bool {
	if ta, ok := other.(atom); ok {
		return ta == t
	}
	return false
}

func (t atom) Name() string                   { return t.String() }
func (t atom) Format(state fmt.State, c rune) { fmt.Fprintf(state, "atom(%d)", byte(t)) }
func (t atom) String() string                 { return fmt.Sprintf("atom(%d)", byte(t)) }
func (t atom) Types() Types                   { return nil }
func (t atom) SetTypes(...Type) TypeOp        { return t }
func (t atom) IsAtom() bool                   { return true }

// list is a mock type op. Think of it as `List a`
type list struct {
	t Type
}

func (t list) Contains(tv TypeVariable) bool {
	ttv, ok := t.t.(TypeVariable)
	if !ok {
		return false
	}

	return ttv.Eq(tv)
}

func (t list) Eq(other Type) bool {
	if tl, ok := other.(list); ok {
		return t.t.Eq(tl.t)
	}
	return false
}

func (t list) Name() string                   { return "List" }
func (t list) Format(state fmt.State, c rune) { fmt.Fprintf(state, "List %v", t.t) }
func (t list) String() string                 { return fmt.Sprintf("List %v", t.t) }
func (t list) Types() Types                   { return Types{t.t} }
func (t list) SetTypes(ts ...Type) TypeOp     { return list{ts[0]} }

// mirrorUniverseList is a List with a different name
type mirrorUniverseList struct {
	list
}

func (t mirrorUniverseList) Name() string { return "GoateeList" }

// malformed is an incomplete Type
type malformed struct{}

func (t malformed) Name() string                   { return "malformed" }
func (t malformed) Contains(tv TypeVariable) bool  { return false }
func (t malformed) Eq(other Type) bool             { return false }
func (t malformed) Format(state fmt.State, c rune) { fmt.Fprintf(state, "malformed") }
func (t malformed) String() string                 { return "malformed" }
