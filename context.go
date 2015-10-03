package ringo

import (
	"net/http"

	"github.com/jjyr/ringo/render"
)

type Context struct {
	*http.Request
	http.ResponseWriter
	Params Params
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) String(statusCode int, format string, contents ...interface{}) {
	c.Render(statusCode, &render.TextData{Format: format, Contents: contents})
}

func (c *Context) Render(statusCode int, r render.Renderable) {
	c.WriteHeader(statusCode)
	if err := r.Render(c.ResponseWriter); err != nil {
		panic(err)
	}
	c.ResponseWriter.(*ResponseWriter).Flush()
}

func (c *Context) Rendered() bool {
	return c.ResponseWriter.(*ResponseWriter).Flushed()
}

func (c *Context) JSON(statusCode int, content interface{}) {
	c.Render(statusCode, &render.JSONData{Content: content})
}
