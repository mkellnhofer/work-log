package web

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/message"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/web/model"
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

var templates = loadTemplates("header.tmpl", "footer.tmpl", "error.tmpl", "login.tmpl",
	"password_change.tmpl", "entries_list.tmpl", "list_entries.tmpl", "entry_form.tmpl",
	"create_entry.tmpl", "edit_entry.tmpl", "copy_entry.tmpl", "search_entries.tmpl",
	"list_search_entries.tmpl", "list_overview_entries.tmpl")

func loadTemplates(filenames ...string) *template.Template {
	var t *template.Template
	for _, filename := range filenames {
		t = loadTemplate(t, "web/templates/"+filename)
	}
	return t
}

func loadTemplate(t *template.Template, filename string) *template.Template {
	// Read template
	b, rErr := os.ReadFile(filename)
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
func RenderErrorTemplate(r *echo.Response, model *model.Error) error {
	return renderTemplate(r, "error", model)
}

// RenderLoginTemplate renders the login page.
func RenderLoginTemplate(r *echo.Response, model *model.Login) error {
	return renderTemplate(r, "login", model)
}

// RenderPasswordChangeTemplate renders the password change page.
func RenderPasswordChangeTemplate(r *echo.Response, model *model.PasswordChange) error {
	return renderTemplate(r, "password_change", model)
}

// RenderListEntriesTemplate renders a page of entries.
func RenderListEntriesTemplate(r *echo.Response, model *model.ListEntries) error {
	return renderTemplate(r, "list_entries", model)
}

// RenderCreateEntryTemplate renders the page to create a entry.
func RenderCreateEntryTemplate(r *echo.Response, model *model.CreateEntry) error {
	return renderTemplate(r, "create_entry", model)
}

// RenderEditEntryTemplate renders the page to edit a entry.
func RenderEditEntryTemplate(r *echo.Response, model *model.EditEntry) error {
	return renderTemplate(r, "edit_entry", model)
}

// RenderCopyEntryTemplate renders the page to copy a entry.
func RenderCopyEntryTemplate(r *echo.Response, model *model.CopyEntry) error {
	return renderTemplate(r, "copy_entry", model)
}

// RenderSearchEntriesTemplate renders the page to search entries.
func RenderSearchEntriesTemplate(r *echo.Response, model *model.SearchEntries) error {
	return renderTemplate(r, "search_entries", model)
}

// RenderListSearchEntriesTemplate renders a page with searched entries.
func RenderListSearchEntriesTemplate(r *echo.Response, model *model.ListSearchEntries) error {
	return renderTemplate(r, "list_search_entries", model)
}

// RenderListOverviewEntriesTemplate renders the page to overview entries.
func RenderListOverviewEntriesTemplate(r *echo.Response, model *model.ListOverviewEntries) error {
	return renderTemplate(r, "list_overview_entries", model)
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

// FormatMinutes returns the minutes string for a duration.
func FormatMinutes(d time.Duration) string {
	m := d.Minutes()
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf("%d", int(m))
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

// CreateDaysString return a formated days string
func CreateDaysString(days float32) string {
	return loc.CreateString("daysValue", days)
}

// CreateHoursString return a formated hours string
func CreateHoursString(hours float32) string {
	return loc.CreateString("hoursValue", hours)
}

// --- Helper functions ---

func renderTemplate(r *echo.Response, tmpl string, model interface{}) error {
	rw := r.Writer
	tErr := templates.ExecuteTemplate(rw, tmpl+".tmpl", model)
	if tErr != nil {
		err := e.WrapError(e.SysUnknown, fmt.Sprintf("Could not render template '%s'.", tmpl), tErr)
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}
