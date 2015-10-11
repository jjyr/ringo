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

type anotherUsersController struct {
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
		if n := getControllerName(cn.c, nil); n != cn.name {
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
		{"POST", "/user/2/custome", "MemberCustome", "2"},
		{"POST", "/users-customize/custome", "CollectionCustome", ""},
	}

	r := NewApp()
	users := &usersController{}
	users.AddRoutes(
		ControllerRouteOption{Member: true, Path: "custome", Method: "POST", Handler: "MemberCustome"},
		ControllerRouteOption{Path: "custome", Suffix: "-customize", Method: "POST", Handler: "CollectionCustome"},
	)
	r.AddController(users, nil)
	r.AddController(users, &ControllerOption{Name: "user"})
	context := NewContext()
	for i, c := range cases {
		id = "nil"
		actionName = "nil"
		h, _, _ := r.MatchRoute(c.path, c.method)
		if h != nil {
			h(context)
		}
		if h == nil || actionName != c.handler || id != c.id {
			t.Errorf("Test case %d failed, expect action: %s, id: %v; get handler %s, id: %v", i+1, c.handler, c.id, actionName, id)
		}
	}

	// find controller route path case
	pathCases := []struct {
		controller, handler, path string
		paramvalue                []string
	}{
		{"users", "List", "/users", nil},
		{"users", "Get", "/users/1", []string{"id", "1"}},
		{"users", "Edit", "/users/editor/edit", []string{"id", "editor"}},
		{"users", "Edit", "/users/3/edit", []string{"id", "3"}},
		{"users", "New", "/new-users", nil},
		{"users", "Create", "/users", nil},
		{"users", "Update", "/users/lastman", []string{"id", "lastman"}},
		{"users", "Delete", "/users/2", []string{"id", "2"}},
		{"users", "MemberCustome", "/users/2/custome", []string{"id", "2"}},
		{"users", "CollectionCustome", "/users-customize/custome", nil},
		{"user", "CollectionCustome", "/user-customize/custome", nil},
	}

	for i, c := range pathCases {
		routePath := r.FindControllerRoutePath(c.controller, c.handler, c.paramvalue...)
		if routePath != c.path {
			t.Errorf("Test case %d failed, find route path '%s#%s', expect '%s' get '%s'", i+1, c.controller, c.handler, c.path, routePath)
		}
	}

	func() {
		defer func() { recover() }()
		r.FindControllerRoutePath("void", "nothing")
		t.Errorf("Find no exists route path should panic")
	}()

	func() {
		defer func() { recover() }()
		users := &anotherUsersController{}
		r.AddController(users, &ControllerOption{Name: "users"})
		t.Errorf("Add same name controller should panic")
	}()
}
