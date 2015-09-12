package main

import (
	"flag"
	"log"

	"github.com/jjyr/ringo"
	"github.com/jjyr/ringo/middleware"
)

func main() {
	flag.Parse()
	app := ringo.NewApp()
	// r := app.GetRouter()
	app.GET("/ping", func(c *ringo.Context) {
		log.Print("pong!")
		c.Render(200, "pong!")
	})

	app.GET("/numbers/:n/echo", func(c *ringo.Context) {
		log.Print(c.URL.Query())
		log.Print(c.Params.Get("n"))
	})

	app.Use(middleware.Recover())

	app.Run("localhost:8020")
}
