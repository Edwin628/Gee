package gee

import "net/http"

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	// request info
	Method string
	Path   string
	// response info
	StatusCode int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Method:  r.Method,
		Path:    r.URL.Path,
	}
}
