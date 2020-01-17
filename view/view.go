package view

import (
	"html/template"
	"net/http"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/view/model"
)

const dateFormat = "02.01.2006"
const timeFormat = "15:04"

var weekdays = map[int]string{
	0: "Sonntag",
	1: "Montag",
	2: "Dienstag",
	3: "Mittwoch",
	4: "Donnerstag",
	5: "Freitag",
	6: "Samstag",
}

var printer = message.NewPrinter(language.German)

var templates = template.Must(template.ParseFiles("templates/header.tmpl", "templates/footer.tmpl",
	"templates/error.tmpl", "templates/login.tmpl", "templates/list_entries.tmpl"))

// --- Render functions ---

// RenderErrorTemplate renders the error page.
func RenderErrorTemplate(w http.ResponseWriter, model *model.Error) {
	renderTemplate(w, "error", model)
}

// RenderLoginTemplate renders the login page.
func RenderLoginTemplate(w http.ResponseWriter, model *model.Login) {
	renderTemplate(w, "login", model)
}

// RenderListEntriesTemplate renders the page to list work entries.
func RenderListEntriesTemplate(w http.ResponseWriter, model *model.ListEntries) {
	renderTemplate(w, "list_entries", model)
}

// --- Time formatting functions ---

// FormatDate returns the date string for a time.
func FormatDate(t time.Time) string {
	return t.Format(dateFormat)
}

// FormatTime returns the time string for a time.
func FormatTime(t time.Time) string {
	return t.Format(timeFormat)
}

// FormatWeekday returns the weekday string for a time.
func FormatWeekday(t time.Time) string {
	wd := t.Weekday()
	return weekdays[int(wd)]
}

// FormatHours returns the hours string for a duration.
func FormatHours(d time.Duration) string {
	h := d.Hours()
	return printer.Sprintf("%.2f", h)
}

// --- Helper functions ---

func renderTemplate(w http.ResponseWriter, tmpl string, model interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", model)
	if err != nil {
		log.Errorf("Could not render template! %s", err)
		http.Error(w, "Could not render template!", http.StatusInternalServerError)
	}
}
