package main

import "github.com/jjyr/ringo"

func main() {
	app := ringo.NewApp()
	app.GET("/", func(c ringo.Context) {
		c.String(200, "hello")
	})
	app.Run("localhost:8000")
}

