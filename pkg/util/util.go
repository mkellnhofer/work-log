package util

import (
	"encoding/base64"
	"math"
	"time"
)

// EncodeBase64 encodes a string into a Base64 string.
func EncodeBase64(s string) string {
	b := []byte(s)
	return base64.URLEncoding.EncodeToString(b)
}

// DecodeBase64 decodes a Base64 string into a string.
func DecodeBase64(s string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

// CalculateWorkingDays calulates the number of working days between two dates
func CalculateWorkingDays(startTime time.Time, endTime time.Time) int {
	// Reduce dates to previous Mondays
	startOffset := weekday(startTime)
	startTime = startTime.AddDate(0, 0, -startOffset)
	endOffset := weekday(endTime)
	endTime = endTime.AddDate(0, 0, -endOffset)

	// Calculate weeks and days
	dif := endTime.Sub(startTime)
	weeks := int(math.Round((dif.Hours() / 24) / 7))
	days := -min(startOffset, 5) + min(endOffset, 5)

	// Calculate total days
	return weeks*5 + days
}

func weekday(d time.Time) int {
	wd := d.Weekday()
	if wd == time.Sunday {
		return 6
	}
	return int(wd) - 1
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
