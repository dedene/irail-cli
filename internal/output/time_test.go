package output

import (
	"testing"
	"time"
)

func TestParseUnixTimestamp(t *testing.T) {
	tests := []struct {
		input string
		want  time.Time
	}{
		{"1609459200", time.Unix(1609459200, 0)},
		{"0", time.Unix(0, 0)},
		{"", time.Time{}},
		{"invalid", time.Time{}},
	}

	for _, tt := range tests {
		got := ParseUnixTimestamp(tt.input)
		if !got.Equal(tt.want) {
			t.Errorf("ParseUnixTimestamp(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		input time.Time
		want  string
	}{
		{time.Date(2025, 1, 1, 9, 5, 0, 0, time.UTC), "09:05"},
		{time.Date(2025, 1, 1, 14, 30, 0, 0, time.UTC), "14:30"},
		{time.Time{}, ""},
	}

	for _, tt := range tests {
		got := FormatTime(tt.input)
		if got != tt.want {
			t.Errorf("FormatTime(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		seconds int
		want    string
	}{
		{0, "0m"},
		{60, "1m"},
		{3600, "1h"},
		{3660, "1h1m"},
		{5400, "1h30m"},
		{7200, "2h"},
		{-1, "0m"},
	}

	for _, tt := range tests {
		got := FormatDuration(tt.seconds)
		if got != tt.want {
			t.Errorf("FormatDuration(%d) = %q, want %q", tt.seconds, got, tt.want)
		}
	}
}

func TestParseDelay(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"0", 0},
		{"60", 60},
		{"300", 300},
		{"", 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		got := ParseDelay(tt.input)
		if got != tt.want {
			t.Errorf("ParseDelay(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestConvertDateForAPI(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"2025-01-15", "150125"},
		{"2025-12-31", "311225"},
		{"", ""},
		{"invalid", ""},
	}

	for _, tt := range tests {
		got := ConvertDateForAPI(tt.input)
		if got != tt.want {
			t.Errorf("ConvertDateForAPI(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestConvertTimeForAPI(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"09:05", "0905"},
		{"14:30", "1430"},
		{"", ""},
		{"invalid", ""},
	}

	for _, tt := range tests {
		got := ConvertTimeForAPI(tt.input)
		if got != tt.want {
			t.Errorf("ConvertTimeForAPI(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
