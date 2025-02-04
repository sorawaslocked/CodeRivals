package app

import (
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"github.com/sorawaslocked/CodeRivals/internal/validator"
	"net/http"
)

func (app *Application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = dtos.UserLoginForm{}

	app.render(w, r, "auth/login.gohtml", data)
}

func (app *Application) loginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	form := &dtos.UserLoginForm{
		Username: username,
		Password: password,
	}

	userId, err := app.AuthService.Login(form)

	if err != nil || !form.Valid() {
		data := app.newTemplateData(r)

		app.ErrorLog.Println(err)
		app.render(w, r, "auth/login.gohtml", data)
		return
	}

	err = app.Session.RenewToken(r.Context())

	app.Session.Put(r.Context(), "authenticatedUserId", userId)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Application) register(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = dtos.UserRegisterForm{}

	app.render(w, r, "auth/register.gohtml", data)
}

func (app *Application) registerPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("username")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	confirmPassword := r.PostForm.Get("confirmPassword")

	form := &dtos.UserRegisterForm{
		Username:        username,
		Email:           email,
		Password:        password,
		ConfirmPassword: confirmPassword,
		Validator:       validator.Validator{},
	}

	err = app.AuthService.Register(form)

	if err != nil || !form.Valid() {
		data := app.newTemplateData(r)

		app.render(w, r, "auth/register.gohtml", data)
		app.ErrorLog.Print(err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Application) problems(w http.ResponseWriter, r *http.Request) {
	problems, err := app.ProblemService.GetAllProblems()

	if err != nil {
		app.ErrorLog.Print(err)
	}

	data := app.newTemplateData(r)
	data.Problems = problems

	app.render(w, r, "problems.gohtml", data)
}
