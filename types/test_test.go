package hmtypes

import "github.com/chewxy/hm"

const (
	proton  hm.TypeConst = "proton"
	neutron hm.TypeConst = "neutron"
	quark   hm.TypeConst = "quark"

	electron hm.TypeConst = "electron"
	positron hm.TypeConst = "positron"
	muon     hm.TypeConst = "muon"

	photon hm.TypeConst = "photon"
	higgs  hm.TypeConst = "higgs"
)

// useful copy pasta from the hm package
type mSubs map[hm.TypeVariable]hm.Type

func (s mSubs) Get(tv hm.TypeVariable) (hm.Type, bool)    { retVal, ok := s[tv]; return retVal, ok }
func (s mSubs) Add(tv hm.TypeVariable, t hm.Type) hm.Subs { s[tv] = t; return s }
func (s mSubs) Remove(tv hm.TypeVariable) hm.Subs         { delete(s, tv); return s }

func (s mSubs) Iter() []hm.Substitution {
	retVal := make([]hm.Substitution, len(s))
	var i int
	for k, v := range s {
		retVal[i] = hm.Substitution{k, v}
		i++
	}
	return retVal
}

func (s mSubs) Size() int { return len(s) }
func (s mSubs) Clone() hm.Subs {
	retVal := make(mSubs)
	for k, v := range s {
		retVal[k] = v
	}
	return retVal
}
