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

func toURL(url string) templ.SafeURL {
	return templ.URL(url)
}
