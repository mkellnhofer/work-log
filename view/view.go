package view

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"golang.org/x/text/message"

	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/loc"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/view/model"
)

const dateFormat = "02.01.2006"
const dateFormatShort = "02.01."
const timeFormat = "15:04"

var weekdayKeys = map[int]string{
	0: "weekdaySun",
	1: "weekdayMon",
	2: "weekdayTue",
	3: "weekdayWed",
	4: "weekdayThu",
	5: "weekdayFri",
	6: "weekdaySat",
}

var monthKeys = map[int]string{
	1:  "monthJan",
	2:  "monthFeb",
	3:  "monthMar",
	4:  "monthApr",
	5:  "monthMay",
	6:  "monthJun",
	7:  "monthJul",
	8:  "monthAug",
	9:  "monthSep",
	10: "monthOct",
	11: "monthNov",
	12: "monthDec",
}

// --- Template loading functions ---

var templates = loadTemplates("templates/header.tmpl", "templates/footer.tmpl",
	"templates/error.tmpl", "templates/login.tmpl", "templates/entries_list.tmpl",
	"templates/list_entries.tmpl", "templates/entry_form.tmpl", "templates/create_entry.tmpl",
	"templates/edit_entry.tmpl", "templates/copy_entry.tmpl", "templates/search_entries.tmpl",
	"templates/list_search_entries.tmpl", "templates/list_overview_entries.tmpl")

func loadTemplates(filenames ...string) *template.Template {
	var t *template.Template
	for _, filename := range filenames {
		t = loadTemplate(t, filename)
	}
	return t
}

func loadTemplate(t *template.Template, filename string) *template.Template {
	// Read template
	b, rErr := ioutil.ReadFile(filename)
	if rErr != nil {
		err := e.WrapError(e.SysUnknown, fmt.Sprintf("Could load template '%s'.", filename), rErr)
		log.Debug(err.StackTrace())
		panic(err)
	}
	s := string(b)
	name := filepath.Base(filename)

	// Register template
	var tmpl *template.Template
	if t == nil {
		t = template.New(name)
		tmpl = t
	} else {
		tmpl = t.New(name)
	}

	// Add functions
	tmpl.Funcs(templateFuncs)

	// Parse template
	_, pErr := tmpl.Parse(s)
	if pErr != nil {
		err := e.WrapError(e.SysUnknown, fmt.Sprintf("Could parse template '%s'.", filename), pErr)
		log.Debug(err.StackTrace())
		panic(err)
	}

	return t
}

// --- Template functions ---

var templateFuncs = template.FuncMap{"text": getText}

func getText(key string) string {
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf(key)
}

// --- Render functions ---

// RenderErrorTemplate renders the error page.
func RenderErrorTemplate(w http.ResponseWriter, model *model.Error) {
	renderTemplate(w, "error", model)
}

// RenderLoginTemplate renders the login page.
func RenderLoginTemplate(w http.ResponseWriter, model *model.Login) {
	renderTemplate(w, "login", model)
}

// RenderListEntriesTemplate renders a page of entries.
func RenderListEntriesTemplate(w http.ResponseWriter, model *model.ListEntries) {
	renderTemplate(w, "list_entries", model)
}

// RenderCreateEntryTemplate renders the page to create a entry.
func RenderCreateEntryTemplate(w http.ResponseWriter, model *model.CreateEntry) {
	renderTemplate(w, "create_entry", model)
}

// RenderEditEntryTemplate renders the page to edit a entry.
func RenderEditEntryTemplate(w http.ResponseWriter, model *model.EditEntry) {
	renderTemplate(w, "edit_entry", model)
}

// RenderCopyEntryTemplate renders the page to copy a entry.
func RenderCopyEntryTemplate(w http.ResponseWriter, model *model.CopyEntry) {
	renderTemplate(w, "copy_entry", model)
}

// RenderSearchEntriesTemplate renders the page to search entries.
func RenderSearchEntriesTemplate(w http.ResponseWriter, model *model.SearchEntries) {
	renderTemplate(w, "search_entries", model)
}

// RenderListSearchEntriesTemplate renders a page with searched entries.
func RenderListSearchEntriesTemplate(w http.ResponseWriter, model *model.ListSearchEntries) {
	renderTemplate(w, "list_search_entries", model)
}

// RenderListOverviewEntriesTemplate renders the page to overview entries.
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

// FormatHours returns the hours string for a duration.
func FormatHours(d time.Duration) string {
	h := d.Hours()
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf("%.2f", h)
}

// GetWeekdayName returns the weekday string for a time.
func GetWeekdayName(t time.Time) string {
	wd := t.Weekday()
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf(weekdayKeys[int(wd)])
}

// GetShortWeekdayName returns the shortend weekday string for a time.
func GetShortWeekdayName(t time.Time) string {
	d := GetWeekdayName(t)
	return fmt.Sprintf("%s.", d[0:2])
}

// GetMonthName returns the name of a month.
func GetMonthName(m int) string {
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf(monthKeys[m])
}

// --- Helper functions ---

func renderTemplate(w http.ResponseWriter, tmpl string, model interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", model)
	if err != nil {
		log.Errorf("Could not render template! %s", err)
		http.Error(w, "Could not render template!", http.StatusInternalServerError)
	}
}
