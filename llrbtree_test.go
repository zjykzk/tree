package tree

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var nameSufix = 1

func dotName(n *node) string {
	return fmt.Sprintf("no_%v%d", n.key, nameSufix)
}

// dot related functions
func dotNode(n *node) string {
	return fmt.Sprintf(`%v [shape=circle,label="%v:%v"];`, dotName(n), n.key, n.value)
}

func dotNodes(n *node) string {
	if n == nil {
		return ""
	}
	return dotNode(n) + dotNodes(n.left) + dotNodes(n.right)
}

func color(n *node) string {
	if isRed(n) {
		return "red"
	}
	return "black"
}

func dotEdges(n *node) string {
	if n == nil {
		return ""
	}

	s := ""
	if n.left != nil {
		s += fmt.Sprintf("%v->%v[color=%s];", dotName(n), dotName(n.left), color(n.left))
	}
	if n.right != nil {
		s += fmt.Sprintf("%v->%v[color=%s];", dotName(n), dotName(n.right), color(n.right))
	}

	s += dotEdges(n.left)
	s += dotEdges(n.right)
	return s
}

func (t *LLRBTree) nodeAndEdges() string {
	b := &strings.Builder{}

	b.WriteString(dotNodes(t.root))
	b.WriteString(dotEdges(t.root))

	nameSufix++

	return b.String()
}

func printGraph(s string) {
	f, err := os.Create("dot.dot")
	if err != nil {
		fmt.Printf("print graph error:%s\n", err)
		return
	}
	f.WriteString(fmt.Sprintf("digraph G{%s}\n", s))
	f.Close()
}

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

	graph := ""
	graph += tree.nodeAndEdges()

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
	graph += tree.nodeAndEdges()
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

	graph += tree.nodeAndEdges()
	// first
	k, v := tree.First()
	assert.Equal(t, mockTree("1"), k.(mockTree))
	assert.Equal(t, 1, v.(int))
	// last
	k, v = tree.Last()
	assert.Equal(t, mockTree("90"), k.(mockTree))
	assert.Equal(t, 90, v.(int))

	t.Log(graph)
	printGraph(graph)
}

func TestPutAll(t *testing.T) {
	tree := &LLRBTree{}

	// Put & Get
	tree.Put(mockTree("1"), 1)
	tree.Put(mockTree("2"), 2)
	tree.Put(mockTree("6"), 6)
	tree.Put(mockTree("7"), 7)
	tree.Put(mockTree("8"), 8)
	tree.Put(mockTree("9"), 9)
	tree.Put(mockTree("3"), 3)
	tree.Put(mockTree("4"), 4)

	newTree := &LLRBTree{}
	newTree.PutAll(tree)

	checkerTrue := func(i int) {
		v, ok := newTree.Get(mockTree(strconv.Itoa(i)))
		assert.True(t, ok)
		assert.Equal(t, i, v.(int))
	}
	checkerTrue(1)
	checkerTrue(2)
	checkerTrue(3)
	checkerTrue(4)
	checkerTrue(6)
	checkerTrue(7)
	checkerTrue(8)
	checkerTrue(9)

	checkerFalse := func(i int) {
		_, ok := newTree.Get(mockTree(strconv.Itoa(i)))
		assert.False(t, ok)
	}
	checkerFalse(5)

	g := newTree.nodeAndEdges() + tree.nodeAndEdges()

	tree.Clear()
	tree.Put(mockTree("20"), 20)
	newTree.PutAll(tree)
	checkerTrue(20)
	g += newTree.nodeAndEdges()

	tree.Clear()
	newTree.PutAll(tree)

	printGraph(g)
}
