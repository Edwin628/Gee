package gee

import (
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	result := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	result = result && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	result = result && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})

	if !result {
		t.Errorf("Failed TestParsePattern")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, params := r.getRoute("GET", "/hello/aibin")

	if n == nil {
		t.Errorf("node nil happened")
	}

	if n.pattern != "/hello/:name" {
		t.Errorf("node pattern not mattched")
	}

	if !(params["name"] == "aibin") {
		t.Errorf("Params name not matched")
	}
}
