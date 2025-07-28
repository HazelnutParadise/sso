package utils

import "time"

func PtrString(s string) *string {
	return &s
}

func Deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
