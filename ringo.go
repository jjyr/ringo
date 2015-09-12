package ringo

type HandlerFunc func(c *Context)

type MiddlewareFunc func(HandlerFunc) HandlerFunc
