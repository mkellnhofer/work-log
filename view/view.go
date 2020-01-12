package view

import (
	"html/template"
	"net/http"

	"kellnhofer.com/work-log/view/model"
)

var templates = template.Must(template.ParseFiles("templates/header.tmpl", "templates/footer.tmpl",
	"templates/error.tmpl"))

// RenderErrorTemplate renders the error page.
func RenderErrorTemplate(w http.ResponseWriter, model *model.Error) {
	renderTemplate(w, "error", model)
}

// --- Helper methods ---

func renderTemplate(w http.ResponseWriter, tmpl string, model interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", model)
	if err != nil {
		http.Error(w, "Failed to render template!", http.StatusInternalServerError)
	}
}