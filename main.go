package main

import (
	"lycoris"
	"net/http"
)

func main() {
	r := lycoris.New()

	r.GET("/", func(c *lycoris.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World!</h1>")
	})

	r.GET("/hello/:name", func(c *lycoris.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s.\n", c.Params["name"], c.Path)
	})

	r.Run(":8080")
}
