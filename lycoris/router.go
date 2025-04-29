package lycoris

import (
	"net/http"
	"strings"
)

type routeInfo struct {
	Method      string
	Path        string
	HandlerName string
}

type router struct {
	roots map[string]*node // 每种HTTP方法对应一个前缀树
}

func newRouter() *router {
	return &router{
		roots: make(map[string]*node),
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

func (r *router) addRoute(method, pattern string, handlers []HandlerFunc) {
	parts := parsePattern(pattern)

	if _, exist := r.roots[method]; !exist {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0, handlers)
}

func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern)
	params := make(map[string]string)

	root, exist := r.roots[method]
	if !exist {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = parts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(parts[index:], "/")
				break
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
		c.handlers = n.handlers
	} else {
		c.handlers = []HandlerFunc{func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		}}
	}
	c.Next()
}
