package output

import (
	"fmt"
	"strconv"
	"time"
)

// ParseUnixTimestamp converts a Unix timestamp string to time.Time.
func ParseUnixTimestamp(ts string) time.Time {
	if ts == "" {
		return time.Time{}
	}

	sec, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(sec, 0)
}

// FormatTime formats a time as HH:MM.
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("15:04")
}

// FormatTimeFromTimestamp converts a Unix timestamp to HH:MM.
func FormatTimeFromTimestamp(ts string) string {
	return FormatTime(ParseUnixTimestamp(ts))
}

// FormatRelativeTime formats a time relative to now (e.g., "in 5m", "2m ago").
func FormatRelativeTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	now := time.Now()
	diff := t.Sub(now)

	if diff < 0 {
		diff = -diff
		mins := int(diff.Minutes())

		if mins < 1 {
			return "now"
		}

		if mins < 60 {
			return fmt.Sprintf("%dm ago", mins)
		}

		hours := mins / 60
		mins %= 60

		if mins == 0 {
			return fmt.Sprintf("%dh ago", hours)
		}

		return fmt.Sprintf("%dh%dm ago", hours, mins)
	}

	mins := int(diff.Minutes())

	if mins < 1 {
		return "now"
	}

	if mins < 60 {
		return fmt.Sprintf("in %dm", mins)
	}

	hours := mins / 60
	mins %= 60

	if mins == 0 {
		return fmt.Sprintf("in %dh", hours)
	}

	return fmt.Sprintf("in %dh%dm", hours, mins)
}

// FormatRelativeFromTimestamp converts a Unix timestamp to relative time.
func FormatRelativeFromTimestamp(ts string) string {
	return FormatRelativeTime(ParseUnixTimestamp(ts))
}

// FormatDuration formats a duration in seconds as "1h27m".
func FormatDuration(seconds int) string {
	if seconds <= 0 {
		return "0m"
	}

	hours := seconds / 3600
	mins := (seconds % 3600) / 60

	if hours == 0 {
		return fmt.Sprintf("%dm", mins)
	}

	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}

	return fmt.Sprintf("%dh%dm", hours, mins)
}

// FormatDurationFromString parses a duration string (seconds) and formats it.
func FormatDurationFromString(s string) string {
	secs, err := strconv.Atoi(s)
	if err != nil {
		return s
	}

	return FormatDuration(secs)
}

// ParseDelay parses a delay string (seconds) and returns the integer value.
func ParseDelay(s string) int {
	if s == "" {
		return 0
	}

	delay, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return delay
}

// FormatDate formats a time as a date string for API requests (DDMMYY).
func FormatDateForAPI(t time.Time) string {
	return t.Format("020106")
}

// FormatTimeForAPI formats a time for API requests (HHMM).
func FormatTimeForAPI(t time.Time) string {
	return t.Format("1504")
}

// ParseDate parses a date string in YYYY-MM-DD format.
func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// ParseTimeString parses a time string in HH:MM format.
func ParseTimeString(s string) (time.Time, error) {
	return time.Parse("15:04", s)
}

// ConvertDateForAPI converts YYYY-MM-DD to DDMMYY format.
func ConvertDateForAPI(s string) string {
	if s == "" {
		return ""
	}

	t, err := ParseDate(s)
	if err != nil {
		return ""
	}

	return FormatDateForAPI(t)
}

// ConvertTimeForAPI converts HH:MM to HHMM format.
func ConvertTimeForAPI(s string) string {
	if s == "" {
		return ""
	}

	t, err := ParseTimeString(s)
	if err != nil {
		return ""
	}

	return FormatTimeForAPI(t)
}
