package gee

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

func (n *node) insert(pattern string, parts []string, height int) {
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
	if height == len(parts) {
		return n
	}

	part := parts[height]
	child := n.findChild(part)

	if child == nil {
		return nil
	}

	return child.search(parts, height+1)
}
