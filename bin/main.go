package main

import (
	"flag"
	"log"

	"github.com/jjyr/ringo"
)

func main() {
	flag.Parse()
	app := ringo.NewApp()
	r := app.GetRouter()
	r.GET("/ping", func(c *ringo.Context) {
		log.Print("pong!")
		c.Render(200, "pong!")
	})

	r.GET("/numbers/:n/echo", func(c *ringo.Context) {
		log.Print(c.URL.Query())
		log.Print(c.Params.Get("n"))
	})

	app.Run("localhost:8020")
}
