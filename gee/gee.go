package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(c *Context)
type DataStruct map[string]interface{}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix     string
	middleware []HandlerFunc
	parent     *RouterGroup
	engine     *Engine
}

// New is the constructor of gee.Engine
func New() *Engine {
	eng := &Engine{
		router: newRouter(),
	}
	eng.RouterGroup = &RouterGroup{
		engine: eng,
		prefix: "",
	}
	eng.groups = append(eng.groups, eng.RouterGroup)
	return eng
}

func (routergroup *RouterGroup) Group(pre string) *RouterGroup {
	group := &RouterGroup{
		prefix: routergroup.prefix + pre,
		parent: routergroup,
		engine: routergroup.engine,
	}
	group.engine.groups = append(group.engine.groups, group)
	return group
}

func (routergroup *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	// If user write the "/path" in a casual way
	pattern = routergroup.prefix + pattern
	log.Printf("Route %4s - %s", method, pattern)
	routergroup.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (routergroup *RouterGroup) GET(pattern string, handler HandlerFunc) {
	routergroup.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (routergroup *RouterGroup) POST(pattern string, handler HandlerFunc) {
	routergroup.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, engine)
	return err
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	engine.router.handle(ctx)
}
