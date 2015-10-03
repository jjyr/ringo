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

func (*usersController) List(c *ringo.Context) {
	c.String(200, "list users")
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
		log.Print(c.URL.Query())
		log.Print(c.Params.ByName("n"))
	})

	u := usersController{}

	app.AddController(&u, ringo.ControllerRouterOption{Handler: "SayHehe", Method: "GET", Member: true, Path: "hehe"})

	app.Use(middleware.Recover())

	app.Run("localhost:8020")
}
