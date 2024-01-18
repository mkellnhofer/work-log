package mapper

import (
	"fmt"
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

func parseMinutesDuration(m int) time.Duration {
	d, pErr := time.ParseDuration(fmt.Sprintf("%dm", m))
	if pErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not parse minutes duration.", pErr)
		log.Error(err.StackTrace())
		panic(err)
	}
	return d
}

func formatMinutesDuration(d time.Duration) int {
	md := d.Round(time.Minute)
	return int(md.Minutes())
}
