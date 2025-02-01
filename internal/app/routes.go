package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	// Create a file server handler for static files
	fileServer := http.FileServer(http.Dir("./web/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.Session.LoadAndSave)

	router.Handler("GET", "/login", dynamic.ThenFunc(app.login))
	router.Handler("POST", "/login", dynamic.ThenFunc(app.loginPost))
	router.Handler("GET", "/register", dynamic.ThenFunc(app.register))
	router.Handler("POST", "/register", dynamic.ThenFunc(app.registerPost))

	standard := alice.New(app.logRequest)

	return standard.Then(router)
}
