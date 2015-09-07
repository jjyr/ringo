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
}
