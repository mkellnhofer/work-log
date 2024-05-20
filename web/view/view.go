package view

import (
	"strconv"
	"strings"

	"golang.org/x/text/message"

	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/web/model"
)

// GetText returns a localized text.
func GetText(key string) string {
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf(key)
}

// --- Functions to render the user icon SVG ---

var userIconColors = []string{
	"#df5755", "#e24970", "#977dc8", "#439dde", "#3eaabd",
	"#57aa5a", "#8dbe5a", "#f7c04c", "#f7a951", "#f76b4f",
}

// CreateUserIconSvg creates a user icon from the user's initials.
func CreateUserIconSvg(initials string) string {
	color := userIconColors[(int(initials[0])+int(initials[1]))%len(userIconColors)]

	return createSvg("0 0 32 32", "32", "32", func() string {
		return `
			<circle
				style="fill:` + color + `;"
				cx="16"
				cy="16"
				r="16"
			/>
			<text
				style="font-size:0.7em;text-anchor:middle;fill:#FFF;"
				x="16"
				y="20"
			>` + initials + `</text>
		`
	})
}

// --- Functions to render the log summary progress SVG ---

var LogSummaryProgressColorLog = "#999999"
var LogSummaryProgressColorRem = "#d6d6d6"
var LogSummaryProgressColorOvr = "#0c63e4"
var LogSummaryProgressColorUnd = "#f1b523"

// CreateLogSummaryProgressSvg creates a log summary progress bar.
func CreateLogSummaryProgressSvg(logged int, overtime int, undertime int) string {
	return createProgressSvg(func() string {
		return createLogSummaryProgressSvgLogRem(logged) +
			createLogSummaryProgressSvgOvrUnd(logged, overtime, undertime)
	})
}

func createLogSummaryProgressSvgLogRem(logged int) string {
	return createProgressSvgRect(0, 100, LogSummaryProgressColorRem) +
		createProgressSvgRect(0, logged, LogSummaryProgressColorLog)
}

func createLogSummaryProgressSvgOvrUnd(logged int, overtime int, undertime int) string {
	if overtime > undertime {
		return createProgressSvgRect(logged-overtime, logged, LogSummaryProgressColorOvr)
	} else {
		return createProgressSvgRect(logged, logged+undertime, LogSummaryProgressColorUnd)
	}
}

// --- Functions to render the overview summary progress SVG ---

var OverviewSummaryProgressColorRem = "#d6d6d6"

// CreateOverviewSummaryProgressSvg creates a overview summary progress bar.
func CreateOverviewSummaryProgressSvg(typeProgress map[int]int) string {
	return createProgressSvg(func() string {
		ws, we := 0, typeProgress[model.EntryTypeIdWork]
		ts, te := we, we+typeProgress[model.EntryTypeIdTravel]
		vs, ve := te, te+typeProgress[model.EntryTypeIdVacation]
		hs, he := ve, ve+typeProgress[model.EntryTypeIdHoliday]
		is, ie := he, he+typeProgress[model.EntryTypeIdIllness]
		return createProgressSvgRect(0, 100, OverviewSummaryProgressColorRem) +
			createProgressSvgRect(ws, we, model.EntryTypeColors[model.EntryTypeIdWork]) +
			createProgressSvgRect(ts, te, model.EntryTypeColors[model.EntryTypeIdTravel]) +
			createProgressSvgRect(vs, ve, model.EntryTypeColors[model.EntryTypeIdVacation]) +
			createProgressSvgRect(hs, he, model.EntryTypeColors[model.EntryTypeIdHoliday]) +
			createProgressSvgRect(is, ie, model.EntryTypeColors[model.EntryTypeIdIllness])
	})
}

// --- Progress SVG functions ---

func createProgressSvg(bodyFunc func() string) string {
	return createSvg("0 0 1000 10", "100%", "10", func() string {
		return createProgressSvgMask() + bodyFunc()
	})
}

func createProgressSvgMask() string {
	return `
		<mask id="msk">
			<rect width="100%" height="100%" fill="black" />
			<rect width="100%" height="100%" ry="100%" fill="white" />
		</mask>
	`
}

func createProgressSvgRect(start int, end int, color string) string {
	return `
		<rect
			x="` + strconv.Itoa(start) + `%"
			width="` + strconv.Itoa(end-start) + `%"
			height="100%"
			fill="` + color + `"
			mask="url(#msk)" />
	`
}

// --- Helper functions ---

func createSvg(viewbox string, width string, height string, bodyFunc func() string) string {
	return trimSvg(`
		<svg
			viewbox="` + viewbox + `"
			width="` + width + `"
			height="` + height + `"
			preserveAspectRatio="none"
			xmlns="http://www.w3.org/2000/svg"
		>
			` + bodyFunc() + `
		</svg>
	`)
}

func trimSvg(svg string) string {
	svg = strings.ReplaceAll(svg, "\n", " ")
	svg = strings.ReplaceAll(svg, "\t", "")
	return svg
}
