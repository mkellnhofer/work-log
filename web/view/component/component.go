package component

import (
	"strconv"

	"github.com/a-h/templ"
	"kellnhofer.com/work-log/web/view"
)

func getText(key string) string {
	return view.GetText(key)
}

func toString(value int) string {
	return strconv.Itoa(value)
}

func toQuotedString(value string) string {
	return "\"" + value + "\""
}

func toQuotedStrings(values []string) []string {
	quotedValues := make([]string, len(values))
	for i, value := range values {
		quotedValues[i] = toQuotedString(value)
	}
	return quotedValues
}

func toURL(url string) templ.SafeURL {
	return templ.URL(url)
}

func hx(url string) string {
	return "/hx" + url
}

func createColorStyleAttributes(color string) templ.Attributes {
	return templ.Attributes{"style": "color:" + color + " !important;"}
}

func createBorderColorStyleAttributes(color string) templ.Attributes {
	return templ.Attributes{"style": "border-color:" + color + " !important;"}
}
