package ringo

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/jjyr/ringo/route"
	"github.com/jjyr/ringo/common"
	"github.com/jjyr/ringo/context"
)

type App struct {
	middlewares       []common.MiddlewareFunc
	handleHTTPRequest common.HandlerFunc
	*route.RouteManage
	*TemplateManage
}

func NewApp() *App {
	app := App{}
	app.RouteManage = route.NewRouteManage()
	app.TemplateManage = newTemplateManage()
	return &app
}

func (app *App) Use(middlewreFunc common.MiddlewareFunc) {
	app.middlewares = append(app.middlewares, middlewreFunc)
}

func (app *App) Run(addr string) error {
	log.Infof("Listening on %s, start serve HTTP", addr)
	app.initForServe()
	err := http.ListenAndServe(addr, app)
	if err != nil {
		log.Errorln(err)
	}
	return err
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)
	c := context.NewContext()
	c.SetRequest(r)
	c.ResponseWriter = context.NewResponseWriter(w)
	c.TemplateManage = app.TemplateManage
	app.handleHTTPRequest(c)
}

func (app *App) initForServe() {
	// use middlewares
	app.handleHTTPRequest = app.applyMiddlewares(app.defaultHandleHTTPRequest)
}

func (app *App) applyMiddlewares(handler common.HandlerFunc) common.HandlerFunc {
	for _, middleware := range app.middlewares {
		handler = middleware(handler)
	}

	return handler
}

func (app *App) defaultHandleHTTPRequest(c common.Context) {
	requestPath := c.Request().URL.Path
	handler, params, redirect := app.RouteManage.MatchRoute(requestPath, c.Request().Method)
	c.SetParams(params)

	if redirect {
		// handle redirect without trailing slash
		handler = func(c common.Context) {
			c.Redirect(301, requestPath[:len(requestPath)-1])
		}
	}

	if handler == nil {
		handler = handleNotFound
		log.Printf("Not found route for %s", c.Request().RequestURI)
	}

	handler(c)

	if !c.HasRendered() {
		log.Printf("Empty render")
	}
}

func handleNotFound(c common.Context) {
	statusCode := http.StatusNotFound
	text := http.StatusText(statusCode)
	c.String(statusCode, text)
}
