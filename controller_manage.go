package ringo

import (
	"fmt"
	"path"
	"reflect"
	"strings"
)

type ControllerManage struct {
	controllers              map[string]controllerRegistry
	router                   *Router
	controllerRoutePathCache map[string]string
}

func newControllerManage(router *Router) *ControllerManage {
	controllerManage := &ControllerManage{router: router, controllers: make(map[string]controllerRegistry)}
	controllerManage.controllerRoutePathCache = make(map[string]string)
	return controllerManage
}

type ControllerOption struct {
	// controller path prefix
	Prefix string
	// override controller name
	Name string
}

type controllerRegistry struct {
	controller       Controllerable
	controllerOption *ControllerOption
}

// AddController add newcontroller
func (m *ControllerManage) AddController(c Controllerable, controllerOption *ControllerOption) {
	m.registerController(c, controllerOption)
	// register to router
	controllerName := getControllerName(c, controllerOption)
	if controllerName == "" {
		panic(fmt.Errorf("Controller Name is empty, %+v", c))
	}
	controller := reflect.ValueOf(c)
	for _, routeOption := range append(controllerDefaultRouteOptions, c.GetRoutes()...) {
		handlerValue := controller.MethodByName(routeOption.Handler)
		if handlerValue.IsValid() {
			for _, method := range routeOption.Methods {
				routePath := controllerRoutePathFromOption(c, controllerOption, routeOption)
				m.router.AddRoute(routePath, method, func(context *Context) {
					handlerValue.Call([]reflect.Value{reflect.ValueOf(context)})
				})
			}
		}
	}
}

// generate controller route path from options
func controllerRoutePathFromOption(c Controllerable, controllerOption *ControllerOption, routeOption ControllerRouteOption) string {
	var controllerPrefix string
	if controllerOption != nil {
		controllerPrefix = controllerOption.Prefix
	}
	controllerName := getControllerName(c, controllerOption)
	controllerName = routeOption.Prefix + controllerName + routeOption.Suffix
	routePath := path.Join("/", controllerName)
	if routeOption.Member {
		routePath = path.Join(routePath, "/:id")
	}
	routePath = path.Join(routePath, routeOption.Path)
	routePath = controllerPrefix + routePath
	return routePath
}

func (m *ControllerManage) registerController(c Controllerable, controllerOption *ControllerOption) {
	name := getControllerName(c, controllerOption)
	if _, exists := m.controllers[name]; exists {
		panic(fmt.Errorf("Controller name '%s' duplicate", name))
	}
	m.controllers[name] = controllerRegistry{controller: c, controllerOption: controllerOption}
}

func (m *ControllerManage) findRoutePath(name string, handler string) string {
	controllerEntry := m.controllers[name]
	for _, routeOption := range append(controllerDefaultRouteOptions, controllerEntry.controller.GetRoutes()...) {
		if routeOption.Handler == handler {
			return controllerRoutePathFromOption(controllerEntry.controller, controllerEntry.controllerOption, routeOption)
		}
	}

	return ""
}

// FindControllerRoutePath Find controller action path, return first path if multi route registered, panic if not found
// paramvalue are pairs of param/value like "id", "2"
func (m *ControllerManage) FindControllerRoutePath(name string, handler string, paramvalue ...string) string {
	key := fmt.Sprintf("%s#%s", name, handler)
	routePath, ok := m.controllerRoutePathCache[key]
	if !ok {
		routePath = m.findRoutePath(name, handler)
		if routePath == "" {
			panic(fmt.Errorf("Cannot find handler '%s' in controller '%s'", handler, name))
		}
		m.controllerRoutePathCache[key] = routePath
	}
	if len(paramvalue) > 0 {
		for i := 0; i < len(paramvalue); i += 2 {
			paramvalue[i] = ":" + paramvalue[i]
		}
		routePath = strings.NewReplacer(paramvalue...).Replace(routePath)
	}
	return routePath
}
