package app

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"net/http"
)

func (app *Application) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Login page"))
}

func (app *Application) loginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()

	if err != nil {
		app.ErrorLog.Print(err)
	}

	loginForm := &dtos.UserLoginForm{
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
	}

	var userId uint64
	userId, err = app.AuthService.Login(loginForm)

	if err != nil {
		app.ErrorLog.Print(err)
	}

	app.Session.Put(r.Context(), "user_id", userId)

	w.Write([]byte(fmt.Sprintf("Logged in userId %d", userId)))
}
