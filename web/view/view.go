package view

import (
	"golang.org/x/text/message"

	"kellnhofer.com/work-log/pkg/loc"
)

// GetText returns a localized text.
func GetText(key string) string {
	printer := message.NewPrinter(loc.LngTag)
	return printer.Sprintf(key)
}
