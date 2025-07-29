package utils

import "time"

func PtrString(s string) *string {
	return &s
}

func Deref[T any](s *T) T {
	if s == nil {
		var zero T
		return zero
	}
	return *s
}

func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
