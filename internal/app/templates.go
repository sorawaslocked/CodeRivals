package app

import (
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear         int
	Form                any
	AuthenticatedUserId uint64
	Problems            []*entities.Problem
}

func (app *Application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:         time.Now().Year(),
		AuthenticatedUserId: app.Session.Get(r.Context(), "authenticatedUserId").(uint64),
	}
}

func (app *Application) InitTemplates() error {
	cache := map[string]*template.Template{}

	// Get the base layout template
	baseTemplate := "./web/templates/layout/base.html"

	// Get all page templates
	pages, err := filepath.Glob("./web/templates/*/*.html")
	if err != nil {
		return err
	}

	// Create template for each page
	for _, page := range pages {
		name := filepath.Base(page)

		// Skip base.html as it's our layout template
		if name == "base.html" {
			continue
		}

		// Create template set with base template
		ts, err := template.ParseFiles(baseTemplate)
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
