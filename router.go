package ringo

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
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

func (r *Router) Any(path string, handler HandlerFunc) {
	r.AddRoute(path, "*", handler)
}

func (router *Router) Mount(path string, mountedRouter *Router) {
	for _, r := range mountedRouter.routes {
		router.AddRoute(path+r.path, r.method, r.handler)
	}
}

// Mount(path string, r Router)

type routeHandler struct {
	method     string
	pathRegexp *regexp.Regexp
	path       string
	handler    HandlerFunc
}

type Router struct {
	routes    []routeHandler
	anyRoutes []routeHandler
}

var pathParamRegexp *regexp.Regexp

func init() {
	pathParamRegexp = regexp.MustCompile("\\:(\\w+)(/|\\z)")
}

func compilePathExp(path string) (*regexp.Regexp, error) {
	escaped := regexp.QuoteMeta(path)
	log.Printf("escaped route path: %s", escaped)
	replaced := pathParamRegexp.ReplaceAllString(escaped, "(?P<$1>\\w+)$2")
	log.Printf("replaced route path: %s", replaced)
	return regexp.Compile(fmt.Sprintf("\\A%s\\z", replaced))
}

func (r *Router) AddRoute(path string, method string, handler HandlerFunc) {
	pathRegexp, err := compilePathExp(path)
	if err != nil {
		log.Fatalf("Path compile error, please check syntax:\n[%s]%s\n%s", method, path, err)
	}

	route := routeHandler{path: path, method: method, handler: handler, pathRegexp: pathRegexp}

	if method == "*" {
		r.anyRoutes = append(r.anyRoutes, route)
	} else {
		r.routes = append(r.routes, route)
	}
	// if _, ok = r.routes[key]; ok {
	// 	// log.Fatalf.Fatalf("[%s]%s already exists, duplicate route!", method, path)
	// } else {
	// 	r.routes[key] = handler
	// }
}

func (r *Router) MatchRoute(path string, method string) (HandlerFunc, *url.Values) {
	handler, params := matchRouteHandler(path, method, r.routes)
	if handler != nil {
		return handler, params
	}
	return matchRouteHandler(path, "*", r.anyRoutes)
}

func matchRouteHandler(path string, method string, routes []routeHandler) (HandlerFunc, *url.Values) {
	for _, route := range routes {
		log.Printf("try match [%s]%s by route %+v", method, path, route)
		if route.method == method {
			subMatch := route.pathRegexp.FindAllStringSubmatch(path, -1)
			if subMatch != nil {
				subMatchNames := route.pathRegexp.SubexpNames()[1:]
				pathParams := url.Values{}
				for i, matchName := range subMatchNames {
					pathParams.Set(matchName, subMatch[0][i+1])
				}

				return route.handler, &pathParams
			}
		}
	}

	return nil, nil
}

// NewRouter new router
func NewRouter() *Router {
	r := Router{}
	return &r
}
