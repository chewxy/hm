package hm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	num      = NewSimpleTypeClass("Num")
	floating = NewSimpleTypeClass("Floating")
	ord      = NewSimpleTypeClass("Ord")

	eq = NewSimpleTypeClass("Eq")
)

var typeclasses1 = []TypeClass{
	num, num, floating, floating, ord, ord,
}

var typeclasses2 = []TypeClass{
	ord, eq, ord, eq,
}

func TestTypeClassSet(t *testing.T) {
	assert := assert.New(t)
	set := NewTypeClassSet()
	set2 := NewTypeClassSet(typeclasses2...)
	assert.Equal(0, len(set.s), "Expected empty set")
	assert.Equal(2, len(set2.s), "Expected a set with 2 elements")

	// add
	for _, tc := range typeclasses1 {
		set = set.Add(tc)
	}

	assert.Equal(NewTypeClassSet(num, floating, ord), set)
	assert.True(set.ContainsAll(typeclasses1...))
	assert.False(set.ContainsAll(typeclasses2...))

	// subset and superset
	subset := NewTypeClassSet(num)
	assert.True(subset.IsSubsetOf(set))
	assert.True(set.IsSupersetOf(subset))
	assert.False(subset.IsSupersetOf(set))
	assert.False(set.IsSubsetOf(subset))
	assert.False(subset.IsSubsetOf(set2))

	// intersect, union and differences
	assert.Equal(NewTypeClassSet(ord), set.Intersect(set2))
	assert.Equal(NewTypeClassSet(num, floating, ord, eq), set.Union(set2))
	assert.Equal(NewTypeClassSet(num, floating), set.Difference(set2))
	assert.Equal(NewTypeClassSet(eq), set2.Difference(set))
	assert.Equal(NewTypeClassSet(num, floating, eq), set.SymmetricDifference(set2))

	// empty sets
	emptySet := NewTypeClassSet()
	assert.Nil(set.Intersect(emptySet))
	assert.Nil(emptySet.Intersect(set))
	assert.Equal(set, set.Union(emptySet))
	assert.Equal(set, emptySet.Union(set))

	// set equality
	assert.True(NewTypeClassSet(ord, eq).Equals(NewTypeClassSet(eq, ord)))
	assert.False(NewTypeClassSet(ord).Equals(NewTypeClassSet(eq, ord)))
	assert.False(NewTypeClassSet(ord).Equals(NewTypeClassSet(eq)))

	// string (for completeness)
	assert.Equal("TypeClassSet[Num, Floating, Ord]", set.String())

	// ToSlices (for completeness)
	assert.Equal([]TypeClass{num, floating, ord}, set.ToSlice())
}
