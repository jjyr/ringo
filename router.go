package ringo

import (
	"fmt"

	"github.com/golang/glog"
)

// Router ringo router interface
type Router interface {
	AddRoute(path string, method string, handler HandleFunc)
}

func GET(r Router) (path string, handler HandleFunc) {
	r.AddRoute(path, "GET", handler)
}
func POST(r Router) (path string, handler HandleFunc) {
	r.AddRoute(path, "POST", handler)
}
func PUT(r Router) (path string, handler HandleFunc) {
	r.AddRoute(path, "PUT", handler)
}
func DELETE(r Router) (path string, handler HandleFunc) {
	r.AddRoute(path, "DELETE", handler)
}
func HEAD(r Router) (path string, handler HandleFunc) {
	r.AddRoute(path, "HEAD", handler)
}
func OPTION(r Router) (path string, handler HandleFunc) {
	r.AddRoute(path, "OPTION", handler)
}
func PATCH(r Router) (path string, handlers HandleFunc) {
	r.AddRoute(path, "PATCH", handler)
}

// Mount(path string, r Router)

// DefaultRouter default Router implement
type DefaultRouter struct {
	routes map[string]HandleFunc
}

func AddRoute(r *DefaultRouter) (path string, method string, handler HandleFunc) {
	key := fmt.Sprint("%s/%s", method, path)
	if _, ok = r.routes[key]; ok {
		glog.Fatalf("[%s]%s already exists, duplicate route!", method, path)
	} else {
		r.routes[key] = handler
	}
}

// NewRouter make a default router
func NewRouter() Router {
	r := DefaltRouter{}
	r.routes = make(map[string]HandleFunc)
	return &r
}
