package cmd

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/dedene/irail-cli/internal/api"
	"github.com/dedene/irail-cli/internal/errfmt"
	"github.com/dedene/irail-cli/internal/output"
)

type StationsCmd struct {
	Search string `help:"Filter stations by name" short:"s"`
}

func (c *StationsCmd) Run(root *RootFlags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := api.NewClient(root.DetectLang())

	resp, err := client.GetStations(ctx)
	if err != nil {
		return wrapAPIError(err)
	}

	stations := resp.Station

	// Filter by search term if provided
	if c.Search != "" {
		stations = filterStations(stations, c.Search)
	}

	if root.JSON {
		return output.WriteJSON(os.Stdout, stations)
	}

	// Table output
	colors := output.NewColors(os.Stdout, root.NoColor)
	table := output.NewTable(os.Stdout)
	table.SetHeaders(
		colors.Bold("NAME"),
		colors.Bold("STANDARD NAME"),
		colors.Bold("ID"),
	)
	table.WriteHeader()

	for _, s := range stations {
		table.AddRow(s.Name, s.StandardName, s.ID)
	}

	if err := table.Render(); err != nil {
		return &ExitError{Code: ExitGeneralError, Err: err}
	}

	if len(stations) == 0 && c.Search != "" {
		_, _ = os.Stderr.WriteString(errfmt.Format(&api.NotFoundError{Resource: "station", ID: c.Search}))
	}

	return nil
}

// filterStations performs case-insensitive fuzzy filtering of stations.
func filterStations(stations []api.Station, search string) []api.Station {
	search = strings.ToLower(search)
	var result []api.Station

	for _, s := range stations {
		name := strings.ToLower(s.Name)
		stdName := strings.ToLower(s.StandardName)

		if strings.Contains(name, search) || strings.Contains(stdName, search) {
			result = append(result, s)
		}
	}

	return result
}
