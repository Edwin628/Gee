package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(c *Context)
type DataStruct map[string]interface{}

type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
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

func (routergroup *RouterGroup) createStaticHandle(relativepath string, fs http.FileSystem) HandlerFunc {
	prefix := path.Join(routergroup.prefix, relativepath)
	fileServer := http.StripPrefix(prefix, http.FileServer(fs))
	return func(c *Context) {
		file := c.Params["filepath"]
		log.Printf("assets file: %s", file)
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.StatusCode = http.StatusNotFound
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (routergroup *RouterGroup) Static(relativepath string, root string) {
	handle := routergroup.createStaticHandle(relativepath, http.Dir(root))
	urlPattern := path.Join(relativepath, "/*filepath")
	routergroup.GET(urlPattern, handle)
}

// Use defines the middleware added
func (routergroup *RouterGroup) Use(handlers ...HandlerFunc) {
	routergroup.middleware = append(routergroup.middleware, handlers...)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, engine)
	return err
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	var middleware []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middleware = append(middleware, group.middleware...)
		}
	}
	ctx.Handles = middleware
	ctx.Engine = engine
	engine.router.handle(ctx)
}
