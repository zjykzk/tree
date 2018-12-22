package tree

import (
	"fmt"
	"strings"
)

// Key the key in the rb tree
type Key interface {
	// CompareTo returns a negetive integer, zero, or a positive integer as this key is less than,
	// equal to, or greater than the specified key
	//
	// The implementor must ensure sgn(x.compareTo(y)) == -sgn(y.compareTo(x)) for all x and y.
	//
	// The implementor must also ensure that the relation is transitive:
	// (x.compareTo(y)>0 && y.compareTo(z)>0) implies x.compareTo(z)>0.
	//
	// Finally, the implementor must ensure that x.compareTo(y)==0
	// implies that sgn(x.compareTo(z)) == sgn(y.compareTo(z)), for all z.
	//
	// It is strongly recommended, but not strictly required that
	// (x.compareTo(y)==0) == (x.equals(y)).  Generally speaking, any
	//
	// In the foregoing description, the notation
	// sgn(expression) designates the mathematical signum function, which is defined to return one of
	// -1, 0, or 1 according to whether the value of expression is negative, zero or positive.
	CompareTo(k Key) int
}

// LLRBTree the red-black tree. The algorithms are adaptations of those in
// http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf
type LLRBTree struct {
	root *node
}

// Get returns the value to which the specified key is mapped, and false if this tree contains no
// mapping for th key.
func (t *LLRBTree) Get(k Key) (v interface{}, ok bool) {
	n := t.root
	for n != nil {
		switch cmp := k.CompareTo(n.key); {
		case cmp == 0:
			return n.value, true
		case cmp > 0:
			n = n.right
		default:
			n = n.left
		}
	}
	return nil, false
}

// Put associates the specified value with the specified key in this tree.
// If the tree previously contained a mapping for the key, the old value is replaced.
func (t *LLRBTree) Put(k Key, v interface{}) (old interface{}) {
	t.root, old = insert(t.root, k, v)
	t.root.color = black
	return
}

// Remove removes the mapping for a key from this map if it is present (optional operation)
func (t *LLRBTree) Remove(k Key) (v interface{}) {
	v, ok := t.Get(k)
	if !ok {
		return
	}

	t.root = remove(t.root, k)
	t.root.color = black
	return v
}

const (
	red   = true
	black = false
)

type node struct {
	key   Key
	value interface{}

	left, right *node
	color       bool
}

func rotateLeft(n *node) *node {
	x := n.right
	n.right = x.left
	x.left = n
	x.color = n.color
	n.color = red
	return x
}

func rotateRight(n *node) *node {
	x := n.left
	n.left = x.right
	x.right = n
	x.color = n.color
	n.color = red
	return x
}

func colorFlip(n *node) *node {
	n.color = !n.color
	n.left.color = !n.left.color
	n.right.color = !n.right.color
	return n
}

func isRed(n *node) bool {
	return n != nil && n.color == red
}

func insert(n *node, k Key, v interface{}) (h *node, old interface{}) {
	if n == nil {
		return &node{key: k, value: v, color: red}, nil
	}

	switch cmp := k.CompareTo(n.key); {
	case cmp == 0:
		old, n.value = n.value, v
	case cmp > 0:
		n.right, old = insert(n.right, k, v)
	default:
		n.left, old = insert(n.left, k, v)
	}

	h = fixUp(n)
	return
}

func remove(n *node, k Key) *node {
	cmp := k.CompareTo(n.key)

	// left
	if cmp < 0 {
		if !isRed(n.left) && !isRed(n.left.left) { // push red right if necessary
			n = moveRedLeft(n) // move down (left)
		}
		n.left = remove(n.left, k)

		return fixUp(n)
	}

	// right or equal
	if isRed(n.left) {
		n = rotateRight(n) // rotate to push red right
		cmp = k.CompareTo(n.key)
	}

	if cmp == 0 && n.right == nil { // equal at bottom, delete it
		return nil
	}

	if !isRed(n.right) && !isRed(n.right.left) { // push red right if necessary
		n = moveRedRight(n)
		cmp = k.CompareTo(n.key)
	}

	if cmp == 0 { // EQUAL (not at bottom)
		n.key, n.value = min(n.right)   // replace current node with successor key, value
		n.right, _ = removeMin(n.right) // delete successor
		return fixUp(n)
	}

	n.right = remove(n.right, k) // move down (right)
	return fixUp(n)              // fix right-leaning red links and eliminate 4-nodes on the way up
}

func min(n *node) (k Key, v interface{}) {
	for n.left != nil {
		n = n.left
	}
	return n.key, n.value
}

func removeMin(n *node) (r *node, v interface{}) {
	if n.left == nil {
		return nil, n.value
	}

	if !isRed(n.left) && !isRed(n.left.left) {
		n = moveRedLeft(n)
	}
	n.left, v = removeMin(n.left)

	return fixUp(n), v
}

func fixUp(n *node) *node {
	if isRed(n.right) && !isRed(n.left) {
		n = rotateLeft(n)
	}

	if isRed(n.left) && isRed(n.left.left) {
		n = rotateRight(n)
	}

	if isRed(n.left) && isRed(n.right) {
		colorFlip(n)
	}

	return n
}

func moveRedRight(n *node) *node {
	colorFlip(n)
	if isRed(n.left.left) {
		n = rotateRight(n)
		colorFlip(n)
	}

	return n
}

func moveRedLeft(n *node) *node {
	colorFlip(n)
	if isRed(n.right.left) {
		n.right = rotateRight(n.right)
		n = rotateLeft(n)
		colorFlip(n)
	}

	return n
}

// dot related functions
func dotNode(n *node) string {
	c := "gray"
	if isRed(n) {
		c = "red"
	}
	return fmt.Sprintf(`%v [shape=circle,color=%s,label="%v:%v"];`, n.key, c, n.key, n.value)
}

func dotNodes(n *node) string {
	if n == nil {
		return ""
	}
	return dotNode(n) + dotNodes(n.left) + dotNodes(n.right)
}

func dotEdges(n *node) string {
	if n == nil {
		return ""
	}

	s := ""
	if n.left != nil {
		s += fmt.Sprintf("%v->%v;", n.key, n.left.key)
	}
	if n.right != nil {
		s += fmt.Sprintf("%v->%v;", n.key, n.right.key)
	}

	s += dotEdges(n.left)
	s += dotEdges(n.right)
	return s
}

func (t *LLRBTree) dotString() string {
	b := &strings.Builder{}

	b.WriteString(`digraph G {`)
	b.WriteString(dotNodes(t.root))
	b.WriteString(dotEdges(t.root))
	b.WriteByte('}')

	return b.String()
}
