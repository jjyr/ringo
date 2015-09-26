package ringo

import (
	"fmt"
	"log"
	"net/url"
	"path"
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

func (router *Router) Mount(mountPath string, mountedRouter *Router) {
	for _, r := range mountedRouter.routes {
		router.AddRoute(path.Join(mountPath, r.path), r.method, r.handler)
	}
}

func (r *Router) AddController(controller Controllerable) {
	registerToRouter(r, controller)
}

type routeHandler struct {
	method     string
	pathRegexp *regexp.Regexp
	path       string
	handler    HandlerFunc
}

type routeRegisterMap map[string]bool

type Router struct {
	routes        []routeHandler
	anyRoutes     []routeHandler
	routesPathMap routeRegisterMap
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
	key := fmt.Sprintf("%s %s", method, path)
	if _, exist := r.routesPathMap[key]; exist {
		log.Panicf("router [%s]%s already exists", method, path)
	} else {
		r.routesPathMap[key] = true
	}

	pathRegexp, err := compilePathExp(path)
	if err != nil {
		log.Panicf("Path compile error, please check syntax:\n[%s]%s\n%s", method, path, err)
	}

	route := routeHandler{path: path, method: method, handler: handler, pathRegexp: pathRegexp}

	if method == "*" {
		r.anyRoutes = append(r.anyRoutes, route)
	} else {
		r.routes = append(r.routes, route)
	}
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
	r.routesPathMap = make(routeRegisterMap)
	return &r
}
