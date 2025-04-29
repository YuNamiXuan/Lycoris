package lycoris

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	parts := make([]string, 0)
	for _, item := range strings.Split(pattern, "/") {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	if _, exist := r.roots[method]; !exist {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)

	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	parts := parsePattern(pattern)

	params := make(map[string]string)
	root, exist := r.roots[method]
	if !exist {
		return nil, nil
	}

	n := root.search(parts, 0)
	if n != nil {
		temp := parsePattern(n.pattern)
		for index, part := range temp {
			if part[0] == ':' {
				params[part[1:]] = parts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(parts[index:], "/")
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
