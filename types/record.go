package hmtypes

import (
	"fmt"

	"github.com/chewxy/hm"
)

// Tuple is a basic tuple type. It takes an optional name
type Tuple struct {
	ts   []hm.Type
	name string
}

// NewTupleType creates a new Tuple
func NewTupleType(name string, ts ...hm.Type) *Tuple {
	return &Tuple{
		ts:   ts,
		name: name,
	}
}

func (t *Tuple) Apply(subs hm.Subs) hm.Substitutable {
	ts := t.apply(subs)
	return NewTupleType(t.name, ts...)
}

func (t *Tuple) FreeTypeVar() hm.TypeVarSet {
	var tvs hm.TypeVarSet
	for _, v := range t.ts {
		tvs = v.FreeTypeVar().Union(tvs)
	}
	return tvs
}

func (t *Tuple) Name() string {
	if t.name != "" {
		return t.name
	}
	return t.String()
}

func (t *Tuple) Normalize(k, v hm.TypeVarSet) (T hm.Type, err error) {
	var ts []hm.Type
	if ts, err = t.normalize(k, v); err != nil {
		return nil, err
	}
	return NewTupleType(t.name, ts...), nil
}

func (t *Tuple) Types() hm.Types {
	ts := hm.BorrowTypes(len(t.ts))
	copy(ts, t.ts)
	return ts
}

func (t *Tuple) Format(f fmt.State, c rune) {
	f.Write([]byte("("))
	for i, v := range t.ts {
		if i < len(t.ts)-1 {
			fmt.Fprintf(f, "%v, ", v)
		} else {
			fmt.Fprintf(f, "%v)", v)
		}
	}
}

func (t *Tuple) String() string { return fmt.Sprintf("%v", t) }

func (t *Tuple) Eq(other hm.Type) bool {
	if ot, ok := other.(*Tuple); ok {
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

// Clone implements Cloner
func (t *Tuple) Clone() interface{} {
	retVal := new(Tuple)
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

// internal function to be used by Tuple.Apply and Record.Apply
func (t *Tuple) apply(subs hm.Subs) []hm.Type {
	ts := make([]hm.Type, len(t.ts))
	for i, v := range t.ts {
		ts[i] = v.Apply(subs).(hm.Type)
	}
	return ts
}

// internal function to be used by Tuple.Normalize and Record.Normalize
func (t *Tuple) normalize(k, v hm.TypeVarSet) ([]hm.Type, error) {
	ts := make([]hm.Type, len(t.ts))
	var err error
	for i, tt := range t.ts {
		if ts[i], err = tt.Normalize(k, v); err != nil {
			return nil, err
		}
	}
	return ts
}

// Field is a name-type pair.
type Field struct {
	Name string
	Type hm.Type
}

// Record is a basic record type. It's like Tuple except there are named fields. It takes an optional name.
type Record struct {
	Tuple
	ns []string // field names
}

// NewRecordType creates a new Record hm.Type
func NewRecordType(name string, fields ...Field) *Record {
	ts := make([]hm.Type, len(fields))
	ns := make([]string, len(fields))
	for i := range fields {
		ts[i] = fields[i].Name
		ns[i] = fields[i].Type
	}
	return &Record{
		Tuple: Tuple{
			ts:   ts,
			name: name,
		},
		ns: ns,
	}
}

func (t *Record) Apply(subs hm.Subs) hm.Substitutable {
	ts := t.apply(subs)
	return &Record{
		Tuple: Tuple{
			ts:   ts,
			name: t.name,
		},
		ns: t.ns,
	}
}

func (t *Record) Normalize(k, v hm.TypeVarSet) (T hm.Type, err error) {
	var ts []hm.Type
	if ts, err = t.normalize(k, v); err != nil {
		return nil, err
	}
	return &Record{
		Tuple: Tuple{
			ts:   ts,
			name: t.name,
		},
		ns: t.ns,
	}, nil
}

func (t *Record) Format(f fmt.State, c rune) {
	if t.name != "" {
		f.Write([]byte(t.name))
	}
	f.Write([]byte("{"))
	for i, v := range t.ts {
		if i < len(t.ts)-1 {
			fmt.Fprintf(f, "%v: %v, ", t.ns[i], v)
		} else {
			fmt.Fprintf(f, "%v: %v}", t.ns[i], v)
		}
	}
}

func (t *Record) Eq(other hm.Type) bool {
	if ot, ok := other.(*Record); ok {
		if len(ot.ts) != len(t.ts) {
			return false
		}
		for i, v := range t.ts {
			if t.ns[i] != ot.ns[i] {
				return false
			}
			if !v.Eq(ot.ts[i]) {
				return false
			}
		}
		return true
	}
	return false
}
