package ringo

import (
	"fmt"
	"path"
	"reflect"
)

type ControllerManage struct {
	controllers map[string]controllerRegistry
	router      *Router
}

func newControllerManage(router *Router) *ControllerManage {
	return &ControllerManage{router: router, controllers: make(map[string]controllerRegistry)}
}

type ControllerOption struct {
	Prefix string
	Name   string
}

type controllerRegistry struct {
	controller       Controllerable
	controllerOption *ControllerOption
}

func (m *ControllerManage) AddController(c Controllerable, controllerOption *ControllerOption) {
	m.registerController(c, controllerOption)
	// register to router
	controllerName := GetControllerName(c)
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
	var controllerName, controllerPrefix string
	controllerName = routeOption.Name
	if controllerOption != nil {
		controllerPrefix = controllerOption.Prefix
		if controllerName == "" {
			controllerName = controllerOption.Name
		}
	}

	if controllerName == "" {
		controllerName = GetControllerName(c)
	}
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
	name := GetControllerName(c)
	if _, exists := m.controllers[name]; exists {
		panic(fmt.Errorf("Controller name '%s' duplicate", name))
	}
	m.controllers[name] = controllerRegistry{controller: c, controllerOption: controllerOption}
}

// FindControllerRoutePath Find controller action path, return first path if multi route registered, panic if not found
func (m *ControllerManage) FindControllerRoutePath(name string, handler string) string {
	controllerEntry := m.controllers[name]
	for _, routeOption := range append(controllerDefaultRouteOptions, controllerEntry.controller.GetRoutes()...) {
		if routeOption.Handler == handler {
			return controllerRoutePathFromOption(controllerEntry.controller, controllerEntry.controllerOption, routeOption)
		}
	}

	panic(fmt.Errorf("Cannot find handler '%s' in controller '%s'", handler, name))
}