package mapper

import (
	"time"

	"kellnhofer.com/work-log/pkg/constant"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
)

func parseDate(d string) time.Time {
	t, pErr := time.ParseInLocation(constant.ApiDateFormat, d, time.Local)
	if pErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not parse date.", pErr)
		log.Error(err.StackTrace())
		panic(err)
	}
	return t
}

func formatDate(t time.Time) string {
	return t.Format(constant.ApiDateFormat)
}

func parseTimestamp(ts string) time.Time {
	t, pErr := time.ParseInLocation(constant.ApiTimestampFormat, ts, time.Local)
	if pErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not parse timestamp.", pErr)
		log.Error(err.StackTrace())
		panic(err)
	}
	return t
}

func formatTimestamp(t time.Time) string {
	return t.Format(constant.ApiTimestampFormat)
}
