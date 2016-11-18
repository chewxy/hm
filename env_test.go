package hm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleEnv(t *testing.T) {
	assert := assert.New(t)
	var orig, env Env
	var expected SimpleEnv

	// Add
	orig = make(SimpleEnv)
	orig = orig.Add("foo", &Scheme{
		tvs: TypeVarSet{'a', 'b', 'c'},
		t:   TypeVariable('a'),
	})
	orig = orig.Add("bar", &Scheme{
		tvs: TypeVarSet{'b', 'c', 'd'},
		t:   TypeVariable('a'),
	})
	orig = orig.Add("baz", &Scheme{
		tvs: TypeVarSet{'a', 'b', 'c'},
		t:   neutron,
	})
	qs := &Scheme{
		tvs: TypeVarSet{'a', 'b'},
		t:   proton,
	}
	orig = orig.Add("qux", qs)

	expected = SimpleEnv{
		"foo": &Scheme{
			tvs: TypeVarSet{'a', 'b', 'c'},
			t:   TypeVariable('a'),
		},
		"bar": &Scheme{
			tvs: TypeVarSet{'b', 'c', 'd'},
			t:   TypeVariable('a'),
		},
		"baz": &Scheme{
			tvs: TypeVarSet{'a', 'b', 'c'},
			t:   neutron,
		},
		"qux": &Scheme{
			tvs: TypeVarSet{'a', 'b'},
			t:   proton,
		},
	}
	assert.Equal(expected, orig)

	// Get
	s, ok := orig.SchemeOf("qux")
	if s != qs || !ok {
		t.Error("Expected to get scheme of \"qux\"")
	}

	// Remove
	orig = orig.Remove("qux")
	delete(expected, "qux")
	assert.Equal(expected, orig)

	// Clone
	env = orig.Clone()
	assert.Equal(orig, env)

	subs := mSubs{
		'a': proton,
		'b': neutron,
		'd': electron,
		'e': proton,
	}

	env = env.Apply(subs).(Env)
	expected = SimpleEnv{
		"foo": &Scheme{
			tvs: TypeVarSet{'a', 'b', 'c'},
			t:   TypeVariable('a'),
		},
		"bar": &Scheme{
			tvs: TypeVarSet{'b', 'c', 'd'},
			t:   proton,
		},
		"baz": &Scheme{
			tvs: TypeVarSet{'a', 'b', 'c'},
			t:   neutron,
		},
	}
	assert.Equal(expected, env)

	env = orig.Clone()
	ftv := env.FreeTypeVar()
	correctFTV := TypeVarSet{'a'}

	if !correctFTV.Equals(ftv) {
		t.Errorf("Expected freetypevars to be equal. Got %v instead", ftv)
	}
}
