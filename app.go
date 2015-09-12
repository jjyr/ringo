package ringo

import (
	"log"
	"net/http"
)

type App struct {
	*Router
	middlewares []MiddlewareFunc
}

func NewApp() *App {
	return &App{Router: NewRouter()}
}

func (app *App) Use(middlewreFunc MiddlewareFunc) {
	app.middlewares = append(app.middlewares, middlewreFunc)
}

func (app *App) Run(addr string) error {
	log.Printf("Listening on %s, start serve HTTP", addr)
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
	app.applyMiddlewares(app.handleHTTPRequest)(c)
}

func (app *App) applyMiddlewares(handler HandlerFunc) HandlerFunc {
	for _, middleware := range app.middlewares {
		handler = middleware(handler)
	}

	return handler
}

func (app *App) handleHTTPRequest(c *Context) {
	handler, params := app.Router.MatchRoute(c.Request.URL.Path, c.Request.Method)
	c.Params = params

	if handler == nil {
		handler = handleNotFound
		log.Printf("Not found route for %s", c.RequestURI)
	}

	if !c.Rendered() {
		panic("Empty response, missing render call in handler")
	}
}

func handleNotFound(c *Context) {
	statusCode := http.StatusNotFound
	text := http.StatusText(statusCode)
	c.Render(statusCode, text)
}
