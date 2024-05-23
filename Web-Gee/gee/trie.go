package gee

import "strings"

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func NewNode(str string, wild bool) *node {
	return &node{
		pattern: "",
		part:    str,
		isWild:  wild,
	}
}

func (n *node) findChild(part string) *node {
	for _, node := range n.children {
		if node.part == part || node.isWild {
			return node
		}
	}
	return nil
}

func (n *node) findChilden(part string) []*node {
	var res []*node
	for _, node := range n.children {
		if node.part == part || node.isWild {
			res = append(res, node)
		}
	}
	return res
}

func (n *node) insert(pattern string, parts []string, height int) {
	// this condition is important
	// /p/:lang/doc, set the node "doc" pattern
	if height == len(parts) {
		n.pattern = pattern
		return
	}
	part := parts[height]

	child := n.findChild(part)
	if child == nil {
		child = NewNode(part, part[0] == ':' || part[0] == '*')
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// use strings.HasPrefix(n.part, "*") instead of n.part[0](index 0 of length 0)
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
		if n.pattern != "" {
			return n
		}
		return nil
	}

	part := parts[height]
	children := n.findChilden(part)

	for _, child := range children {
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}

	}

	return nil
}
