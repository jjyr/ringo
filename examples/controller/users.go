package main

import (
	"github.com/jjyr/ringo"
	"github.com/jjyr/ringo/middleware"
)

// user model
// validation details: gopkg.in/go-playground/validator.v8
type user struct {
	Name string `json:"name" validate:"required"`
}

var users []user

// controllers
type Users struct {
	ringo.Controller
}

func (_ *Users) Get(c ringo.Context) {
	c.JSON(200, ringo.H{"users": users})
}

func (ctl *Users) Post(c ringo.Context) {
	u := user{}
	if err := c.BindJSON(&u); err == nil {
		users = append(users, u)
		c.JSON(200, u)
	} else {
		c.JSON(400, ringo.H{"message": "format error"})
	}
}

func DisplayList(c ringo.Context) {
	// render html template
	c.HTML(200, "list.html", users)
}

type User struct {
	ringo.Controller
}

func (ctl *User) Delete(c ringo.Context) {
	idx := -1
	name := c.Params().ByName("id")
	for i, u := range users {
		if name == u.Name {
			// found
			idx = i
			break
		}
	}
	if idx >= 0 {
		users = append(users[:idx], users[idx+1:]...)
		c.JSON(200, ringo.H{"message": "ok"})
	} else {
		c.JSON(404, ringo.H{"message": "not found"})
	}
}

func main() {
	app := ringo.NewApp()

	// register controllers
	app.Add("/users", &Users{})
	app.GET("/users-list", DisplayList)
	app.Add("/user/:id", &User{})
	// use recover middleware, response 500 if handler panic
	app.Use(middleware.Recover())
	// setup templates path
	app.SetTemplatePath("templates")
	app.Run("localhost:8000")
}
