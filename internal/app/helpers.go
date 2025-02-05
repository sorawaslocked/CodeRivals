package app

import (
	"net/http"
	"runtime/debug"
)

func (app *Application) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	app.render(w, r, "errors/not_found.gohtml", nil)
}

func (app *Application) serverError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	app.render(w, r, "errors/internal_server_error.gohtml", nil)
}

func (app *Application) panicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	app.ErrorLog.Printf("panic: %v\n%s", err, debug.Stack())

	app.Session.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "errors/internal_server_error.gohtml", nil)
	})).ServeHTTP(w, r)
}
