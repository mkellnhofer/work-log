package mapper

import (
	"fmt"
	"time"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
)

func parseDate(d string) time.Time {
	t, pErr := time.Parse(constant.ApiDateFormat, d)
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
	t, pErr := time.Parse(constant.ApiTimestampFormat, ts)
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

func parseHoursDuration(h float32) time.Duration {
	min := int(h * 60.0)
	d, pErr := time.ParseDuration(fmt.Sprintf("%dm", min))
	if pErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not parse hours duration.", pErr)
		log.Error(err.StackTrace())
		panic(err)
	}
	return d
}

func formatHoursDuration(d time.Duration) float32 {
	md := d.Round(time.Minute)
	min := int(md.Minutes())
	return float32(min) / 60.0
}
