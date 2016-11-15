package hm

import "fmt"

type particle byte

const (
	proton particle = iota
	neutron
	quark

	electron
	positron
	muon

	photon
	higgs
)

func (t particle) Contains(tv TypeVariable) bool { return false }
func (t particle) Eq(other Type) bool {
	if ta, ok := other.(particle); ok {
		return ta == t
	}
	return false
}

func (t particle) Name() string                   { return t.String() }
func (t particle) Format(state fmt.State, c rune) { fmt.Fprintf(state, t.String()) }
func (t particle) Types() Types                   { return nil }
func (t particle) Clone() TypeOp                  { return t }
func (t particle) Replace(Type, Type) TypeOp      { return t }
func (t particle) IsConstant() bool               { return true }
func (t particle) String() string {
	switch t {
	case proton:
		return "proton"
	case neutron:
		return "neutron"
	case quark:
		return "quark"
	case electron:
		return "electron"
	case positron:
		return "positron"
	case muon:
		return "muon"
	case photon:
		return "photon"
	case higgs:
		return "higgs"
	default:
		return fmt.Sprintf("atom(%d)", byte(t))
	}
}

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
func (t list) Clone() TypeOp {
	retVal := list{}
	switch tt := t.t.(type) {
	case TypeVariable:
		retVal.t = tt
	case TypeOp:
		retVal.t = tt.Clone()
	}
	return retVal
}

func (t list) Replace(what, with Type) TypeOp {
	switch tt := t.t.(type) {
	case TypeVariable:
		if tt.Eq(what) {
			t.t = with
		}
	case TypeConst:
		// do nothing
	case TypeOp:
		t.t = tt.Replace(what, with)
	default:
		panic("WTF")
	}
	return t
}

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

func typeEqAnyVar(a, b Type) bool {
	switch at := a.(type) {
	case *FunctionType:
		if bt, ok := b.(*FunctionType); ok {
			if !typeEqAnyVar(at.ts[0], bt.ts[0]) {
				return false
			}
			if !typeEqAnyVar(at.ts[1], bt.ts[1]) {
				return false
			}
			return true
		}
		return false
	case TypeVariable:
		if bt, ok := b.(TypeVariable); ok {
			if at.name == "∀" || bt.name == "∀" {
				return true
			}
			return at.name == bt.name
		}
		return false
	default:
		return a.Eq(b)
	}
}
