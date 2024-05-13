package main

import (
	"net/http"

	"gee"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	// curl "http://localhost:9999/hello?name=geektutu"
	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Request.URL.Query().Get("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Params["name"], c.Path)
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(gee.DataStruct{
			"filepath": c.Params["filepath"],
		})
	})

	// curl "http://localhost:9999/login" -X POST -d 'username=geektutu&password=1234'
	r.POST("/login", func(c *gee.Context) {
		c.JSON(gee.DataStruct{
			"username": c.Request.FormValue("username"),
			"password": c.Request.FormValue("password"),
		})
	})

	r.Run(":9999")
}
