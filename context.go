package ringo

import (
	"net/http"

	"github.com/jjyr/ringo/binding"
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

func (c *Context) ContentType() string {
	return c.Request.Header.Get("Content-Type")
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

func (c *Context) Redirect(code int, location string) {
	// make rendered check correct
	c.Render(code, &render.Redirect{
		Code:     code,
		Location: location,
		Request:  c.Request,
	})
}

func (c *Context) JSON(statusCode int, content interface{}) {
	c.Render(statusCode, &render.JSONData{Content: content})
}

// Bind
func (c *Context) Bind(obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.BindWith(obj, b)
}

func (c *Context) BindWith(obj interface{}, b binding.Binder) error {
	return binding.BindWith(c.Request, obj, b)
}
