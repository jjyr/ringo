package ringo

import (
	"net/url"
	"reflect"
	"testing"
)

func TestRouter(t *testing.T) {
	// TODO add same name params test
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
		{"PUT", "/get/secret/info", h2, &url.Values{"secret": []string{"secret"}}},
		{"OPTIONS", "/try/2/params/3", h1, &url.Values{"two": []string{"2"}, "togather": []string{"3"}}},
		{"HEAD", "/try/2/params/3", nil, nil},
		{"OPTIONS", "/tests/any", h4, &url.Values{}},
		{"GET", "/tests/any", h4, &url.Values{}},
	}

	for i, c := range cases {
		h, v := r.MatchRoute(c.path, c.method)
		if !equalHandler(c.handler, h) || !reflect.DeepEqual(c.params, v) {
			t.Errorf("Test case %d failed, expect handler: %v, params: %v; get handler %v, params: %v", i+1, c.handler, c.params, h, v)
		}
	}
}
