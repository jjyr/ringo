package ringo

import "testing"

var actionName, id string

type usersController struct {
	Controller
}

type NetWorkController struct {
	Controller
}

type SECURYController struct {
	Controller
}

func (*usersController) New(c *Context) {
	actionName = "New"
	id = c.Params.ByName("id")
}
func (*usersController) Edit(c *Context) {
	actionName = "Edit"
	id = c.Params.ByName("id")
}
func (*usersController) List(c *Context) {
	actionName = "List"
	id = c.Params.ByName("id")
}
func (*usersController) Create(c *Context) {
	actionName = "Create"
	id = c.Params.ByName("id")
}
func (*usersController) Delete(c *Context) {
	actionName = "Delete"
	id = c.Params.ByName("id")
}
func (*usersController) Update(c *Context) {
	actionName = "Update"
	id = c.Params.ByName("id")
}
func (*usersController) Get(c *Context) {
	actionName = "Get"
	id = c.Params.ByName("id")
}
func (*usersController) MemberCustome(c *Context) {
	actionName = "MemberCustome"
	id = c.Params.ByName("id")
}
func (*usersController) CollectionCustome(c *Context) {
	actionName = "CollectionCustome"
	id = c.Params.ByName("id")
}

func TestController(t *testing.T) {

	// test controller name detective
	conrollerNamesMap := []struct {
		c    Controllerable
		name string
	}{
		{&usersController{}, "users"},
		{&NetWorkController{}, "net-work"},
		{&SECURYController{}, "secury"},
	}

	for _, cn := range conrollerNamesMap {
		if n := GetControllerName(cn.c); n != cn.name {
			t.Errorf("Controller name detect error, the name '%s' should be '%s'", n, cn.name)
		}
	}

	// test controller default actions
	cases := []struct {
		method, path, handler, id string
	}{
		{"GET", "/users", "List", ""},
		{"GET", "/users/tester", "Get", "tester"},
		{"GET", "/users/editor/edit", "Edit", "editor"},
		{"GET", "/new-users", "New", ""},
		{"POST", "/users", "Create", ""},
		{"PUT", "/users/1", "Update", "1"},
		{"PATCH", "/users/1", "Update", "1"},
		{"DELETE", "/users/2", "Delete", "2"},
		{"POST", "/users/2/custome", "MemberCustome", "2"},
		{"POST", "/user/custome", "CollectionCustome", ""},
	}

	r := NewRouter()
	r.AddController(&usersController{},
		ControllerRouterOption{Member: true, Path: "custome", Method: "POST", Handler: "MemberCustome"},
		ControllerRouterOption{Collection: true, Path: "custome", Name: "user", Method: "POST", Handler: "CollectionCustome"},
	)
	context := NewContext()
	for i, c := range cases {
		id = "nil"
		actionName = "nil"
		h, _, _ := r.MatchRoute(c.path, c.method)
		h(context)
		if actionName != c.handler || id != c.id {
			t.Errorf("Test case %d failed, expect action: %s, id: %v; get handler %s, id: %v", i+1, c.handler, c.id, actionName, id)
		}
	}
}
