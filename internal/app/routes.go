package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	router.GET("/login", app.login)
	router.POST("/login", app.loginPost)
	router.GET("/register", app.register)
	router.POST("/register", app.registerPost)

	standard := alice.New(app.logRequest)

	return standard.Then(router)
}
