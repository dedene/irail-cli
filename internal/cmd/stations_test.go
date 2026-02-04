package cmd

import (
	"testing"

	"github.com/dedene/irail-cli/internal/api"
)

func TestFilterStations(t *testing.T) {
	stations := []api.Station{
		{Name: "Brussel-Centraal", StandardName: "Brussel-Centraal/Bruxelles-Central", ID: "BE.NMBS.008813003"},
		{Name: "Brussel-Noord", StandardName: "Brussel-Noord/Bruxelles-Nord", ID: "BE.NMBS.008812005"},
		{Name: "Brugge", StandardName: "Brugge", ID: "BE.NMBS.008891009"},
		{Name: "Gent-Sint-Pieters", StandardName: "Gent-Sint-Pieters", ID: "BE.NMBS.008892007"},
	}

	tests := []struct {
		search    string
		wantCount int
	}{
		{"brussel", 2},
		{"brugge", 1},
		{"gent", 1},
		{"xyz", 0},
		{"", 4},        // Empty search returns all
		{"BRUSSEL", 2}, // Case insensitive
		{"bru", 3},     // Partial match
	}

	for _, tt := range tests {
		got := filterStations(stations, tt.search)
		if len(got) != tt.wantCount {
			t.Errorf("filterStations(%q) returned %d stations, want %d", tt.search, len(got), tt.wantCount)
		}
	}
}
