package ringo

import (
	"net/http"

	"github.com/golang/glog"
)

type App struct {
	router Router
}

func (app *App) GetRouter() Router { return app.router }

func NewApp() *App {
	return &App{router: *NewRouter()}
}

func (app *App) Run(addr string) error {
	err := http.ListenAndServe(addr, app)
	if err != nil {
		glog.Errorf("%s", err)
	}
	return err
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infof("%+v", r)
	c := NewContext()
	c.Request = r
	c.ResponseWriter = newResponseWriter(w)
	app.handleHTTPRequest(c)
}

func (app *App) handleHTTPRequest(c *Context) {
	handler, pathParams := app.router.MatchRoute(c.Request.RequestURI, c.Request.Method)
	c.PathParams = pathParams
	handler(c)

	w := c.ResponseWriter.(*ResponseWriter)
	if !w.Written() {
		w.WriteHeader(500)
	}
}
