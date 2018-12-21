package tree

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTree string

func (m mockTree) CompareTo(m1 interface{}) int {
	return strings.Compare(string(m), string(m1.(mockTree)))
}

func TestPutAndGet(t *testing.T) {
	tree := &LLRBTree{}
	tree.Put(mockTree("1"), 1)
	tree.Put(mockTree("2"), 2)
	v, ok := tree.Get(mockTree("1"))
	assert.Equal(t, 1, v.(int))
	assert.True(t, ok)

	v, ok = tree.Get(mockTree("2"))
	assert.Equal(t, 2, v.(int))
	assert.True(t, ok)

	_, ok = tree.Get(mockTree("3"))
	assert.False(t, ok)

	tree.Put(mockTree("2"), 4)
	v, ok = tree.Get(mockTree("2"))
	assert.Equal(t, 4, v.(int))
	assert.True(t, ok)
}
