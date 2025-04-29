package main

import (
	"fmt"
	"log"
	"lycoris"
	"net/http"
)

func main() {
	// 1. 创建框架实例
	r := lycoris.New()

	// 2. 添加全局中间件
	r.Use(func(c *lycoris.Context) {
		fmt.Println("全局中间件: 请求开始")
		c.Next()
		fmt.Println("全局中间件: 请求结束")
	})

	// 3. 创建API路由组
	api := r.Group("/api", func(c *lycoris.Context) {
		fmt.Println("API分组中间件: 验证Token")
		c.Next()
	})

	// 4. 添加用户路由组
	user := api.Group("/user", func(c *lycoris.Context) {
		fmt.Println("用户分组中间件: 记录访问日志")
		c.Next()
	})

	// 5. 注册路由
	user.GET("/:id", func(c *lycoris.Context) {
		// 获取路由参数
		id := c.GetParam("id")
		// 返回JSON响应
		c.JSON(http.StatusOK, lycoris.H{
			"code":    0,
			"message": "success",
			"data": lycoris.H{
				"id":   id,
				"name": "张三",
				"age":  25,
			},
		})
	})

	user.POST("/create", func(c *lycoris.Context) {
		// 获取POST表单数据
		name := c.PostForm("name")
		age := c.PostForm("age")

		// 返回创建成功的响应
		c.JSON(http.StatusOK, lycoris.H{
			"code":    0,
			"message": "用户创建成功",
			"data": lycoris.H{
				"name": name,
				"age":  age,
			},
		})
	})

	// 6. 静态文件路由
	r.GET("/static/*filepath", func(c *lycoris.Context) {
		filepath := c.GetParam("filepath")
		c.String(http.StatusOK, "请求的静态文件是: %s", filepath)
	})

	// 7. 启动服务器
	fmt.Println("服务器启动在 http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
