package ringo

import (
	"log"
	"net/http"
)

type App struct {
	router *Router
}

func (app *App) GetRouter() *Router { return app.router }

func NewApp() *App {
	return &App{router: NewRouter()}
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
	app.handleHTTPRequest(c)
}

func (app *App) handleHTTPRequest(c *Context) {
	handler, pathParams := app.router.MatchRoute(c.Request.RequestURI, c.Request.Method)
	c.PathParams = pathParams
	if handler != nil {
		handler(c)
	} else {
		// handler 404 not found
		log.Printf("Not found route for %s", c.RequestURI)
	}

	w := c.ResponseWriter.(*ResponseWriter)
	if !w.Written() {
		w.WriteHeader(500)
	}
}
