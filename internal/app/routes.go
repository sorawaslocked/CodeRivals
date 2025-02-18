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
	router.Handler("POST", "/logout", dynamic.ThenFunc(app.logout))
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

	router.Handler("POST", "/comments/edit", dynamic.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		app.editComment(w, r, params)
	})))
	router.Handler("GET", "/problems/:url/solutions", dynamic.Then(http.HandlerFunc(app.solutions)))
	router.Handler("GET", "/solutions/:id", dynamic.ThenFunc(app.solution))
	router.Handler("POST", "/api/solutions/:id/vote", dynamic.ThenFunc(app.handleSolutionVote))
	router.Handler("POST", "/solutions/:id/comments", dynamic.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		app.createSolutionComment(w, r, params)
	})))
	router.Handler("POST", "/solutions/:id/comments/edit", dynamic.ThenFunc(app.editSolutionComment))
	router.Handler("POST", "/solutions/:id/comments/delete", dynamic.ThenFunc(app.deleteSolutionComment))
	router.Handler("GET", "/learning-materials", dynamic.ThenFunc(app.learningMaterials))

	// Admin routes
	router.Handler("GET", "/admin", dynamic.ThenFunc(app.adminDashboard))

	router.Handler("GET", "/admin/problems", dynamic.ThenFunc(app.adminProblems))
	router.Handler("GET", "/admin/problems/:id", dynamic.ThenFunc(app.adminGetProblem))
	router.Handler("POST", "/admin/problems", dynamic.ThenFunc(app.adminCreateProblem))
	router.Handler("POST", "/admin/problems/:id", dynamic.ThenFunc(app.adminUpdateProblem))
	router.Handler("DELETE", "/admin/problems/:id", dynamic.ThenFunc(app.adminDeleteProblem))
	router.Handler("GET", "/admin/problems/:id/testcases", dynamic.ThenFunc(app.adminGetTestCases))
	router.Handler("POST", "/admin/problems/:id/testcases", dynamic.ThenFunc(app.adminUpdateTestCases))

	router.Handler("GET", "/admin/users", dynamic.ThenFunc(app.adminUsers))
	//router.Handler("GET", "/admin/users/:id/details", dynamic.ThenFunc(app.adminGetUserDetails))
	router.Handler("POST", "/admin/users/:id/role", dynamic.ThenFunc(app.adminUpdateUserRole))
	router.Handler("POST", "/admin/users/:id/reset-password", dynamic.ThenFunc(app.adminResetUserPassword))

	router.Handler("GET", "/admin/topics", dynamic.ThenFunc(app.adminTopics))
	router.Handler("POST", "/admin/topics", dynamic.ThenFunc(app.adminTopicsCreate))
	router.Handler("POST", "/admin/topics/:id", dynamic.ThenFunc(app.adminTopicsUpdate))
	router.Handler("POST", "/admin/topics/:id/delete", dynamic.ThenFunc(app.adminTopicsDelete))

	standard := alice.New(app.logRequest)

	return standard.Then(router)
}
