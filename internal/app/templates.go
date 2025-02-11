package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type templateData struct {
	CurrentYear         int
	Form                any
	AuthenticatedUserId int
	Problems            []*entities.Problem
	Topics              []*entities.Topic
	Pagination          Pagination
	Examples            []entities.ProblemExample
	UserErrorMessage    string
	Users               []*entities.User
	Submissions         []*entities.FullProblemSubmission
}

func (app *Application) newTemplateData(r *http.Request) *templateData {
	userIdFromSession := app.Session.Get(r.Context(), "authenticatedUserId")

	authenticatedUserId, ok := userIdFromSession.(int)
	if !ok {
		authenticatedUserId = 0
	}

	return &templateData{
		CurrentYear:         time.Now().Year(),
		AuthenticatedUserId: authenticatedUserId,
	}
}

func (app *Application) InitTemplates() error {
	cache := map[string]*template.Template{}

	// Get the base layout template
	baseTemplate := "./web/templates/layout/base.gohtml"

	partials, err := filepath.Glob("./web/templates/partials/*.gohtml")
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"toLowerCase": strings.ToLower,
		"add":         func(a, b int) int { return a + b },
		"subtract":    func(a, b int) int { return a - b },
		"formatJSON": func(v interface{}) string {
			if v == nil {
				return ""
			}
			switch v := v.(type) {
			case json.RawMessage:
				var prettyJSON bytes.Buffer
				err := json.Indent(&prettyJSON, v, "", "  ")
				if err != nil {
					return string(v)
				}
				return prettyJSON.String()
			default:
				data, err := json.MarshalIndent(v, "", "  ")
				if err != nil {
					return fmt.Sprintf("%v", v)
				}
				return string(data)
			}
		},
	}

	// Get all page templates
	pages, err := filepath.Glob("./web/templates/*/*.gohtml")
	if err != nil {
		return err
	}

	// Create template for each page
	for _, page := range pages {
		name := filepath.Base(page)

		// Skip base.html as it's our layout template
		if name == "base.gohtml" {
			continue
		}

		// Create template set with base template
		ts, err := template.New("").Funcs(funcMap).ParseFiles(baseTemplate)
		if err != nil {
			return err
		}

		// Parse all partials
		ts, err = ts.ParseFiles(partials...)
		if err != nil {
			return err
		}

		// Add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return err
		}

		cache[name] = ts
	}

	app.templateCache = cache
	return nil
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Get requested template from cache
	ts, ok := app.templateCache[filepath.Base(name)]
	if !ok {
		app.ErrorLog.Printf("Template %s not found", name)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create template data if none provided
	if td == nil {
		td = app.newTemplateData(r)
	}

	// Execute the template
	err := ts.ExecuteTemplate(w, "base", td)
	if err != nil {
		app.ErrorLog.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
