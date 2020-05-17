package view

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/view/model"
)

const dateFormat = "02.01.2006"
const dateFormatShort = "02.01."
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

var months = map[int]string{
	1:  "Januar",
	2:  "Februar",
	3:  "März",
	4:  "April",
	5:  "Mai",
	6:  "Juni",
	7:  "Juli",
	8:  "August",
	9:  "September",
	10: "Oktober",
	11: "November",
	12: "Dezember",
}

var printer = message.NewPrinter(language.German)

var templates = template.Must(template.ParseFiles("templates/header.tmpl", "templates/footer.tmpl",
	"templates/error.tmpl", "templates/login.tmpl", "templates/entries_list.tmpl",
	"templates/list_entries.tmpl", "templates/entry_form.tmpl", "templates/create_entry.tmpl",
	"templates/edit_entry.tmpl", "templates/copy_entry.tmpl", "templates/search_entries.tmpl",
	"templates/list_search_entries.tmpl", "templates/list_overview_entries.tmpl"))

// --- Render functions ---

// RenderErrorTemplate renders the error page.
func RenderErrorTemplate(w http.ResponseWriter, model *model.Error) {
	renderTemplate(w, "error", model)
}

// RenderLoginTemplate renders the login page.
func RenderLoginTemplate(w http.ResponseWriter, model *model.Login) {
	renderTemplate(w, "login", model)
}

// RenderListEntriesTemplate renders a page of work entries.
func RenderListEntriesTemplate(w http.ResponseWriter, model *model.ListEntries) {
	renderTemplate(w, "list_entries", model)
}

// RenderCreateEntryTemplate renders the page to create a work entry.
func RenderCreateEntryTemplate(w http.ResponseWriter, model *model.CreateEntry) {
	renderTemplate(w, "create_entry", model)
}

// RenderEditEntryTemplate renders the page to edit a work entry.
func RenderEditEntryTemplate(w http.ResponseWriter, model *model.EditEntry) {
	renderTemplate(w, "edit_entry", model)
}

// RenderCopyEntryTemplate renders the page to copy a work entry.
func RenderCopyEntryTemplate(w http.ResponseWriter, model *model.CopyEntry) {
	renderTemplate(w, "copy_entry", model)
}

// RenderSearchEntriesTemplate renders the page to search work entries.
func RenderSearchEntriesTemplate(w http.ResponseWriter, model *model.SearchEntries) {
	renderTemplate(w, "search_entries", model)
}

// RenderListSearchEntriesTemplate renders a page with searched work entries.
func RenderListSearchEntriesTemplate(w http.ResponseWriter, model *model.ListSearchEntries) {
	renderTemplate(w, "list_search_entries", model)
}

// RenderListOverviewEntriesTemplate renders the page to overview work entries.
func RenderListOverviewEntriesTemplate(w http.ResponseWriter, model *model.ListOverviewEntries) {
	renderTemplate(w, "list_overview_entries", model)
}

// --- Time formatting functions ---

// FormatDate returns the date string for a time.
func FormatDate(t time.Time) string {
	return t.Format(dateFormat)
}

// FormatShortDate returns the short date string for a time.
func FormatShortDate(t time.Time) string {
	return t.Format(dateFormatShort)
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

// FormatShortWeekday returns the shortend weekday string for a time.
func FormatShortWeekday(t time.Time) string {
	wd := t.Weekday()
	d := weekdays[int(wd)]
	return fmt.Sprintf("%s.", d[0:2])
}

// FormatHours returns the hours string for a duration.
func FormatHours(d time.Duration) string {
	h := d.Hours()
	return printer.Sprintf("%.2f", h)
}

// GetMonthName returns the name of a month.
func GetMonthName(m int) string {
	return months[m]
}

// --- Helper functions ---

func renderTemplate(w http.ResponseWriter, tmpl string, model interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", model)
	if err != nil {
		log.Errorf("Could not render template! %s", err)
		http.Error(w, "Could not render template!", http.StatusInternalServerError)
	}
}
