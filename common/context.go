package common

import "net/http"

type Context interface {
	http.ResponseWriter
	SetParams(Params)
	Params() *Params
	SetRequest(*http.Request)
	Request() *http.Request

	// rendering
	Redirect(int, string)
	HasRendered() bool
	String(statusCode int, format string, contents ...interface{})
	JSON(statusCode int, content interface{})
	HTML(statusCode int, name string, data interface{})

	// binding
	BindJSON(interface{}) error
}
