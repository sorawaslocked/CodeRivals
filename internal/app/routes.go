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

	router.GET("/login", app.login)
	router.POST("/login", app.loginPost)
	router.GET("/register", app.register)
	router.POST("/register", app.registerPost)

	standard := alice.New(app.logRequest)

	return standard.Then(router)
}
