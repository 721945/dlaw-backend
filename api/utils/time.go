package utils

import "time"

func CovertStringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func CovertTimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}
