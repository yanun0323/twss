package util

import "time"

/* 2006-01-02 */
func LogDate(date time.Time) string {
	return date.Format("2006-01-02")
}

/* 20060102 */
func FormatDate(date time.Time) string {
	return date.Format("20060102")
}

func NextDate(date time.Time) time.Time {
	return date.Add(24 * time.Hour)
}
