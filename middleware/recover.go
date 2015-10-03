package middleware

import (
	"log"
	"net/http"

	"github.com/jjyr/ringo"
)

func Recover() ringo.MiddlewareFunc {
	return func(handler ringo.HandlerFunc) ringo.HandlerFunc {
		return func(c *ringo.Context) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("error: %s", r)
					if !c.Rendered() {
						statusCode := http.StatusInternalServerError
						text := http.StatusText(statusCode)
						c.String(statusCode, text)
					}
				}
			}()
			handler(c)
		}
	}
}
