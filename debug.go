package tree

import (
	"fmt"
	"strings"
)

// dot related functions
func dotNode(n *node) string {
	return fmt.Sprintf(`%v [shape=circle,label="%v:%v"];`, n.key, n.key, n.value)
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
		s += fmt.Sprintf("%v->%v[color=%s];", n.key, n.left.key, color(n.left))
	}
	if n.right != nil {
		s += fmt.Sprintf("%v->%v[color=%s];", n.key, n.right.key, color(n.right))
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
