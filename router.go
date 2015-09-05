package ringo

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
type DefaultRouter struct{}

// NewRouter make a default router
func NewRouter() Router {
	return &DefaltRouter{}
}
