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
func (r *Router) OPTION(path string, handler HandlerFunc) {
	r.AddRoute(path, "OPTION", handler)
}
func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.AddRoute(path, "PATCH", handler)
}

// Mount(path string, r Router)

type routeHandler struct {
	Method     string
	PathRegexp *regexp.Regexp
	Path       string
	Handler    HandlerFunc
}

type Router struct {
	routes []routeHandler
}

var pathParamRegexp *regexp.Regexp

func init() {
	pathParamRegexp = regexp.MustCompile("\\:(\\w+)([/$])")
}

func compilePathExp(path string) (*regexp.Regexp, error) {
	escaped := regexp.QuoteMeta(path)
	log.Printf("escaped route path: %s", escaped)
	replaced := pathParamRegexp.ReplaceAllString(escaped, "(?P<$1>\\w+)$2")
	log.Printf("replaced route path: %s", replaced)
	return regexp.Compile(fmt.Sprintf("^%s$", replaced))
}

func (r *Router) AddRoute(path string, method string, handler HandlerFunc) {
	// key := fmt.Sprint("%s/%s", method, path)
	pathRegexp, err := compilePathExp(path)
	if err != nil {
		log.Fatalf("Path compile error, please check syntax:\n[%s]%s\n%s", method, path, err)
	}
	r.routes = append(r.routes, routeHandler{Path: path, Method: method, Handler: handler, PathRegexp: pathRegexp})
	// if _, ok = r.routes[key]; ok {
	// 	// log.Fatalf.Fatalf("[%s]%s already exists, duplicate route!", method, path)
	// } else {
	// 	r.routes[key] = handler
	// }
}

func (r *Router) MatchRoute(path string, method string) (HandlerFunc, *url.Values) {
	for _, route := range r.routes {
		log.Printf("try match [%s]%s by route %+v", method, path, route)
		if route.Method == method {
			subMatch := route.PathRegexp.FindAllStringSubmatch(path, -1)
			if subMatch != nil {
				subMatchNames := route.PathRegexp.SubexpNames()[1:]
				pathParams := url.Values{}
				for i, matchName := range subMatchNames {
					pathParams.Set(matchName, subMatch[i][1])
				}

				return route.Handler, &pathParams
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
