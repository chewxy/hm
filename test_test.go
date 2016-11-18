package hm

import "fmt"

const (
	proton  TypeConst = "proton"
	neutron TypeConst = "neutron"
	quark   TypeConst = "quark"

	electron TypeConst = "electron"
	positron TypeConst = "positron"
	muon     TypeConst = "muon"

	photon TypeConst = "photon"
	higgs  TypeConst = "higgs"
)

type list struct {
	t Type
}

func (l list) Name() string                  { return "List" }
func (l list) Apply(subs Subs) Substitutable { l.t = l.t.Apply(subs).(Type); return l }
func (l list) FreeTypeVar() TypeVarSet       { return l.t.FreeTypeVar() }
func (l list) Format(s fmt.State, c rune)    { fmt.Fprintf(s, "List %v", l.t) }
func (l list) String() string                { return fmt.Sprintf("%v", l) }
func (l list) Normalize(k, v TypeVarSet) (Type, error) {
	var t Type
	var err error
	if t, err = l.t.Normalize(k, v); err != nil {
		return nil, err
	}
	l.t = t
	return l, nil
}
func (l list) Types() Types { return Types{l.t} }

type mirrorUniverseList struct {
	t Type
}

func (l mirrorUniverseList) Name() string                  { return "GoateeList" }
func (l mirrorUniverseList) Apply(subs Subs) Substitutable { l.t = l.t.Apply(subs).(Type); return l }
func (l mirrorUniverseList) FreeTypeVar() TypeVarSet       { return l.t.FreeTypeVar() }
func (l mirrorUniverseList) Format(s fmt.State, c rune)    { fmt.Fprintf(s, "List %v", l.t) }
func (l mirrorUniverseList) String() string                { return fmt.Sprintf("%v", l) }
func (l mirrorUniverseList) Normalize(k, v TypeVarSet) (Type, error) {
	var t Type
	var err error
	if t, err = l.t.Normalize(k, v); err != nil {
		return nil, err
	}
	l.t = t
	return l, nil
}
func (l mirrorUniverseList) Types() Types { return Types{l.t} }
