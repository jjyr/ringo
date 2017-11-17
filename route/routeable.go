package route

type Routeable interface {
	GetRoutes() []RouteOption
}

// RouteOption options to customize route
type RouteOption struct {
	// handler method name
	Handler string
	// http methods
	Methods []string
}
