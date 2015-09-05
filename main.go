package main

import "github.com/jjyr/ringo"

func main() {
	r = ringo.NewRouter()
	r.GET("/ping", func(c *C) {
		c.String(200, "pong")
	})
}
