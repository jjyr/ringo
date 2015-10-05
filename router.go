package ringo

import (
	"log"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
)

func (r *Router) GET(path string, handler HandlerFunc) {
	r.AddRoute(path, "GET", handler)
}
func (r *Router) POST(path string, handler HandlerFunc) {
	r.AddRoute(path, "POST", handler)
}
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.AddRoute(path, "PUT", handler)
}
func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.AddRoute(path, "DELETE", handler)
}
func (r *Router) HEAD(path string, handler HandlerFunc) {
	r.AddRoute(path, "HEAD", handler)
}
func (r *Router) OPTIONS(path string, handler HandlerFunc) {
	r.AddRoute(path, "OPTIONS", handler)
}
func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.AddRoute(path, "PATCH", handler)
}

func (router *Router) Mount(mountPath string, mountedRouter *Router) {
	for _, r := range mountedRouter.routes {
		router.AddRoute(path.Join(mountPath, r.path), r.method, r.handler)
	}
}

type routeHandler struct {
	method  string
	path    string
	handler HandlerFunc
}

type routeRegisterMap map[string]bool

type Router struct {
	routes []routeHandler
	httprouter.Router
}

func (r *Router) AddRoute(path string, method string, handler HandlerFunc) {
	log.Printf("Add handler '%s' -> [%s]%s", handler, method, path)
	r.Router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		c := w.(*Context)
		c.Request = r
		c.Params = Params(params)
		handler(c)
	})

	route := routeHandler{path: path, method: method, handler: handler}
	r.routes = append(r.routes, route)
}

func (r *Router) MatchRoute(path string, method string) (handler HandlerFunc, params Params, redirect bool) {
	rawHandler, rawParams, redirect := r.Lookup(method, path)
	params = Params(rawParams)

	if rawHandler != nil {
		handler = func(c *Context) {
			rawHandler(c, c.Request, rawParams)
		}
	}

	return
}

// NewRouter new router
func NewRouter() *Router {
	r := Router{}
	return &r
}
