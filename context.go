package ringo

import (
	"net/http"
	"time"

	"github.com/jjyr/ringo/binding"
	"github.com/jjyr/ringo/render"
	"golang.org/x/net/context"
)

// Context Application context during requests
type Context struct {
	Request *http.Request
	http.ResponseWriter
	Params   Params
	metadata map[string]interface{}
	app      *App
}

// check context.Context interface
var _ context.Context = &Context{}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) ContentType() string {
	return c.Request.Header.Get("Content-Type")
}

// FindControllerRoutePath Find controller action path, return first path if multi route registered, panic if not found
func (c *Context) FindControllerRoutePath(name string, handler string) string {
	return c.app.FindControllerRoutePath(name, handler)
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

func (c *Context) Rendered() bool {
	return c.ResponseWriter.(*ResponseWriter).Written()
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

func (c *Context) File(file string) {
	http.ServeFile(c.ResponseWriter, c.Request, file)
}

func (c *Context) HTML(statusCode int, name string, data interface{}) {
	c.Render(statusCode, &render.HTMLTemplateData{Data: data, Template: c.app.TemplateManage.FindTemplate(name)})
}

// Bind bind object to request
func (c *Context) Bind(obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.BindWith(obj, b)
}

// BindWith bind Binder object to requests
func (c *Context) BindWith(obj interface{}, b binding.Binder) error {
	return binding.BindWith(c.Request, obj, b)
}

func (c *Context) BindJSON(obj interface{}) error {
	return c.BindWith(obj, binding.JsonBinding)
}

/************************************/
/*implement golang.org/x/net/context*/
/************************************/

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *Context) Done() <-chan struct{} {
	return nil
}

func (c *Context) Err() error {
	return nil
}

func (c *Context) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}
