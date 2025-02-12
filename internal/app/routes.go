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

	router.NotFound = dynamic.ThenFunc(app.notFound)
	router.PanicHandler = app.panicHandler

	router.Handler("GET", "/", dynamic.ThenFunc(app.home))
	router.Handler("GET", "/login", dynamic.ThenFunc(app.login))
	router.Handler("POST", "/login", dynamic.ThenFunc(app.loginPost))
	router.Handler("GET", "/register", dynamic.ThenFunc(app.register))
	router.Handler("POST", "/register", dynamic.ThenFunc(app.registerPost))
	router.Handler("GET", "/problems", dynamic.ThenFunc(app.problems))
	router.Handler("GET", "/problems/:url", dynamic.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		app.showProblem(w, r, params)
	})))
	router.Handler("GET", "/leaderboard", dynamic.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.showLeaderboard(w, r, nil)
	})))
	router.Handler("GET", "/profile", dynamic.ThenFunc(app.profile))
	router.Handler("POST", "/problems/:url/submit", dynamic.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		app.postSubmission(w, r, params)
	})))
	router.Handler("GET", "/submissions", dynamic.ThenFunc(app.userSubmissions))
	router.Handler("POST", "/createSolution", dynamic.ThenFunc(app.createSolution))
	router.Handler("POST", "/solution", dynamic.ThenFunc(app.postSolution))
	router.Handler("POST", "/problems/:url/comments", dynamic.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		app.createComment(w, r, params)
	})))
	router.Handler("POST", "/comments/delete", dynamic.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		app.deleteComment(w, r, params)
	})))

	standard := alice.New(app.logRequest)

	return standard.Then(router)
}
