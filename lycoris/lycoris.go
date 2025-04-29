package lycoris

// H 是通用的map类型别名，用于方便的构建JSON响应
type H map[string]interface{}

// HandlerFunc 定义请求处理器类型
type HandlerFunc func(*Context)

// 框架的其他全局功能可以在这里添加
// 例如默认中间件、全局配置等
