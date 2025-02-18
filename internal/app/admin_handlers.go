package app

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/sorawaslocked/CodeRivals/internal/dtos"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"net/http"
	"strconv"
	"strings"
)

func (app *Application) adminProblems(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	if td.AuthenticatedUserId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	isAdmin, err := app.AdminService.IsUserAdmin(td.AuthenticatedUserId)
	if err != nil || !isAdmin {
		app.notFound(w, r)
		return
	}

	problems, err := app.ProblemService.GetAllProblems()
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}
	td.Problems = problems

	topics, err := app.TopicService.GetAllTopics()
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}
	td.Topics = topics

	app.render(w, r, "admin/admin_problems.gohtml", td)
}

func (app *Application) adminGetProblem(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	problem, err := app.ProblemService.GetProblem(id)
	if err != nil {
		http.Error(w, "Problem not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problem)
}

func (app *Application) adminCreateProblem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get topics
	var topics []*entities.Topic
	topicIDs := r.Form["topics"]
	for _, idStr := range topicIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		topics = append(topics, &entities.Topic{ID: id})
	}

	// Parse input types
	inputTypesStr := r.FormValue("inputTypes")
	inputTypes := strings.Split(inputTypesStr, ",")
	for i := range inputTypes {
		inputTypes[i] = strings.TrimSpace(inputTypes[i])
	}

	title := r.FormValue("title")
	problem := &entities.Problem{
		Title:       title,
		Description: r.FormValue("description"),
		Difficulty:  r.FormValue("difficulty"),
		Topics:      topics,
		InputTypes:  inputTypes,
		OutputType:  r.FormValue("outputType"),
		MethodName:  r.FormValue("methodName"),
		Url:         createURLFromTitle(title),
	}

	createReq := &dtos.ProblemCreateRequest{
		Title:       problem.Title,
		Description: problem.Description,
		Difficulty:  problem.Difficulty,
		Topics:      problem.Topics,
		InputTypes:  problem.InputTypes,
		OutputType:  problem.OutputType,
		MethodName:  problem.MethodName,
		Url:         problem.Url,
	}

	if err := app.ProblemService.CreateProblem(createReq); err != nil {
		app.ErrorLog.Printf("Failed to create problem: %v", err)
		http.Error(w, "Failed to create problem", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/problems", http.StatusSeeOther)
}

func (app *Application) adminUpdateProblem(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get topics
	var topics []*entities.Topic
	topicIDs := r.Form["topics"]
	for _, idStr := range topicIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		topics = append(topics, &entities.Topic{ID: id})
	}

	// Parse input types
	inputTypesStr := r.FormValue("inputTypes")
	inputTypes := strings.Split(inputTypesStr, ",")
	for i := range inputTypes {
		inputTypes[i] = strings.TrimSpace(inputTypes[i])
	}

	title := r.FormValue("title")
	problem := &entities.Problem{
		ID:          id,
		Title:       title,
		Description: r.FormValue("description"),
		Difficulty:  r.FormValue("difficulty"),
		Topics:      topics,
		InputTypes:  inputTypes,
		OutputType:  r.FormValue("outputType"),
		MethodName:  r.FormValue("methodName"),
		Url:         createURLFromTitle(title),
	}

	if err := app.ProblemService.UpdateProblem(problem); err != nil {
		app.ErrorLog.Printf("Failed to update problem: %v", err)
		http.Error(w, "Failed to update problem", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/problems", http.StatusSeeOther)
}

func (app *Application) adminDeleteProblem(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	if err := app.ProblemService.DeleteProblem(id); err != nil {
		app.ErrorLog.Printf("Failed to delete problem: %v", err)
		http.Error(w, "Failed to delete problem", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *Application) adminGetTestCases(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	testCases, err := app.ProblemService.GetTestCases(id)
	if err != nil {
		app.ErrorLog.Printf("Failed to get test cases: %v", err)
		http.Error(w, "Failed to get test cases", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testCases)
}

func (app *Application) adminUpdateTestCases(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	problemID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.ErrorLog.Printf("Invalid problem ID: %v", err)
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	// Parse multipart form data with a reasonable max memory
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		app.ErrorLog.Printf("Failed to parse multipart form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	inputs := r.MultipartForm.Value["inputs[]"]
	outputs := r.MultipartForm.Value["outputs[]"]

	// Log the received data
	app.InfoLog.Printf("Received inputs: %v", inputs)
	app.InfoLog.Printf("Received outputs: %v", outputs)

	if len(inputs) != len(outputs) {
		app.ErrorLog.Printf("Mismatched inputs (%d) and outputs (%d)", len(inputs), len(outputs))
		http.Error(w, "Mismatched inputs and outputs", http.StatusBadRequest)
		return
	}

	// Create empty test cases array - might be empty if all test cases were removed
	var testCases []*entities.ProblemTestCase

	// Only process test cases if there are any
	for i := range inputs {
		testCase := &entities.ProblemTestCase{
			ProblemID:  problemID,
			OrderIndex: i,
			Input:      inputs[i],
			Output:     outputs[i],
		}
		testCases = append(testCases, testCase)
	}

	// Log operation
	if len(testCases) == 0 {
		app.InfoLog.Printf("Removing all test cases for problem %d", problemID)
	} else {
		app.InfoLog.Printf("Saving %d test cases for problem %d", len(testCases), problemID)
	}

	// Update test cases (might be empty array = removing all test cases)
	err = app.ProblemService.UpdateTestCases(problemID, testCases)
	if err != nil {
		app.ErrorLog.Printf("Failed to update test cases: %v", err)
		http.Error(w, "Failed to update test cases", http.StatusInternalServerError)
		return
	}

	app.InfoLog.Printf("Successfully updated test cases for problem %d", problemID)
	w.WriteHeader(http.StatusOK)
}

func createURLFromTitle(title string) string {
	// Convert to lowercase and replace spaces with hyphens
	url := strings.ToLower(title)
	url = strings.ReplaceAll(url, " ", "-")

	// Remove any special characters
	url = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, url)

	return url
}

func (app *Application) adminUsers(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	if td.AuthenticatedUserId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	isAdmin, err := app.AdminService.IsUserAdmin(td.AuthenticatedUserId)
	if err != nil || !isAdmin {
		app.notFound(w, r)
		return
	}

	users, err := app.AdminService.GetAllUsersWithDetails()
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}

	td.AdminUsers = users

	app.render(w, r, "admin/admin_users.gohtml", td)
}

func (app *Application) adminGetUserDetails(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	details, err := app.AdminService.GetUserDetails(userId)
	if err != nil {
		http.Error(w, "Failed to get user details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

func (app *Application) adminUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = app.AdminService.UpdateUserRole(userId, req.Role == "admin")
	if err != nil {
		http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *Application) adminResetUserPassword(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	newPassword, err := app.AdminService.ResetUserPassword(userId)
	if err != nil {
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"password": newPassword})
}

func (app *Application) adminTopics(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	if td.AuthenticatedUserId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	isAdmin, err := app.AdminService.IsUserAdmin(td.AuthenticatedUserId)
	if err != nil || !isAdmin {
		app.notFound(w, r)
		return
	}

	topics, err := app.TopicService.GetAllTopics()
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}

	var topicsWithCount []*entities.TopicWithCount
	for _, topic := range topics {
		count, err := app.TopicService.GetProblemCountForTopic(topic.ID)
		if err != nil {
			app.ErrorLog.Print(err)
			continue
		}
		topicsWithCount = append(topicsWithCount, &entities.TopicWithCount{
			Topic:         topic,
			ProblemsCount: count,
		})
	}

	td.TopicsWithCount = topicsWithCount
	app.render(w, r, "admin/admin_topics.gohtml", td)
}

func (app *Application) adminTopicsCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}

	name := r.PostFormValue("name")
	if name == "" {
		app.userError(w, r, "Topic name cannot be empty")
		return
	}

	err = app.TopicService.Create(name)
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}

	http.Redirect(w, r, "/admin/topics", http.StatusSeeOther)
}

func (app *Application) adminTopicsUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params := httprouter.ParamsFromContext(r.Context())
	topicId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.notFound(w, r)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}

	name := r.PostFormValue("name")
	if name == "" {
		app.userError(w, r, "Topic name cannot be empty")
		return
	}

	err = app.TopicService.UpdateTopic(topicId, name)
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}

	http.Redirect(w, r, "/admin/topics", http.StatusSeeOther)
}

func (app *Application) adminTopicsDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params := httprouter.ParamsFromContext(r.Context())
	topicId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.notFound(w, r)
		return
	}

	err = app.TopicService.Delete(topicId)
	if err != nil {
		app.ErrorLog.Print(err)
		app.serverError(w, r)
		return
	}

	http.Redirect(w, r, "/admin/topics", http.StatusSeeOther)
}
