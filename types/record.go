package hmtypes

import (
	"fmt"

	"github.com/chewxy/hm"
)

// Record is a basic record/tuple type. It takes an optional name.
type Record struct {
	ts   []hm.Type
	name string
}

// NewRecordType creates a new Record hm.Type
func NewRecordType(name string, ts ...hm.Type) *Record {
	return &Record{
		ts:   ts,
		name: name,
	}
}

func (t *Record) Apply(subs hm.Subs) hm.Substitutable {
	ts := make([]hm.Type, len(t.ts))
	for i, v := range t.ts {
		ts[i] = v.Apply(subs).(hm.Type)
	}
	return NewRecordType(t.name, ts...)
}

func (t *Record) FreeTypeVar() hm.TypeVarSet {
	var tvs hm.TypeVarSet
	for _, v := range t.ts {
		tvs = v.FreeTypeVar().Union(tvs)
	}
	return tvs
}

func (t *Record) Name() string {
	if t.name != "" {
		return t.name
	}
	return t.String()
}

func (t *Record) Normalize(k, v hm.TypeVarSet) (hm.Type, error) {
	ts := make([]hm.Type, len(t.ts))
	var err error
	for i, tt := range t.ts {
		if ts[i], err = tt.Normalize(k, v); err != nil {
			return nil, err
		}
	}
	return NewRecordType(t.name, ts...), nil
}

func (t *Record) Types() hm.Types {
	ts := hm.BorrowTypes(len(t.ts))
	copy(ts, t.ts)
	return ts
}

func (t *Record) Eq(other hm.Type) bool {
	if ot, ok := other.(*Record); ok {
		if len(ot.ts) != len(t.ts) {
			return false
		}
		for i, v := range t.ts {
			if !v.Eq(ot.ts[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (t *Record) Format(f fmt.State, c rune) {
	f.Write([]byte("("))
	for i, v := range t.ts {
		if i < len(t.ts)-1 {
			fmt.Fprintf(f, "%v, ", v)
		} else {
			fmt.Fprintf(f, "%v)", v)
		}
	}

}

func (t *Record) String() string { return fmt.Sprintf("%v", t) }

// Clone implements Cloner
func (t *Record) Clone() interface{} {
	retVal := new(Record)
	ts := hm.BorrowTypes(len(t.ts))
	for i, tt := range t.ts {
		if c, ok := tt.(Cloner); ok {
			ts[i] = c.Clone().(hm.Type)
		} else {
			ts[i] = tt
		}
	}
	retVal.ts = ts
	retVal.name = t.name

	return retVal
}
