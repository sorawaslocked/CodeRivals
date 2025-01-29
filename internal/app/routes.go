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

	standard := alice.New(app.logRequest)

	return standard.Then(router)
}
