package lycoris

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	// HTTP info
	Writer  http.ResponseWriter
	Request *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response info
	StatusCode int

	// middlewares info
	handlers []HandlerFunc
	index    int
	engin    *Engine
}

// newContext: 创建并初始化context实例
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Path:    r.URL.Path,
		Method:  r.Method,
		index:   -1, // 初始化为-1 表示未开始执行
	}
}

// Next: 执行处理器链中的下一个处理器 该方法允许中间件控制处理器执行流程
func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

// PostForm: 获取POST表单参数
func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

// SetHeader: 设置响应头
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// SetStatus: 设置HTTP状态码
func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// String: 返回纯文本响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON: 返回JSON格式响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data: 返回二进制数据响应
func (c *Context) Data(code int, data []byte) {
	c.SetStatus(code)
	c.Writer.Write(data)
}

// HTML: 返回HTML响应
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	c.Writer.Write([]byte(html))
}

// GetParam: 获取路由参数
func (c *Context) GetParam(key string) string {
	value, _ := c.Params[key]
	return value
}
