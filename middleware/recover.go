package middleware

import (
	"log"
	"net/http"

	"github.com/jjyr/ringo/common"
)

func Recover() common.MiddlewareFunc {
	return func(handler common.HandlerFunc) common.HandlerFunc {
		return func(c common.Context) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("error: %s", r)
					if !c.HasRendered() {
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
