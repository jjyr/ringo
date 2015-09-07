package main

import "github.com/jjyr/ringo"

func main() {
	app := ringo.NewApp()
	r := app.GetRouter()
	r.GET("/ping", func(c *ringo.Context) {
		return
	})

	app.Run("localhost:8020")
}
