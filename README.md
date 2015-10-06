# Ringo [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/jjyr/ringo) [![Build Status](https://travis-ci.org/jjyr/ringo.svg?branch=master)](https://travis-ci.org/jjyr/ringo)

Yet another MVC web framework for Go, inspired from rails, gin.

## Install

`go get github.com/jjyr/ringo`

## TODO
- [ ] Rails-like generator.
- [ ] Templates support.
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
package main

import (
	"fmt"

	"github.com/jjyr/ringo"
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
	c.JSON(200, ctl.users)
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
		c.String(200, "ok")
	} else {
		c.String(404, "not found")
	}
}

func (ctl *UserController) Create(c *ringo.Context) {
	u := user{}
	if err := c.BindJSON(&u); err == nil {
		ctl.users = append(ctl.users, u)
		c.JSON(200, u)
	} else {
		c.String(400, "format error")
	}
}

// customized action
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
		c.String(404, "not found")
	} else {
		c.String(200, fmt.Sprintf("Hello, I'm %s", u.Name))
	}
}

var userController *UserController

// init userController
func init() {
	userController = &UserController{}
	// register customized route
	// Post -> [POST]/user/:id/greeting
	userController.AddRoutes(ringo.ControllerRouteOption{Handler: "Greet", Member: true, Method: "POST", Path: "greeting"})
}

func main() {
	app := ringo.NewApp()
	app.AddController(userController, nil)
	app.Run("localhost:8000")
}
```


## Contribute

* Feel free to open issue
* Feel free to send pull request
* New feature pull request should include test case
