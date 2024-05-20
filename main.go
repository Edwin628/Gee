package main

import (
	"log"
	"net/http"
	"time"

	"gee"
)

func Logger() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(Logger())

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>index page</h1>")
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

	// test group feature
	v1 := r.Group("/v1")
	v1.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>v1 page</h1>")
	})
	v1.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Request.URL.Query().Get("name"), c.Path)
	})

	v2 := r.Group("/v2")
	v2.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Params["name"], c.Path)
	})
	v2.POST("/login", func(c *gee.Context) {
		c.JSON(gee.DataStruct{
			"username": c.Request.FormValue("name"),
			"password": c.Request.FormValue("password"),
		})
	})

	// middleware feature
	middlewareV2 := func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		// c.String(500, "Internal Server Error\n")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
	v2.Use(middlewareV2)

	r.Run(":9999")
}
