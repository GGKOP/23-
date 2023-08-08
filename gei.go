package gei

import (
	"net/http"
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

func (routergroup *RouterGroup) NewRouterGroup(prefix string, handler HandlerFunc) *RouterGroup {
	engine := routergroup.engine
	newGroup := &RouterGroup{
		name:     routergroup.name + prefix,
		parent:   routergroup,
		engine:   routergroup.engine,
		middware: []HandlerFunc{handler}, //将提前handler 注册到 分组路由里面
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
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
	c := newContext(w, req)
	engine.router.handle(c)
}
