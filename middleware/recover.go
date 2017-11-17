package middleware

import (
	"net/http"

	"github.com/jjyr/ringo/common"
	log "github.com/sirupsen/logrus"
)

func Recover() common.MiddlewareFunc {
	return func(handler common.HandlerFunc) common.HandlerFunc {
		return func(c common.Context) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("recovery error: %s", r)
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
