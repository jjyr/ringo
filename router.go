package ringo

import (
	"regexp"

	"github.com/golang/glog"
)

func (r *Router) GET(path string, handler HandleFunc) {
	r.AddRoute(path, "GET", handler)
}
func (r *Router) POST(path string, handler HandleFunc) {
	r.AddRoute(path, "POST", handler)
}
func (r *Router) PUT(path string, handler HandleFunc) {
	r.AddRoute(path, "PUT", handler)
}
func (r *Router) DELETE(path string, handler HandleFunc) {
	r.AddRoute(path, "DELETE", handler)
}
func (r *Router) HEAD(path string, handler HandleFunc) {
	r.AddRoute(path, "HEAD", handler)
}
func (r *Router) OPTION(path string, handler HandleFunc) {
	r.AddRoute(path, "OPTION", handler)
}
func (r *Router) PATCH(path string, handler HandleFunc) {
	r.AddRoute(path, "PATCH", handler)
}

// Mount(path string, r Router)

type routeHandler struct {
	Method     string
	PathRegexp *regexp.Regexp
	Path       string
	Handler    HandleFunc
}

type Router struct {
	routes []routeHandler
}

var pathParamRegexp *regexp.Regexp

func init() {
	pathParamRegexp = regexp.MustCompile("(\\:w+)(?:/|$)")
}

func compilePathExp(path string) (*regexp.Regexp, error) {
	escaped := regexp.QuoteMeta(path)
	glog.Info("escaped route path: ", escaped)
	replaced := pathParamRegexp.ReplaceAllString(escaped, "(?P<$1>w+)")
	glog.Info("replaced route path: ", replaced)
	return regexp.Compile(replaced)
}

func (r *Router) AddRoute(path string, method string, handler HandleFunc) {
	// key := fmt.Sprint("%s/%s", method, path)
	pathRegexp, err := compilePathExp(path)
	if err != nil {
		glog.Fatalf("Path compile error, please check syntax:\n[%s]%s\n%s", method, path, err)
	}
	r.routes = append(r.routes, routeHandler{Path: path, Method: method, Handler: handler, PathRegexp: pathRegexp})
	// if _, ok = r.routes[key]; ok {
	// 	// glog.Fatalf("[%s]%s already exists, duplicate route!", method, path)
	// } else {
	// 	r.routes[key] = handler
	// }
}

func (r *Router) MatchRoute(path string, method string) (HandleFunc, map[string]string) {
	for _, route := range r.routes {
		if route.Method == method {
			subMatch := route.PathRegexp.FindAllStringSubmatch(path, -1)
			if subMatch != nil {
				subMatchNames := route.PathRegexp.SubexpNames()[1:]
				pathParams := make(map[string]string)
				for i, matchName := range subMatchNames {
					pathParams[matchName] = subMatch[i][1]
				}

				return route.Handler, pathParams
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
