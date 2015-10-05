package main

import (
	"flag"
	"log"

	"github.com/jjyr/ringo"
	"github.com/jjyr/ringo/middleware"
)

type usersController struct {
	ringo.Controller
}

type User struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age"`
}

func (*usersController) List(c *ringo.Context) {
	c.String(200, "list users")
}

func (*usersController) Create(c *ringo.Context) {
	u := User{}
	if err := c.Bind(&u); err != nil {
		panic(err)
	}
	resp := struct {
		UserInfo *User
		Message  string `json:message`
	}{UserInfo: &u, Message: "hello"}
	c.JSON(200, resp)
}

func (*usersController) SayHehe(c *ringo.Context) {
	c.JSON(200, struct {
		Hello   string
		Message string
	}{Hello: "nihaoo",
		Message: "你好",
	})
}

func main() {
	flag.Parse()
	app := ringo.NewApp()
	r := ringo.NewRouter()
	r.GET("/hello", func(c *ringo.Context) {
		c.String(200, "hello world!")
	})

	app.Mount("/say", r)

	app.GET("/ping", func(c *ringo.Context) {
		log.Print("pong!")
		c.String(200, "pong!")
	})

	app.GET("/numbers/:n/echo/:n", func(c *ringo.Context) {
		log.Print(c.Request.URL.Query())
		log.Print(c.Params.ByName("n"))
	})

	u := usersController{}
	u.AddRoutes(ringo.ControllerRouteOption{Handler: "SayHehe", Method: "GET", Member: true, Path: "hehe"})
	app.AddController(&u, nil)

	app.Use(middleware.Recover())

	app.Run("localhost:8020")
}
