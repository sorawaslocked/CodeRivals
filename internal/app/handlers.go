package app

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"github.com/sorawaslocked/CodeRivals/internal/validator"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *Application) logout(w http.ResponseWriter, r *http.Request) {
	err := app.Session.RenewToken(r.Context())

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	app.Session.Remove(r.Context(), "authenticatedUserId")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Get top 5 users for the leaderboard section
	topUsers, err := app.LeaderBoardService.GetLeaderboard()
	if err != nil {
		app.ErrorLog.Print(err)
	}
	data.Users = topUsers

	// Get first 5 problems for the featured problems section
	problems, _, err := app.ProblemService.GetPaginatedProblems(0, 5)
	if err != nil {
		app.ErrorLog.Print(err)
	}
	data.Problems = problems

	app.render(w, r, "home/home.gohtml", data)
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

	signature := app.ProblemService.GenerateSignature(problem)

	examples, err := app.ProblemService.GetProblemExamples(problem.ID)
	if err != nil {
		app.ErrorLog.Print(err)
	}

	comments, err := app.CommentService.GetProblemComments(problem.ID)
	if err != nil {
		app.ErrorLog.Print(err)
	}

	topLevelComments := make([]entities.Comment, 0)
	replyMap := make(map[int][]entities.Comment)

	for _, comment := range comments {
		if comment.CommentID == nil {
			topLevelComments = append(topLevelComments, comment)
		} else {
			replyMap[*comment.CommentID] = append(replyMap[*comment.CommentID], comment)
		}
	}

	data := app.newTemplateData(r)
	data.Form = problem
	data.Examples = examples
	data.Comments = topLevelComments
	data.ReplyMap = replyMap
	data.Signature = signature

	app.render(w, r, "problem/problem.gohtml", data)
}

func (app *Application) profile(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	if td.AuthenticatedUserId == 0 {
		app.userError(w, r, "You are not authenticated")
		return
	}

	// Get user information using AuthService
	user, err := app.AuthService.GetUser(td.AuthenticatedUserId)
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}

	// Update template to use Form for passing data, following existing pattern
	td.Form = struct {
		Username  string
		Points    int
		CreatedAt time.Time
	}{
		Username:  user.Username,
		Points:    user.Points,
		CreatedAt: user.CreatedAt,
	}

	app.render(w, r, "user/profile.gohtml", td)
}

func (app *Application) userSubmissions(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from session
	td := app.newTemplateData(r)

	uid := td.AuthenticatedUserId

	if uid == 0 {
		app.userError(w, r, "You are not authenticated")
	}

	// Get submissions for the user
	submissions, err := app.SubmissionService.GetAllUserSubmissions(uid)
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}

	// Prepare template data
	data := app.newTemplateData(r)
	data.Submissions = submissions

	app.render(w, r, "user/user_submissions.gohtml", data)
}

func (app *Application) testExecution(w http.ResponseWriter, r *http.Request) {}

func (app *Application) showLeaderboard(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userId := app.Session.GetInt(r.Context(), "authenticatedUserId")
	if userId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

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

func (app *Application) postSubmission(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	userId := app.Session.GetInt(r.Context(), "authenticatedUserId")
	if userId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

	_, err = app.SubmissionService.Submit(userId, problem.ID, code)
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Failed to submit solution", http.StatusInternalServerError)
		return
	}

	app.Session.Put(r.Context(), "flash", "Solution submitted successfully!")
	http.Redirect(w, r, "/submissions", http.StatusSeeOther)
}

func (app *Application) createSolution(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	formId := r.PostFormValue("submissionId")

	id, err := strconv.Atoi(formId)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	var submission *entities.ProblemSubmission
	submission, err = app.SubmissionService.GetSubmission(id)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	var problem *entities.Problem
	problem, err = app.ProblemService.GetProblem(submission.ProblemID)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	td.SubmissionForSolution = &entities.FullProblemSubmission{
		Submission: submission,
		Problem:    problem,
	}

	app.render(w, r, "user/solution_form.gohtml", td)
}

func (app *Application) postSolution(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	formId := r.PostFormValue("submissionId")
	title := r.PostFormValue("title")
	description := r.PostFormValue("description")

	id, err := strconv.Atoi(formId)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	var submission *entities.ProblemSubmission
	submission, err = app.SubmissionService.GetSubmission(id)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	var problem *entities.Problem
	problem, err = app.ProblemService.GetProblem(submission.ProblemID)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	// create solution

	solution := &entities.ProblemSolution{
		ProblemId:   problem.ID,
		UserId:      td.AuthenticatedUserId,
		Title:       title,
		Description: description,
		Code:        submission.Code,
	}

	err = app.ProblemService.CreateProblemSolution(solution)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	http.Redirect(w, r, fmt.Sprintf("/problems/%s", problem.Url), http.StatusSeeOther)
}

func (app *Application) createComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := app.Session.GetInt(r.Context(), "authenticatedUserId")
	if userId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

	text := r.PostFormValue("text")
	if text == "" {
		http.Error(w, "Comment text is required", http.StatusBadRequest)
		return
	}

	// Check if this is a reply to another comment
	parentCommentID := r.PostFormValue("parent_comment_id")
	var parentID *int
	if parentCommentID != "" {
		// This is a reply
		id, err := strconv.Atoi(parentCommentID)
		if err != nil {
			app.ErrorLog.Print(err)
			http.Error(w, "Invalid parent comment ID", http.StatusBadRequest)
			return
		}
		parentID = &id
	}

	err = app.CommentService.CreateReply(userId, problem.ID, text, parentID)

	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/problems/%s", problemUrl), http.StatusSeeOther)
}

func (app *Application) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := app.Session.GetInt(r.Context(), "authenticatedUserId")
	if userId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentId, err := strconv.Atoi(r.PostFormValue("comment_id"))
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	err = app.CommentService.DeleteComment(commentId)
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	// Redirect back to the previous page
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *Application) editComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := app.Session.GetInt(r.Context(), "authenticatedUserId")
	if userId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	commentId, err := strconv.Atoi(r.PostFormValue("comment_id"))
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	text := r.PostFormValue("text")
	if text == "" {
		http.Error(w, "Comment text is required", http.StatusBadRequest)
		return
	}

	err = app.CommentService.UpdateComment(commentId, text)
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	// Redirect back to the previous page
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *Application) solutions(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	url := params.ByName("url")

	problem, err := app.ProblemService.GetProblemByURL(url)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	var solutions []*entities.ProblemSolution
	solutions, err = app.ProblemService.GetSolutionsForProblem(problem.ID)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	var fullSolutions []*entities.ProblemSolutionDisplay

	for _, solution := range solutions {
		user, err := app.AuthService.GetUser(solution.UserId)

		if err != nil {
			app.ErrorLog.Print(err)
			app.serverError(w, r)
		}

		fullSolution := &entities.ProblemSolutionDisplay{
			Solution:    solution,
			SubmittedBy: user.Username,
		}

		fullSolutions = append(fullSolutions, fullSolution)
	}

	td := app.newTemplateData(r)
	td.Solutions = fullSolutions
	td.ProblemTitle = problem.Title

	app.render(w, r, "problem/solutions.gohtml", td)
}

func (app *Application) solution(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	solutionId, err := strconv.Atoi(id)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	var solution *entities.ProblemSolution
	solution, err = app.ProblemService.GetSolutionById(solutionId)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	var problem *entities.Problem
	problem, err = app.ProblemService.GetProblem(solution.ProblemId)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	var user *entities.User
	user, err = app.AuthService.GetUser(solution.UserId)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	comments, err := app.ProblemSolutionCommentService.GetSolutionComments(solutionId)
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	topLevelSolutionComments := make([]entities.ProblemSolutionComment, 0)
	solutionReplyMap := make(map[int][]entities.ProblemSolutionComment)

	for _, comment := range comments {
		if comment.CommentID == nil {
			topLevelSolutionComments = append(topLevelSolutionComments, comment)
		} else {
			solutionReplyMap[*comment.CommentID] = append(solutionReplyMap[*comment.CommentID], comment)
		}
	}

	td := app.newTemplateData(r)
	td.Solution = solution
	td.ProblemTitle = problem.Title
	td.ProblemUrl = problem.Url
	td.SolutionSubmittedBy = user.Username
	td.SolutionComments = topLevelSolutionComments
	td.SolutionReplyMap = solutionReplyMap

	voteStatus := &entities.ProblemSolutionVoteStatus{}

	var upvote bool
	upvote, err = app.ProblemService.GetUpvoteBySolutionIdAndUserId(solutionId, td.AuthenticatedUserId)

	if err == nil && upvote {
		voteStatus.Upvoted = true
	}
	if err == nil && !upvote {
		voteStatus.Downvoted = true
	}

	td.SolutionVoteStatus = voteStatus

	app.render(w, r, "problem/solution.gohtml", td)
}

func (app *Application) handleSolutionVote(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())
	idString := ps.ByName("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	td := app.newTemplateData(r)

	if td.AuthenticatedUserId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	var voteReq entities.SolutionVoteRequest

	if err = json.NewDecoder(r.Body).Decode(&voteReq); err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	switch {
	case voteReq.Upvoted:
		err = app.ProblemService.UpvoteSolution(id, td.AuthenticatedUserId)
		if err != nil {
			app.ErrorLog.Print(err)
			app.serverError(w, r)
		}
	case voteReq.Downvoted:
		err = app.ProblemService.DownvoteSolution(id, td.AuthenticatedUserId)
		if err != nil {
			app.ErrorLog.Print(err)
			app.serverError(w, r)
		}
	default:
		err = app.ProblemService.UnvoteSolution(id, td.AuthenticatedUserId)
		if err != nil {
			app.ErrorLog.Print(err)
			app.serverError(w, r)
		}
	}

	var solution *entities.ProblemSolution
	solution, err = app.ProblemService.GetSolutionById(id)

	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}

	response := entities.SolutionVoteResponse{
		Upvoted:   voteReq.Upvoted,
		Downvoted: voteReq.Downvoted,
		Votes:     solution.Votes,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
	}
}

func (app *Application) createSolutionComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := app.Session.GetInt(r.Context(), "authenticatedUserId")
	if userId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	solutionId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid solution ID", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	text := r.PostFormValue("text")
	if text == "" {
		http.Error(w, "Comment text is required", http.StatusBadRequest)
		return
	}

	parentCommentID := r.PostFormValue("parent_comment_id")
	var parentID *int
	if parentCommentID != "" {
		id, err := strconv.Atoi(parentCommentID)
		if err != nil {
			app.ErrorLog.Print(err)
			http.Error(w, "Invalid parent comment ID", http.StatusBadRequest)
			return
		}
		parentID = &id
	}

	if parentID != nil {
		err = app.ProblemSolutionCommentService.CreateReply(userId, solutionId, text, parentID)
	} else {
		err = app.ProblemSolutionCommentService.CreateComment(userId, solutionId, text)
	}

	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/solutions/%d", solutionId), http.StatusSeeOther)
}

func (app *Application) editSolutionComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := app.Session.GetInt(r.Context(), "authenticatedUserId")
	if userId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	commentId, err := strconv.Atoi(r.PostFormValue("comment_id"))
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	text := r.PostFormValue("text")
	if text == "" {
		http.Error(w, "Comment text is required", http.StatusBadRequest)
		return
	}

	err = app.ProblemSolutionCommentService.UpdateComment(commentId, text)
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	// Redirect back to the previous page
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *Application) deleteSolutionComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := app.Session.GetInt(r.Context(), "authenticatedUserId")
	if userId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	commentId, err := strconv.Atoi(r.PostFormValue("comment_id"))
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	err = app.ProblemSolutionCommentService.DeleteComment(commentId)
	if err != nil {
		app.ErrorLog.Print(err)
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	// Redirect back to the previous page
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
