package util

import (
	"os"
	"time"
)

// Getenv will return an environment variable if exists or default if not
func Getenv(name, def string) string {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	return v
}

// GetDayFromUTCTime is a helper function to get the day from an input utc timestamp
func GetDayFromUTCTime(utcTime time.Time) time.Time {
	return utcTime.Truncate(time.Hour * 24)
}
