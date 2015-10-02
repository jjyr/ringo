package ringo

import (
	"fmt"
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

func (c *Context) String(statusCode int, content interface{}) {
	w := c.ResponseWriter.(*ResponseWriter)
	var byteContent []byte

	switch content.(type) {
	case []byte:
		byteContent = content.([]byte)
	default:
		byteContent = []byte(fmt.Sprint(content))
	}
	w.WriteHeader(statusCode)
	w.Write(byteContent)
}

func (c *Context) Rendered() bool {
	return c.ResponseWriter.(*ResponseWriter).Written()
}

func (c *Context) JSON(statusCode int, content interface{}) {
	render.JSON(c.ResponseWriter, statusCode, content)
}
