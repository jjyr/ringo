package common

type HandlerFunc func(c Context)

type MiddlewareFunc func(HandlerFunc) HandlerFunc
