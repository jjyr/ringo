package route

import (
	"testing"
	"github.com/jjyr/ringo/common"
	"github.com/jjyr/ringo/context"
)

func TestRouter(t *testing.T) {
	r := NewRouter()

	var handlerName string

	h1 := func(c common.Context) {
		handlerName = "h1"
	}
	h2 := func(c common.Context) {
		handlerName = "h2"
	}
	h3 := func(c common.Context) {
		handlerName = "h3"
	}
	h4 := func(c common.Context) {
		handlerName = "h4"
	}

	var h1Name, h2Name string

	equalHandler := func(h1 common.HandlerFunc, h2 common.HandlerFunc) bool {
		h1Name = "nil"
		h2Name = "nil"
		if h1 == nil && h2 == nil {
			return true
		} else if h1 == nil || h2 == nil {
			return false
		}

		c := context.NewContext()

		handlerName = ""
		h1(c)
		h1Name = handlerName

		handlerName = ""
		h2(c)
		h2Name = handlerName
		return h1Name == h2Name
	}

	paramEqual := func(p1 common.Params, p2 common.Params) bool {
		for i, p := range p1 {
			if p2[i] != p {
				return false
			}
		}
		return true
	}

	r.GET("/tests", h1)
	r.POST("/tests", h2)
	r.HEAD("/tests", h3)
	r.DELETE("/tests", h2)
	r.PATCH("/change/:thing", h2)
	r.PUT("/get/:thing/info", h2)
	r.OPTIONS("/try/:two/params/:togather", h1)
	r.OPTIONS("/same/:name/params/:name", h1)

	r2 := NewRouter()
	r2.GET("/hello", h3)
	r2.POST("/echo/:word", h4)
	r.Mount("/mount", r2)

	r3 := NewRouter()
	r3.GET("/root", h2)
	r.Mount("/", r3)

	cases := []struct {
		method, path string
		handler      common.HandlerFunc
		params       common.Params
	}{
		{"GET", "/tests", h1, common.Params{}},
		{"GET", "tests", nil, nil},
		{"GET", "tests/", nil, nil},
		{"POST", "/tests", h2, common.Params{}},
		{"HEAD", "/tests", h3, common.Params{}},
		{"DELETE", "/tests", h2, common.Params{}},
		{"PATCH", "/change", nil, nil},
		{"PATCH", "/change/", nil, nil},
		{"PATCH", "/change/world", h2, common.Params{{"thing", "world"}}},
		{"PUT", "/get/secret/info", h2, common.Params{{"thing", "secret"}}},
		{"OPTIONS", "/try/2/params/3", h1, common.Params{{"two", "2"}, {"togather", "3"}}},
		{"OPTIONS", "/same/test/params/test2", h1, common.Params{{"name", "test"}, {"name", "test2"}}},
		{"HEAD", "/try/2/params/3", nil, nil},
		{"GET", "/mount/hello", h3, common.Params{}},
		{"POST", "/mount/echo/nihao", h4, common.Params{{"word", "nihao"}}},
		{"GET", "/root", h2, common.Params{}},
	}

	for i, c := range cases {
		h, v, _ := r.MatchRoute(c.path, c.method)
		if !equalHandler(c.handler, h) || !paramEqual(c.params, v) {
			t.Errorf("Test case %d failed, expect handler: %s, params: %v; get handler %s, params: %v", i+1, h1Name, c.params, h2Name, v)
		}
	}

	if _, _, redirect := r.MatchRoute("/root/", "GET"); !redirect {
		t.Errorf("redirect path not work!")
	}

	func() {
		defer func() { recover() }()
		r.GET("/tests", h1)
		t.Errorf("Same method & path router should trigger panic, bug not")
	}()
}
