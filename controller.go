package ringo

import (
	"fmt"
	"strings"
)

type Controllerable interface {
	ControllerName() string
	SetControllerName(string)
	AddRoutes(...ControllerRouteOption)
	GetRoutes() []ControllerRouteOption
}

type Controller struct {
	Name   string
	routes []ControllerRouteOption
}

func (c *Controller) ControllerName() string {
	return c.Name
}

func (c *Controller) SetControllerName(name string) {
	if c.ControllerName() != "" {
		panic(fmt.Errorf("Should not override non empty controller name"))
	}
	c.Name = name
}

// AddRoutes add customize route to controller
func (c *Controller) AddRoutes(routeOptions ...ControllerRouteOption) {
	for _, routeOption := range routeOptions {
		// check option
		routeOption.confirm()
		c.routes = append(c.routes, routeOption)
	}
}

func (c *Controller) GetRoutes() []ControllerRouteOption {
	return c.routes
}

func isAlphabetUpper(s string) bool {
	return strings.ToUpper(s) == s
}

// GetControllerName return controller name, generate by type name if not manually set
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

// ControllerRouteOption options to customize controller route
type ControllerRouteOption struct {
	// handler method name
	Handler string
	// http method
	Method string
	// http methods
	Methods []string
	// route path
	Path string
	// route prefix
	Prefix string
	// override controller name
	Name string
	// route suffix
	Suffix string
	// as member route, like: "/users/1/xxx"
	Member bool
	// as collection route, like: "/users/xxx"
	Collection bool
}

// validate value
func (routerOption *ControllerRouteOption) confirm() {
	if routerOption.Member == routerOption.Collection {
		panic(fmt.Errorf("Router option must be member or collection"))
	}

	if routerOption.Method != "" {
		routerOption.Methods = append(routerOption.Methods, routerOption.Method)
	}

	if len(routerOption.Methods) == 0 {
		panic(fmt.Errorf("Router option must provide at least one method"))
	}
}

var controllerDefaultRouteOptions []ControllerRouteOption

func init() {
	controllerDefaultRouteOptions = []ControllerRouteOption{
		{Handler: "List", Method: "GET", Collection: true},
		{Handler: "Create", Method: "POST", Collection: true},
		{Handler: "Get", Method: "GET", Member: true},
		{Handler: "Delete", Method: "DELETE", Member: true},
		{Handler: "Update", Methods: []string{"PUT", "PATCH"}, Member: true},
		{Handler: "New", Method: "GET", Collection: true, Prefix: "new-"},
		{Handler: "Edit", Method: "GET", Member: true, Path: "/edit"},
	}

	for i, routeOption := range controllerDefaultRouteOptions {
		routeOption.confirm()
		controllerDefaultRouteOptions[i] = routeOption
	}
}
