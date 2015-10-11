# Ringo [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/jjyr/ringo) [![Build Status](https://travis-ci.org/jjyr/ringo.svg?branch=master)](https://travis-ci.org/jjyr/ringo)

Lightweight & high customized MVC web framework for Go.

## Install

`go get github.com/jjyr/ringo`

## TODO
- [ ] Rails-like generator.
- [ ] Config file support.
- [ ] Ringo model?

## Usage

### hello world
``` go
// hello.go
package main

import "github.com/jjyr/ringo"

func main() {
	app := ringo.NewApp()
	app.GET("/", func(c *ringo.Context) {
		c.String(200, "hello")
	})
	app.Run("localhost:8000")
}

```

### controller
``` go
// users_sample.go
package main

import (
	"fmt"

	"github.com/jjyr/ringo"
	"github.com/jjyr/ringo/middleware"
)

// user model
// when use Bind(BindJSON) validate will auto perform, details: gopkg.in/go-playground/validator.v8
type user struct {
	Name string `json:"name" validate:"required"`
}

// user controller
// controller contains several built in RESTful routes will auto added if accord methods is defined
// built in routes:
// method(handler) -> route path
// New -> [GET]/new-user
// Create -> [POST]/user
// Edit -> [GET]/user/:id/edit
// Update -> [PUT/PATCH]/user/:id
// List -> [GET]/user
// Get -> [GET]/user/:id
// Delete -> [DELETE]/user/:id
type UserController struct {
	ringo.Controller
	users []user
}

func (ctl *UserController) List(c *ringo.Context) {
	c.JSON(200, ringo.H{"users": ctl.users})
}

func (ctl *UserController) Delete(c *ringo.Context) {
	idx := -1
	name := c.Params.ByName("id")
	for i, u := range ctl.users {
		if name == u.Name {
			// found
			idx = i
			break
		}
	}
	if idx >= 0 {
		ctl.users = append(ctl.users[:idx], ctl.users[idx+1:]...)
		c.JSON(200, ringo.H{"message": "ok"})
	} else {
		c.JSON(404, ringo.H{"message": "not found"})
	}
}

func (ctl *UserController) Create(c *ringo.Context) {
	u := user{}
	if err := c.BindJSON(&u); err == nil {
		ctl.users = append(ctl.users, u)
		c.JSON(200, u)
	} else {
		c.JSON(400, ringo.H{"message": "format error"})
	}
}

// customized actions
func (ctl *UserController) Greet(c *ringo.Context) {
	name := c.Params.ByName("id")
	var u *user = nil
	for _, user := range ctl.users {
		if name == user.Name {
			u = &user
			break
		}
	}
	if u == nil {
		c.JSON(400, ringo.H{"message": "not found"})
	} else {
		c.JSON(200, ringo.H{"message": fmt.Sprintf("Hello, I'm %s", u.Name)})
	}
}

func (ctl *UserController) DisplayList(c *ringo.Context) {
	// render html template
	c.HTML(200, "list.html", ctl.users)
}

var userController *UserController

// init userController
func init() {
	userController = &UserController{}
	// register customized route
	userController.AddRoutes([]ringo.ControllerRouteOption{
		// Greet -> [POST]/user/:id/greeting
		{Handler: "Greet", Member: true, Method: "POST", Path: "greeting"},
		// DisplayList -> [GET]/users/list
		{Handler: "DisplayList", Collection: true, Method: "GET", Suffix: "s", Path: "list"},
	}...)
}

func main() {
	app := ringo.NewApp()

	// register controllers
	app.AddController(userController, nil)
	// use recover middleware, response 500 if handler panic
	app.Use(middleware.Recover())
	// setup templates path
	app.SetTemplatePath("templates")
	app.Run("localhost:8000")
}
```

### template
```html
<!-- template/list.html -->
<ol>
  {{range .}}
  <li> {{.Name}} </li>
  {{end}}
</ol>
```


## Contribute

* Feel free to open issue
* Feel free to send pull request
* New feature pull request should include test case
