package app

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"github.com/sorawaslocked/CodeRivals/internal/validator"
	"net/http"
	"strconv"
	"strings"
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	// Get selected topics from query
	var selectedTopicIDs []int
	if topicsStr := r.URL.Query().Get("topics"); topicsStr != "" {
		topicStrs := strings.Split(topicsStr, ",")
		for _, ts := range topicStrs {
			if id, err := strconv.Atoi(ts); err == nil {
				selectedTopicIDs = append(selectedTopicIDs, id)
			}
		}
	}

	itemsPerPage := 10
	offset := (page - 1) * itemsPerPage

	// Get problems with topic filtering
	problems, totalItems, err := app.ProblemService.GetPaginatedProblemsWithTopics(offset, itemsPerPage, selectedTopicIDs)
	if err != nil {
		app.ErrorLog.Print(err)
		return
	}

	// Get all topics
	topics, err := app.TopicService.GetAllTopics()
	if err != nil {
		app.ErrorLog.Print(err)
		return
	}

	data := app.newTemplateData(r)
	data.Problems = problems
	data.Topics = topics
	data.Pagination = NewPagination(page, totalItems, itemsPerPage, r.URL.Query())

	app.render(w, r, "problem/problems.gohtml", data)
}

func (app *Application) showProblem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := ps.ByName("url")

	problem, err := app.ProblemService.GetProblemByURL(url)
	if err != nil {
		app.ErrorLog.Print(err)
		http.NotFound(w, r)
		return
	}

	examples, err := app.ProblemService.GetProblemExamples(problem.ID)
	if err != nil {
		app.ErrorLog.Print(err)
	}

	data := app.newTemplateData(r)
	data.Form = problem
	data.Examples = examples

	app.render(w, r, "problem/problem.gohtml", data)
}

func (app *Application) profile(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	if td.AuthenticatedUserId == 0 {
		app.userError(w, r, "You are not authenticated")
	}
}

func (app *Application) showLeaderboard(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users, err := app.LeaderBoardService.GetLeaderboard()
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := app.newTemplateData(r)
	data.Users = users

	app.render(w, r, "leaderboard/leaderboard.gohtml", data)
}

func (app *Application) submitSolution(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	userId := app.Session.GetInt(r.Context(), "authenticatedUserId")
	if userId == 0 {
		http.Error(w, "You are not authenticated", http.StatusUnauthorized)
		return
	}

	problemUrl := ps.ByName("url")
	problem, err := app.ProblemService.GetProblemByURL(problemUrl)
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Problem not found", http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	code := r.PostFormValue("code")
	if code == "" {
		http.Error(w, "Code can't be empty", http.StatusBadRequest)
		return
	}

	err = app.SubmissionService.Submit(userId, problem.ID, code)
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Failed to submit solution", http.StatusInternalServerError)
		return
	}

	app.Session.Put(r.Context(), "flash", "Solution submitted successfully!")
	http.Redirect(w, r, fmt.Sprintf("/problems/%s", problemUrl), http.StatusSeeOther)
}
