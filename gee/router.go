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
		handles: make(map[string]HandlerFunc),
	}
}

func (r *router) parsePattern(pattern string) []string {
	tokens := strings.Split(pattern, "/")
	var parts []string
	for _, token := range tokens {
		if token != "" {
			parts = append(parts, token)
			// tbh can we delete this one
			if token[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := r.parsePattern(pattern)
	r.nodes[method].insert(parts)
	// know the conncept of method&pattern here
	key := method + "-" + pattern
	r.handles[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handle, ok := r.handles[key]; !ok {
		fmt.Fprintf(c.Writer, "404 Not Found")
	} else {
		handle(c)
	}
}
