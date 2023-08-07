package gei

import (
	"fmt"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']  目标是搞成这个样子。

// 构造一个新的router
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/") //将vs用‘/’分割开来得到的字符串切片给vs

	parts := make([]string, 0) //定义一个parts  如果有一个* 那么就返回parts 如果没有*  那么就parts和vs一样。
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//如何添加路由 ：要有方法 ， 路径 ，以及handlerfunc，方法直接在context中直接获取 路径就是 parts 然后handerfunc是我们需要定义的部分。
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern) //筛掉有没有多余的* 并且只截取*之前的部分
	key := method + "-" + pattern  // key拼成 GET-/454545这种类型的
	_, ok := r.roots[method]       //将 method找出对应的树并且赋值给ok
	if !ok {
		r.roots[method] = &node{} //如果为空  那么需要根据 method 生成一个root
	}
	r.roots[method].insert(pattern, parts, 0) // 插入 pattern parts
	r.handlers[key] = handler

}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 先筛*
	params := make(map[string]string) //
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0) //按照筛完*的标准去前缀树里找node

	if n != nil {
		parts := parsePattern(n.pattern) //将找的node里面的pattern中继续拆以用来存储参数
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	fmt.Println("Error reading file:", key)
	if handler, ok := r.handlers[key]; ok {
		handler(c)
		return
	}

}
