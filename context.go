package ringo

import (
	"fmt"
	"net/http"
)

type Context struct {
	*http.Request
	http.ResponseWriter
	Params Params
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) Render(statusCode int, content interface{}) {
	w := c.ResponseWriter.(*ResponseWriter)
	var byteContent []byte

	switch content.(type) {
	case []byte:
		byteContent = content.([]byte)
	case string:
		byteContent = []byte(content.(string))
	default:
		byteContent = []byte(fmt.Sprint(content))
	}
	w.WriteHeader(statusCode)
	w.Write(byteContent)
}

func (c *Context) Rendered() bool {
	return c.ResponseWriter.(*ResponseWriter).Written()
}
