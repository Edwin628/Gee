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
	result := reflect.DeepEqual(parsePattern("/"), []string{})
	if !result {
		t.Errorf("Failed case: /")
	}

	result = reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	result = result && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	result = result && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})

	if !result {
		t.Errorf("Failed TestParsePattern")
	}
}

func TestGetRoute(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		pattern   string
		paramName string
		param     string
	}{
		{"test :var", "/hello/aibin", "/hello/:name", "name", "aibin"},
		{"test path", "/hello/b/c", "/hello/b/c", "", ""},
		{"test *var", "/assets/aibin", "/assets/*filepath", "filepath", "aibin"},
		{"test *var", "/assets/aibin/x", "/assets/*filepath", "filepath", "aibin/x"},
	}
	r := newTestRouter()

	for _, tt := range tests {
		// sequence one by one
		t.Run(tt.name, func(t *testing.T) {
			n, params := r.getRoute("GET", tt.path)
			if n == nil {
				t.Errorf("node nil happened")
			}

			if !reflect.DeepEqual(n.pattern, tt.pattern) {
				t.Errorf("parsePattern(%q) got %v, want %v", tt.pattern, n.pattern, tt.pattern)
			}
			if tt.paramName != "" {
				if !reflect.DeepEqual(params[tt.paramName], tt.param) {
					t.Errorf("Params name not matched")
				}
			}
		})
	}

}
