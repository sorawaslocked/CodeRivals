package app

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"net/http"
)

func (app *Application) login(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Write([]byte("Login page"))
}

func (app *Application) loginPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := req.ParseForm()

	if err != nil {
		app.ErrorLog.Print(err)
	}

	loginForm := &dtos.UserLoginForm{
		Username: req.PostForm.Get("username"),
		Password: req.PostForm.Get("password"),
	}

	var userId uint64
	userId, err = app.AuthService.Login(loginForm)

	if err != nil {
		app.ErrorLog.Print(err)
	}

	w.Write([]byte(fmt.Sprintf("Logged in userId %d", userId)))
}
