package ringo

import (
	"fmt"
	"path"
	"reflect"
)

type Controllerable interface {
	GetName() string
}

type Controller struct {
	Name string
}

func (c *Controller) GetName() string {
	return c.Name
}

type controllerHandlerRouteDescription struct {
	handlerName string
	method      string
	path        string
}

var controllerHandlerRouteDescriptions []controllerHandlerRouteDescription

func init() {
	controllerHandlerRouteDescriptions = []controllerHandlerRouteDescription{
		{handlerName: "List", method: "GET", path: ""},
		{handlerName: "Create", method: "POST", path: ""},
		{handlerName: "Get", method: "GET", path: "/:id"},
		{handlerName: "Delete", method: "DELETE", path: "/:id"},
		{handlerName: "Update", method: "PUT", path: "/:id"},
	}
}

func registerToRouter(r *Router, c Controllerable) {
	controllerName := c.GetName()
	if controllerName == "" {
		panic(fmt.Errorf("Controller Name is empty, %+v", c))
	}
	controller := reflect.ValueOf(c)
	for _, handlerRoute := range controllerHandlerRouteDescriptions {
		handlerValue := controller.MethodByName(handlerRoute.handlerName)
		if handlerValue.IsValid() {
			r.AddRoute(path.Join("/", controllerName, handlerRoute.path), handlerRoute.method, func(context *Context) {
				handlerValue.Call([]reflect.Value{reflect.ValueOf(context)})
			})
		}
	}
}
