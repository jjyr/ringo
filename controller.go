package ringo

import (
	"fmt"
	"path"
	"reflect"
	"strings"
)

type Controllerable interface {
	ControllerName() string
}

type Controller struct {
	Name string
}

func isAlphabetUpper(s string) bool {
	return strings.ToUpper(s) == s
}

func getControllerName(c Controllerable) string {
	name := c.ControllerName()

	if name == "" {
		controllerType := fmt.Sprintf("%T", c)
		// *pkg.TestRingoController
		tmpName := controllerType
		dotIndex := strings.LastIndex(tmpName, ".")
		if dotIndex > -1 && len(tmpName) > dotIndex {
			// TestRingoController
			tmpName = tmpName[dotIndex+1:]
			if idx := strings.Index(tmpName, "Controller"); idx > 0 {
				// TestRingo
				tmpName = tmpName[0:idx]

				prev := ""
				for _, w := range tmpName {
					s := string(w)
					if isAlphabetUpper(s) {
						if !isAlphabetUpper(prev) {
							s = "-" + s
						}
					}
					// test-ringo
					name += strings.ToLower(s)
				}
			}
		}
	}

	return name
}

func (c *Controller) ControllerName() string {
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
	controllerName := getControllerName(c)
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
