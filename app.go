package ringo

import (
	"errors"
	"log"
	"net/http"
)

type App struct {
	*Router
	middlewares       []MiddlewareFunc
	handleHTTPRequest HandlerFunc
	*ControllerManage
}

func NewApp() *App {
	defaultRouter := NewRouter()
	return &App{Router: defaultRouter, ControllerManage: newControllerManage(defaultRouter)}
}

func (app *App) Use(middlewreFunc MiddlewareFunc) {
	app.middlewares = append(app.middlewares, middlewreFunc)
}

func (app *App) Run(addr string) error {
	log.Printf("Listening on %s, start serve HTTP", addr)
	app.initForServe()
	err := http.ListenAndServe(addr, app)
	if err != nil {
		log.Printf("%s", err)
	}
	return err
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)
	c := NewContext()
	c.Request = r
	c.ResponseWriter = newResponseWriter(w)
	app.handleHTTPRequest(c)
}

func (app *App) initForServe() {
	// use middlewares
	app.handleHTTPRequest = app.applyMiddlewares(app.defaultHandleHTTPRequest)
}

func (app *App) applyMiddlewares(handler HandlerFunc) HandlerFunc {
	for _, middleware := range app.middlewares {
		handler = middleware(handler)
	}

	return handler
}

func (app *App) defaultHandleHTTPRequest(c *Context) {
	requestPath := c.Request.URL.Path
	handler, params, redirect := app.Router.MatchRoute(requestPath, c.Request.Method)
	c.Params = params

	if redirect {
		// handle redirect without trailing slash
		handler = func(c *Context) {
			c.Redirect(301, requestPath[:len(requestPath)-1])
		}
	}

	if handler == nil {
		handler = handleNotFound
		log.Printf("Not found route for %s", c.Request.RequestURI)
	}

	handler(c)

	if !c.Rendered() {
		panic(errors.New("Empty response, missing render call in handler"))
	}
}

func handleNotFound(c *Context) {
	statusCode := http.StatusNotFound
	text := http.StatusText(statusCode)
	c.String(statusCode, text)
}
