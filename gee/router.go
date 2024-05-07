package gee

import "fmt"

type router struct {
	handles map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handles: make(map[string]HandlerFunc),
	}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
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
