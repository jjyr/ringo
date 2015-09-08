package ringo

import (
	"net/http"
	"net/url"
)

type Context struct {
	*http.Request
	http.ResponseWriter
	PathParams *url.Values
}

func NewContext() *Context {
	return &Context{}
}
