package page

import (
	"kellnhofer.com/work-log/web/view"
)

func getText(key string) string {
	return view.GetText(key)
}

func hx(url string) string {
	return "/hx/" + url
}
