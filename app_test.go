package ringo

import (
	"testing"
	"github.com/jjyr/ringo/context"
)

var actionName, id string

type user struct {
	Controller
}

func (*user) Get(c Context) {
	actionName = "Get"
	id = c.Params().ByName("id")
}
func (*user) Put(c Context) {
	actionName = "Put"
	id = c.Params().ByName("id")
}
func (*user) Patch(c Context) {
	actionName = "Patch"
	id = c.Params().ByName("id")
}
func (*user) Post(c Context) {
	actionName = "Post"
	id = c.Params().ByName("id")
}
func (*user) Delete(c Context) {
	actionName = "Delete"
	id = c.Params().ByName("id")
}

type users struct {
	Controller
}

func (*users) Get(c Context) {
	actionName = "Get"
	id = c.Params().ByName("id")
}
func (*users) Put(c Context) {
	actionName = "Put"
	id = c.Params().ByName("id")
}
func (*users) Delete(c Context) {
	actionName = "Delete"
	id = c.Params().ByName("id")
}
func (*users) Post(c Context) {
	actionName = "Post"
	id = c.Params().ByName("id")
}

type anotherUsersController struct {
	Controller
}

func TestApp(t *testing.T) {

	// test controller default actions
	cases := []struct {
		method, path, handler, id string
	}{
		{"GET", "/users", "Get", ""},
		{"DELETE", "/users", "Delete", ""},
		{"POST", "/users", "Post", ""},
		{"PUT", "/user/1", "Put", "1"},
		{"PATCH", "/user/1", "Patch", "1"},
		{"DELETE", "/user/2", "Delete", "2"},
	}

	app := NewApp()
	app.Add("/user/:id", &user{})
	app.Add("/users", &users{})
	appContext := context.NewContext()
	for i, c := range cases {
		id = "nil"
		actionName = "nil"
		handler, _, _ := app.MatchRoute(c.path, c.method)
		if handler != nil {
			handler(appContext)
		}
		if handler == nil || actionName != c.handler || id != c.id {
			t.Errorf("Test case %d failed, expect action: %s, id: %v; get handler %s, id: %v", i+1, c.handler, c.id, actionName, id)
		}
	}
	//
	//func() {
	//	defer func() { recover() }()
	//	users := &anotherUsersController{}
	//	app.Add("/users", users)
	//	t.Errorf("Add same path controller should panic")
	//}()
}
