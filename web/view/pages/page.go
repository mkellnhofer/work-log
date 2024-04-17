package pages

import (
	"github.com/a-h/templ"

	"kellnhofer.com/work-log/web/view"
)

func getText(key string) string {
	return view.GetText(key)
}

func toURL(url string) templ.SafeURL {
	return templ.URL(url)
}
