package lycoris

import "net/http"

// RouterGroup 路由分组结构
type RouterGroup struct {
	prefix      string        // 分组前缀
	middlewares []HandlerFunc // 分组级中间件
	parent      *RouterGroup  // 父分组
	engine      *Engine       // 关联的Engine实例
}

// Engine 框架核心结构
type Engine struct {
	*RouterGroup                // 默认路由组
	router       *router        // 路由核心
	groups       []*RouterGroup // 所有路由组
}

// New 创建Lycoris框架实例
func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	// 创建默认路由组
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group 创建新的路由分组
// prefix: 分组前缀
// middlewares: 分组级中间件
func (group *RouterGroup) Group(prefix string, middlewares ...HandlerFunc) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix, // 组合完整前缀
		middlewares: make([]HandlerFunc, 0),
		parent:      group,
		engine:      engine,
	}
	newGroup.Use(middlewares...)
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use 添加中间件到分组
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// addRoute 内部路由注册方法
func (group *RouterGroup) addRoute(method, pattern string, handler HandlerFunc) {
	// 组合完整路径
	fullPattern := group.prefix + pattern
	// 组合处理器链（包含中间件和路由处理器）
	handlers := group.combineHandlers(handler)
	// 注册到路由表
	group.engine.router.addRoute(method, fullPattern, handlers)
}

// combineHandlers 组合中间件和路由处理器
func (group *RouterGroup) combineHandlers(handler HandlerFunc) []HandlerFunc {
	// 计算处理器链总长度
	finalSize := len(group.middlewares) + 1
	if parent := group.parent; parent != nil {
		finalSize += len(parent.middlewares)
	}

	// 预分配足够空间的切片
	finalHandlers := make([]HandlerFunc, 0, finalSize)

	// 添加父分组的中间件
	if parent := group.parent; parent != nil {
		finalHandlers = append(finalHandlers, parent.middlewares...)
	}

	// 添加当前分组的中间件
	finalHandlers = append(finalHandlers, group.middlewares...)

	// 添加路由处理器
	finalHandlers = append(finalHandlers, handler)

	return finalHandlers
}

// GET 注册GET路由
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 注册POST路由
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run 启动HTTP服务器
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// ServeHTTP 实现http.Handler接口
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}
