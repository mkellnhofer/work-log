package web

import (
	"fmt"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/message"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
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

// --- Template functions ---

func GetText(key string) string {
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf(key)
}

// --- Render functions ---

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	tErr := t.Render(ctx.Request().Context(), ctx.Response().Writer)
	if tErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not render template.", tErr)
		log.Debug(err.StackTrace())
		return err
	}
	return nil
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
