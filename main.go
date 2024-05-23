package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
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

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

type student struct {
	Name string
	Age  int8
}

func main() {
	r := gee.New()

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	r.Use(Logger())

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css", nil)
	})

	stu1 := &student{Name: "aibin", Age: 24}
	stu2 := &student{Name: "Jack", Age: 22}

	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr", gee.DataStruct{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func", gee.DataStruct{
			"title": "gee",
			"now":   time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
		})
	})

	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css", nil)
	})

	// curl "http://localhost:9999/hello?name=geektutu"
	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Request.URL.Query().Get("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Params["name"], c.Path)
	})

	/* duplicate key in router will cause error
	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(gee.DataStruct{
			"filepath": c.Params["filepath"],
		})
	})
	*/

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
		c.HTML(http.StatusOK, "css.tmpl", nil)
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
