package context

import (
	"net/http"
	"github.com/jjyr/ringo/binding"
	"github.com/jjyr/ringo/render"
	"github.com/jjyr/ringo/common"
)

// Context Application context during requests
type Context struct {
	request        *http.Request
	http.ResponseWriter
	params         common.Params
	metadata       map[string]interface{}
	TemplateManage TemplateManage
}

type TemplateManage interface {
	FindTemplate(string) *render.Template
}

// check context.Context interface
var _ common.Context = &Context{}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) ContentType() string {
	return c.request.Header.Get("Content-Type")
}

func (c *Context) Request() (*http.Request) {
	return c.request
}

func (c *Context) SetRequest(request *http.Request) {
	c.request = request
}

func (c *Context) Params()(*common.Params) {
	return &c.params
}

func (c *Context) SetParams(params common.Params) {
	c.params = params
}

// Set key value pair
func (c *Context) Set(key string, value interface{}) {
	if c.metadata == nil {
		c.metadata = make(map[string]interface{})
	}
	c.metadata[key] = value
}

// Get value by key
func (c *Context) Get(key string) (value interface{}, exists bool) {
	if c.metadata != nil {
		value, exists = c.metadata[key]
	}
	return
}

// render
func (c *Context) String(statusCode int, format string, contents ...interface{}) {
	c.Render(statusCode, &render.TextData{Format: format, Contents: contents})
}

// Render render Renderable object to client
func (c *Context) Render(statusCode int, r render.Renderable) {
	c.WriteHeader(statusCode)
	if err := r.Render(c.ResponseWriter); err != nil {
		panic(err)
	}
	c.ResponseWriter.(*ResponseWriter).Flush()
}

func (c *Context) HasRendered() bool {
	return c.ResponseWriter.(*ResponseWriter).Written()
}

func (c *Context) Redirect(code int, location string) {
	// make rendered check correct
	c.Render(code, &render.Redirect{
		Code:     code,
		Location: location,
		Request:  c.request,
	})
}

func (c *Context) JSON(statusCode int, content interface{}) {
	c.Render(statusCode, &render.JSONData{Content: content})
}

func (c *Context) File(file string) {
	http.ServeFile(c.ResponseWriter, c.request, file)
}

func (c *Context) HTML(statusCode int, name string, data interface{}) {
	c.Render(statusCode, &render.HTMLTemplateData{Data: data, Template: c.TemplateManage.FindTemplate(name)})
}

// Bind bind object to request
func (c *Context) Bind(obj interface{}) error {
	b := binding.Default(c.request.Method, c.ContentType())
	return c.BindWith(obj, b)
}

// BindWith bind Binder object to requests
func (c *Context) BindWith(obj interface{}, b binding.Binder) error {
	return binding.BindWith(c.request, obj, b)
}

func (c *Context) BindJSON(obj interface{}) error {
	return c.BindWith(obj, binding.JsonBinding)
}
