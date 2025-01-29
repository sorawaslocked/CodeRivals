package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"net/http"
)

func (app *Application) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := &templateData{
		Form: &dtos.UserLoginForm{},
	}

	app.render(w, r, "auth/login.html", data)
}

func (app *Application) loginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	form := &dtos.UserLoginForm{
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
	}

	userId, err := app.AuthService.Login(form)
	if err != nil || !form.Valid() {
		data := &templateData{
			Form: form,
		}
		app.render(w, r, "auth/login.html", data)
		return
	}

	app.Session.Put(r.Context(), "user_id", userId)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := &templateData{
		Form: &dtos.UserRegisterForm{},
	}

	app.render(w, r, "auth/register.html", data)
}

func (app *Application) registerPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	form := &dtos.UserRegisterForm{
		Username:        r.PostForm.Get("username"),
		Email:           r.PostForm.Get("email"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirmPassword"),
	}

	err = app.AuthService.Register(form)
	if err != nil || !form.Valid() {
		data := &templateData{
			Form: form,
		}
		app.render(w, r, "auth/register.html", data)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, data *templateData) {
	// Placeholder for template rendering
	// A simple response for now
	w.Write([]byte("Template: " + name))
}
