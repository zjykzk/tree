package tree

import (
	"fmt"
	"strings"
)

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
