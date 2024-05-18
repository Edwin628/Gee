package gee

import (
	"fmt"
	"strings"
)

type router struct {
	nodes   map[string]*node
	handles map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		nodes:   make(map[string]*node),
		handles: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	tokens := strings.Split(pattern, "/")

	// if we write var parts []string, it will cause errors in the case:/

	// case 1
	// var parts []string
	// result := reflect.DeepEqual(parsePattern("/"), []string(nil))

	// case 2
	parts := make([]string, 0)
	// result := reflect.DeepEqual(parsePattern("/"), []string{})
	for _, token := range tokens {
		if token != "" {
			parts = append(parts, token)
			// accept the first *
			if token[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	_, ok := r.nodes[method]
	if !ok {
		r.nodes[method] = NewNode("", false)
	}
	r.nodes[method].insert(pattern, parts, 0)
	// know the conncept of method&pattern here
	key := method + "-" + pattern
	r.handles[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	pathParts := parsePattern(path)
	n, ok := r.nodes[method]
	if !ok {
		return nil, nil
	}
	res := n.search(pathParts, 0)
	if res == nil || res.pattern == "" {
		return nil, nil
	}

	params := make(map[string]string)

	// to assign the pattern value to params
	patternParts := parsePattern(res.pattern)
	for index, part := range patternParts {
		if part[0] == ':' {
			params[part[1:]] = pathParts[index]
		}
		// TBD if only * exists, what should we do?
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(pathParts[index:], "/")
			break
		}
	}

	return res, params

}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n == nil {
		fmt.Fprintf(c.Writer, "404 Not Found, Path Nil")
		return
	}
	c.Params = params
	key := c.Method + "-" + n.pattern
	if handle, ok := r.handles[key]; !ok {
		fmt.Fprintf(c.Writer, "404 Not Found, Pattern Nil")
	} else {
		handle(c)
	}
}
