package view

import (
	"strings"

	"golang.org/x/text/message"

	"kellnhofer.com/work-log/pkg/loc"
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
				style="font-size:13;font-family:FreeSans;text-anchor:middle;fill:#FFF;"
				x="16"
				y="21"
			>` + initials + `</text>
		`
	})
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
