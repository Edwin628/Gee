package gee

type node struct {
	pattern string
	next    map[string]*node
	isEnd   bool
}

func NewNode() *node {
	return &node{
		pattern: "",
		next:    make(map[string]*node),
		isEnd:   false,
	}
}

func (n *node) insert(strs []string) {
	cur := n
	for _, str := range strs {
		if cur.next[str] == nil {
			cur.next[str] = NewNode()
		}
		cur = cur.next[str]
	}
	cur.isEnd = true
}

func (n *node) search(strs []string) bool {
	cur := n
	for _, str := range strs {
		if cur.next[str] != nil {
			cur = cur.next[str]
		} else {
			return false
		}
	}
	return cur.isEnd
}
