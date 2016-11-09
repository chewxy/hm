package hm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// for shits and giggles, here's a carbon atom: 6 protons, 6 neutrons, 4 electrons on the outer shell, 2 electrons in the inner shell
/*
	  **----E----**
	 **           **
	**   *--E---*   **
	|   *        *   |
	|  |   PNPN   |  |
	E  |   NPNP   |  E
	|  |   PNPN   |  |
	|  *         *   |
	**  *---E---*   **
	 **            **
	  **----E----**

*/
var carbon = []Type{
	electron,
	electron,
	proton, neutron, proton, neutron,
	electron, neutron, proton, neutron, proton, electron,
	proton, neutron, proton, neutron,
	electron,
	electron,
}

var lithium = []Type{
	electron,
	neutron, proton, neutron, proton, neutron, proton, electron,
	electron,
}

// here's an impossible thing
var types2 = []Type{
	neutron, photon, neutron, photon,
}

var leptons = Types{
	electron, muon,
}

var bosons = Types{
	photon, higgs,
}

func TestTypes(t *testing.T) {
	assert := assert.New(t)
	set := NewTypes()
	set2 := NewTypes(types2...)
	assert.Equal(0, len(set), "Expected empty set")
	assert.Equal(2, len(set2), "Expected a set with 2 elements")

	// add
	for _, tc := range carbon {
		set = set.Add(tc)
	}

	assert.Equal(Types{electron, proton, neutron}, set)
	assert.True(set.ContainsAll(carbon...))
	assert.False(set.ContainsAll(types2...))

	// subset and superset
	subset := leptons[:1]
	assert.True(subset.IsSubsetOf(set))
	assert.True(set.IsSupersetOf(subset))
	assert.False(subset.IsSupersetOf(set))
	assert.False(set.IsSubsetOf(subset))
	assert.False(subset.IsSubsetOf(set2))

	// intersect, union and differences
	assert.Equal(Types{neutron}, set.Intersect(set2))
	assert.Equal(Types{electron, proton, neutron, photon}, set.Union(set2))
	assert.Equal(Types{electron, proton}, set.Difference(set2))
	assert.Equal(Types{photon}, set2.Difference(set))
	assert.Equal(Types{electron, proton, photon}, set.SymmetricDifference(set2))

	// empty sets
	emptySet := NewTypes()
	assert.Nil(set.Intersect(emptySet))
	assert.Nil(emptySet.Intersect(set))
	assert.Equal(set, set.Union(emptySet))
	assert.Equal(set, emptySet.Union(set))

	// set equality
	assert.True(NewTypes(carbon...).Equals(NewTypes(lithium...)))
	assert.False(NewTypes(carbon...).Equals(NewTypes(lithium...)[:1]))
	assert.False(leptons.Equals(bosons))

	// string (for completeness)
	assert.Equal("Types[electron, proton, neutron]", set.String())
	assert.Equal("Types[electron, neutron, proton]", NewTypes(lithium...).String())

	// ToSlices (for completeness)
	assert.Equal([]Type{electron, proton, neutron}, set.ToSlice())
}
