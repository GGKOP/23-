package gei

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
	*RouterGroup
	groups []*RouterGroup
}

type RouterGroup struct {
	parent   *RouterGroup
	middware []HandlerFunc
	name     string
	engine   *Engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (routergroup *RouterGroup) NewRouterGroup(prefix string) *RouterGroup {
	engine := routergroup.engine
	newGroup := &RouterGroup{
		name:   routergroup.name + prefix,
		parent: routergroup,
		engine: routergroup.engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (routergroup *RouterGroup) RunMiddware(midware ...HandlerFunc) {
	routergroup.middware = append(routergroup.middware, midware...)
}

func (routergroup *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := routergroup.name + comp
	routergroup.engine.router.addRoute(method, pattern, handler)
}

func (routergroup *RouterGroup) GET(pattern string, handler HandlerFunc) {
	routergroup.addRoute("GET", pattern, handler)
}

func (routergroup *RouterGroup) POST(pattern string, handler HandlerFunc) {
	routergroup.addRoute("POST", pattern, handler)
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middwares []HandlerFunc
	for _, routergroup := range engine.groups {
		if strings.HasPrefix(req.URL.Path, routergroup.name) {
			middwares = append(middwares, routergroup.middware...)
		}
	}
	c := newContext(w, req)
	c.handlers = middwares
	engine.router.handle(c)
}
