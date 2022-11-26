package util

import "time"

/* 20060102 */
func FormatToUrlDate(date time.Time) string {
	return date.Format("20060102")
}

func LogDate(date time.Time) string {
	return date.Format("2006-01-02")
}
