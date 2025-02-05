package app

import "net/http"

func (app *Application) notFound(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "errors/notfound.gohtml", nil)
}
