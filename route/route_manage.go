package route

import (
	"reflect"
	"github.com/jjyr/ringo/common"
)

// Handle routing
type RouteManage struct {
	*Router
	controllerRoutePathCache map[string]string
}

func NewRouteManage() *RouteManage {
	router := NewRouter()
	routeManage := &RouteManage{Router: router}
	routeManage.controllerRoutePathCache = make(map[string]string)
	return routeManage
}

// Add Routeable
func (m *RouteManage) Add(routePath string, r Routeable) {
	resource := reflect.ValueOf(r)
	for _, routeOption := range r.GetRoutes() {
		handlerValue := resource.MethodByName(routeOption.Handler)
		if handlerValue.IsValid() {
			for _, method := range routeOption.Methods {
				m.Router.AddRoute(routePath, method, func(context common.Context) {
					handlerValue.Call([]reflect.Value{reflect.ValueOf(context)})
				})
			}
		}
	}
}
