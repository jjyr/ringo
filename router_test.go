package ringo

import (
	"net/url"
	"reflect"
	"testing"
)

func TestRouter(t *testing.T) {
	r := NewRouter()

	var handlerName string

	h1 := func(c *Context) {
		handlerName = "h1"
	}
	h2 := func(c *Context) {
		handlerName = "h2"
	}
	h3 := func(c *Context) {
		handlerName = "h3"
	}
	h4 := func(c *Context) {
		handlerName = "h4"
	}

	equalHandler := func(h1 HandlerFunc, h2 HandlerFunc) bool {
		if h1 == nil && h2 == nil {
			return true
		} else if h1 == nil || h2 == nil {
			return false
		}

		h1(nil)
		h1Name := handlerName
		h2(nil)
		h2Name := handlerName
		return h1Name == h2Name
	}

	r.GET("/tests", h1)
	r.POST("/tests", h2)
	r.HEAD("/tests", h3)
	r.Any("/tests", h4)
	r.Any("/tests/any", h4)
	r.DELETE("/tests", h2)
	r.PATCH("/change/:thing", h2)
	r.PUT("/get/:thing/info", h2)
	r.OPTIONS("/try/:two/params/:togather", h1)
	r.OPTIONS("/same/:name/params/:name", h1)

	r2 := NewRouter()
	r2.GET("/hello", h3)
	r2.POST("/echo/:word", h4)
	r.Mount("/mount", r2)

	cases := []struct {
		method, path string
		handler      HandlerFunc
		params       *url.Values
	}{
		{"GET", "/tests", h1, &url.Values{}},
		{"GET", "tests", nil, nil},
		{"GET", "tests/", nil, nil},
		{"POST", "/tests", h2, &url.Values{}},
		{"OPTIONS", "/tests", h4, &url.Values{}},
		{"HEAD", "/tests", h3, &url.Values{}},
		{"DELETE", "/tests", h2, &url.Values{}},
		{"PATCH", "/change", nil, nil},
		{"PATCH", "/change/", nil, nil},
		{"PATCH", "/change/world", h2, &url.Values{"thing": []string{"world"}}},
		{"PUT", "/get/secret/info", h2, &url.Values{"thing": []string{"secret"}}},
		{"OPTIONS", "/try/2/params/3", h1, &url.Values{"two": []string{"2"}, "togather": []string{"3"}}},
		{"OPTIONS", "/same/test/params/test2", h1, &url.Values{"name": []string{"test2"}}},
		{"HEAD", "/try/2/params/3", nil, nil},
		{"OPTIONS", "/tests/any", h4, &url.Values{}},
		{"GET", "/tests/any", h4, &url.Values{}},
		{"GET", "/mount/hello", h3, &url.Values{}},
		{"POST", "/mount/echo/nihao", h4, &url.Values{"word": []string{"nihao"}}},
	}

	for i, c := range cases {
		h, v := r.MatchRoute(c.path, c.method)
		if !equalHandler(c.handler, h) || !reflect.DeepEqual(c.params, v) {
			t.Errorf("Test case %d failed, expect handler: %v, params: %v; get handler %v, params: %v", i+1, c.handler, c.params, h, v)
		}
	}

	func() {
		defer func() { recover() }()
		r.GET("/tests", h1)
		t.Errorf("Same method & path router should trigger panic, bug not")
	}()
}
