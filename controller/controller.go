package controller

import (
	"github.com/jjyr/ringo/route"
)

type Controller struct{}

func (c *Controller) GetRoutes() []route.RouteOption {
	return controllerDefaultRouteOptions
}

var controllerDefaultRouteOptions []route.RouteOption

func init() {
	controllerDefaultRouteOptions = []route.RouteOption{
		{Handler: "Get", Methods: []string{"GET"}},
		{Handler: "Post", Methods: []string{"POST"}},
		{Handler: "Put", Methods: []string{"PUT"}},
		{Handler: "Delete", Methods: []string{"DELETE"}},
		{Handler: "Patch", Methods: []string{"PATCH"}},
		{Handler: "Option", Methods: []string{"OPTION"}},
		{Handler: "Head", Methods: []string{"HEAD"}},
	}
}
