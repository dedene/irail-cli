package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dedene/irail-cli/internal/api"
	"github.com/dedene/irail-cli/internal/output"
)

type LiveboardCmd struct {
	Station  string `arg:"" help:"Station name or ID"`
	Arrivals bool   `help:"Show arrivals instead of departures" short:"a"`
	Time     string `help:"Time (HH:MM)" short:"t"`
	Date     string `help:"Date (YYYY-MM-DD)" short:"d"`
	Alerts   bool   `help:"Show alerts" default:"true"`
}

func (c *LiveboardCmd) Run(root *RootFlags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := api.NewClient(root.DetectLang())

	// Convert date/time formats for API
	apiDate := output.ConvertDateForAPI(c.Date)
	apiTime := output.ConvertTimeForAPI(c.Time)

	resp, err := client.GetLiveboard(ctx, c.Station, c.Arrivals, apiDate, apiTime)
	if err != nil {
		return wrapAPIError(err)
	}

	if root.JSON {
		return output.WriteJSON(os.Stdout, resp)
	}

	colors := output.NewColors(os.Stdout, root.NoColor)
	hyperlinks := output.DefaultHyperlinks()

	// Print station header
	fmt.Printf("%s - %s\n\n", colors.Bold(resp.Station), modeLabel(c.Arrivals))

	if c.Arrivals {
		return c.printArrivals(resp, colors, hyperlinks)
	}

	return c.printDepartures(resp, colors, hyperlinks)
}

func modeLabel(arrivals bool) string {
	if arrivals {
		return "Arrivals"
	}

	return "Departures"
}

func (c *LiveboardCmd) printDepartures(resp *api.LiveboardResponse, colors *output.Colors, hyperlinks *output.Hyperlinks) error {
	table := output.NewTable(os.Stdout)
	table.SetHeaders(
		colors.Bold("TIME"),
		colors.Bold(""),
		colors.Bold("DELAY"),
		colors.Bold("TRAIN"),
		colors.Bold("PLATFORM"),
		colors.Bold("DESTINATION"),
		colors.Bold("OCCUPANCY"),
	)
	table.WriteHeader()

	for _, dep := range resp.Departures.Departure {
		depTime := output.FormatTimeFromTimestamp(dep.Time)
		relative := output.FormatRelativeFromTimestamp(dep.Time)
		delay := output.ParseDelay(dep.Delay)
		delayStr := colors.FormatDelay(delay)

		trainName := hyperlinks.TrainLink(dep.Vehicle)

		platform := dep.Platform
		platformChanged := dep.PlatformInfo.Normal == "0"
		platformStr := colors.FormatPlatformChange(platform, platformChanged)

		destination := dep.Station

		occupancy := extractOccupancy(dep.Occupancy.Name)
		occupancyStr := colors.FormatOccupancy(occupancy)

		// Handle canceled trains
		if dep.Canceled == "1" {
			table.AddRow(
				depTime,
				colors.Dim(relative),
				colors.FormatCanceled(),
				trainName,
				platformStr,
				destination,
				"",
			)

			continue
		}

		table.AddRow(
			depTime,
			colors.Dim(relative),
			delayStr,
			trainName,
			platformStr,
			destination,
			occupancyStr,
		)
	}

	return table.Render()
}

func (c *LiveboardCmd) printArrivals(resp *api.LiveboardResponse, colors *output.Colors, hyperlinks *output.Hyperlinks) error {
	table := output.NewTable(os.Stdout)
	table.SetHeaders(
		colors.Bold("TIME"),
		colors.Bold(""),
		colors.Bold("DELAY"),
		colors.Bold("TRAIN"),
		colors.Bold("PLATFORM"),
		colors.Bold("ORIGIN"),
	)
	table.WriteHeader()

	for _, arr := range resp.Arrivals.Arrival {
		arrTime := output.FormatTimeFromTimestamp(arr.Time)
		relative := output.FormatRelativeFromTimestamp(arr.Time)
		delay := output.ParseDelay(arr.Delay)
		delayStr := colors.FormatDelay(delay)

		trainName := hyperlinks.TrainLink(arr.Vehicle)

		platform := arr.Platform
		platformChanged := arr.PlatformInfo.Normal == "0"
		platformStr := colors.FormatPlatformChange(platform, platformChanged)

		origin := arr.Station

		// Handle canceled trains
		if arr.Canceled == "1" {
			table.AddRow(
				arrTime,
				colors.Dim(relative),
				colors.FormatCanceled(),
				trainName,
				platformStr,
				origin,
			)

			continue
		}

		table.AddRow(
			arrTime,
			colors.Dim(relative),
			delayStr,
			trainName,
			platformStr,
			origin,
		)
	}

	return table.Render()
}

func extractOccupancy(name string) string {
	switch name {
	case "low":
		return "low"
	case "medium":
		return "medium"
	case "high":
		return "high"
	default:
		return ""
	}
}
