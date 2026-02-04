package cmd

import (
	"context"
	"errors"
	"strings"

	"github.com/dedene/irail-cli/internal/api"
	"github.com/dedene/irail-cli/internal/tui"
)

// ErrStationCanceled is returned when station selection is canceled.
var ErrStationCanceled = errors.New("station selection canceled")

// ResolveStation resolves a station name to a Station object.
// If the name is ambiguous, it shows a picker.
// If JSON mode is enabled, it returns an error instead of showing picker.
func ResolveStation(ctx context.Context, client *api.Client, name string, jsonMode bool) (*api.Station, error) {
	// Fetch all stations
	resp, err := client.GetStations(ctx)
	if err != nil {
		return nil, err
	}

	// Try exact match first (case-insensitive)
	nameLower := strings.ToLower(name)
	for _, s := range resp.Station {
		if strings.ToLower(s.Name) == nameLower || strings.ToLower(s.StandardName) == nameLower {
			return &s, nil
		}
	}

	// Try ID match
	for _, s := range resp.Station {
		if s.ID == name {
			return &s, nil
		}
	}

	// Try partial match
	var matches []api.Station

	for _, s := range resp.Station {
		if strings.Contains(strings.ToLower(s.Name), nameLower) ||
			strings.Contains(strings.ToLower(s.StandardName), nameLower) {
			matches = append(matches, s)
		}
	}

	// No matches
	if len(matches) == 0 {
		return nil, &api.NotFoundError{Resource: "station", ID: name}
	}

	// Single match
	if len(matches) == 1 {
		return &matches[0], nil
	}

	// Multiple matches - show picker if not in JSON mode
	if jsonMode {
		return nil, &api.NotFoundError{Resource: "station", ID: name}
	}

	result, err := tui.PickStation(matches)
	if err != nil {
		return nil, err
	}

	if result.Canceled {
		return nil, ErrStationCanceled
	}

	return &result.Station, nil
}
