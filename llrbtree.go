package tree

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
	CompareTo(k interface{}) int
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
		switch r := n.key.CompareTo(k); {
		case r == 0:
			return n.value, true
		case r < 0:
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

	if isRed(n.left) && isRed(n.right) {
		colorFlip(n)
	}

	switch r := n.key.CompareTo(k); {
	case r == 0:
		old, n.value = n.value, v
	case r < 0:
		n.right, old = insert(n.right, k, v)
	default:
		n.left, old = insert(n.left, k, v)
	}

	if isRed(n.right) && !isRed(n.left) {
		n = rotateLeft(n)
	}

	if isRed(n.left) && isRed(n.left.left) {
		n = rotateRight(n)
	}

	h = n
	return
}
