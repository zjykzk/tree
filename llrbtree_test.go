package tree

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTree string

func (m mockTree) CompareTo(m1 Key) int {
	return strings.Compare(string(m), string(m1.(mockTree)))
}

func TestTree(t *testing.T) {
	tree := &LLRBTree{}

	// remove unpresent key
	tree.Remove(mockTree("1"))

	// Put & Get
	tree.Put(mockTree("1"), 1)
	tree.Put(mockTree("2"), 2)
	tree.Put(mockTree("6"), 6)
	tree.Put(mockTree("7"), 7)
	tree.Put(mockTree("8"), 8)
	tree.Put(mockTree("9"), 9)
	tree.Put(mockTree("3"), 3)
	tree.Put(mockTree("4"), 4)
	assert.Equal(t, 8, tree.Size())

	t.Log(tree.dotString())

	v, ok := tree.Get(mockTree("1"))
	assert.Equal(t, 1, v.(int))
	assert.True(t, ok)

	// get unpresent key
	_, ok = tree.Get(mockTree("_"))
	assert.False(t, ok)

	// update value
	old := tree.Put(mockTree("2"), 4)
	assert.Equal(t, 2, old.(int))
	v, ok = tree.Get(mockTree("2"))
	assert.Equal(t, 4, v.(int))
	assert.True(t, ok)

	// remove smallest
	v = tree.Remove(mockTree("1"))
	assert.Equal(t, 1, v.(int))
	v, ok = tree.Get(mockTree("1"))
	assert.False(t, ok)
	t.Log(tree.dotString())
	assert.Equal(t, 7, tree.Size())

	tree.Put(mockTree("1"), 1)

	// remove greatest
	v = tree.Remove(mockTree("8"))
	assert.Equal(t, 8, v.(int))
	v, ok = tree.Get(mockTree("8"))
	assert.False(t, ok)

	tree.Put(mockTree("5"), 5)
	tree.Put(mockTree("50"), 50)
	tree.Put(mockTree("40"), 40)
	tree.Put(mockTree("90"), 90)

	// remove key
	v = tree.Remove(mockTree("2"))
	assert.Equal(t, 4, v.(int))
	v, ok = tree.Get(mockTree("2"))
	assert.False(t, ok)

	t.Log(tree.dotString())
	// first
	k, v := tree.First()
	assert.Equal(t, mockTree("1"), k.(mockTree))
	assert.Equal(t, 1, v.(int))
	// last
	k, v = tree.Last()
	assert.Equal(t, mockTree("90"), k.(mockTree))
	assert.Equal(t, 90, v.(int))
}
