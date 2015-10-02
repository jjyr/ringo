package ringo

import (
	"fmt"
	"path"
	"reflect"
	"strings"
)

type Controllerable interface {
	ControllerName() string
	SetControllerName(string)
}

type Controller struct {
	Name string
}

func (c *Controller) ControllerName() string {
	return c.Name
}

func (c *Controller) SetControllerName(name string) {
	c.Name = name
}

func isAlphabetUpper(s string) bool {
	return strings.ToUpper(s) == s
}

func GetControllerName(c Controllerable) string {
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
					prev = s
				}
			}
		}

		c.SetControllerName(name)
	}

	return name
}

type ControllerRouterOption struct {
	Handler    string
	Method     []string
	Path       string
	Name       string
	NamePrefix string
	NameSuffix string
	Member     bool
	Collection bool
}

var controllerDefaultRouterOptions []ControllerRouterOption

func init() {
	controllerDefaultRouterOptions = []ControllerRouterOption{
		{Handler: "List", Method: []string{"GET"}, Collection: true},
		{Handler: "Create", Method: []string{"POST"}, Collection: true},
		{Handler: "Get", Method: []string{"GET"}, Member: true},
		{Handler: "Delete", Method: []string{"DELETE"}, Member: true},
		{Handler: "Update", Method: []string{"PUT", "PATCH"}, Member: true},
		{Handler: "New", Method: []string{"GET"}, Collection: true, NamePrefix: "new-"},
		{Handler: "Edit", Method: []string{"GET"}, Member: true, Path: "/edit"},
	}
}

func pathFromRouterOption(c Controllerable, routerOption ControllerRouterOption) string {
	if routerOption.Member == routerOption.Collection {
		panic(fmt.Errorf("Router option must be member or collection"))
	}
	controllerName := routerOption.Name
	if controllerName == "" {
		controllerName = GetControllerName(c)
	}
	controllerName = routerOption.NamePrefix + controllerName + routerOption.NameSuffix
	routerPath := path.Join("/", controllerName)
	if routerOption.Member {
		routerPath = path.Join(routerPath, "/:id")
	}
	routerPath = path.Join(routerPath, routerOption.Path)

	return routerPath
}

func registerToRouter(r *Router, c Controllerable, otherRouters ...ControllerRouterOption) {
	controllerName := GetControllerName(c)
	if controllerName == "" {
		panic(fmt.Errorf("Controller Name is empty, %+v", c))
	}
	controller := reflect.ValueOf(c)
	for _, routerOption := range append(controllerDefaultRouterOptions, otherRouters...) {
		handlerValue := controller.MethodByName(routerOption.Handler)
		if handlerValue.IsValid() {
			for _, m := range routerOption.Method {
				r.AddRoute(pathFromRouterOption(c, routerOption), m, func(context *Context) {
					handlerValue.Call([]reflect.Value{reflect.ValueOf(context)})
				})
			}
		}
	}
}
