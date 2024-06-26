package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	// request info
	Method string
	Path   string

	Params map[string]string
	// response info
	StatusCode int
	// middleware
	Handles []HandlerFunc
	Index   int
	// engine
	Engine *Engine
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Method:  r.Method,
		Path:    r.URL.Path,
		Index:   -1,
	}
}

func (c *Context) Next() {
	c.Index++
	for ; c.Index < len(c.Handles); c.Index++ {
		c.Handles[c.Index](c)
	}
}

func (c *Context) status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) setHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.setHeader("Content-Type", "text/html")
	c.status(code)
	if err := c.Engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.status(500)
	}
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.setHeader("Content-Type", "text/plain")
	c.status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(obj interface{}) {
	c.setHeader("Content-Type", "application/json")
	c.status(500)
	if err := json.NewEncoder(c.Writer).Encode(obj); err != nil {
		http.Error(c.Writer, "Json Failed", 500)
	}
}
