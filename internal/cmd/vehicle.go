package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dedene/irail-cli/internal/api"
	"github.com/dedene/irail-cli/internal/output"
)

type VehicleCmd struct {
	ID     string `arg:"" help:"Vehicle/train ID (e.g., IC1832)"`
	Date   string `help:"Date (YYYY-MM-DD)" short:"d"`
	Stops  bool   `help:"Show all stops" short:"s"`
	Alerts bool   `help:"Show alerts" default:"true"`
}

func (c *VehicleCmd) Run(root *RootFlags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := api.NewClient(root.DetectLang())

	apiDate := output.ConvertDateForAPI(c.Date)

	// Normalize ID - add BE.NMBS. prefix if not present
	vehicleID := c.ID
	if !isFullVehicleID(vehicleID) {
		vehicleID = "BE.NMBS." + vehicleID
	}

	resp, err := client.GetVehicle(ctx, vehicleID, apiDate)
	if err != nil {
		return wrapAPIError(err)
	}

	if root.JSON {
		return output.WriteJSON(os.Stdout, resp)
	}

	colors := output.NewColors(os.Stdout, root.NoColor)

	// Print vehicle header
	trainName := resp.VehicleInfo.ShortName
	if trainName == "" {
		trainName = resp.VehicleInfo.Name
	}

	fmt.Printf("%s\n\n", colors.Bold(trainName))

	if c.Stops {
		return c.printAllStops(resp, colors)
	}

	c.printSummary(resp, colors)

	return nil
}

func isFullVehicleID(id string) bool {
	return len(id) > 8 && id[:3] == "BE."
}

func (c *VehicleCmd) printSummary(resp *api.VehicleResponse, colors *output.Colors) {
	stops := resp.Stops.Stop
	if len(stops) == 0 {
		fmt.Println("No stops found")

		return
	}

	first := stops[0]
	last := stops[len(stops)-1]

	firstTime := output.FormatTimeFromTimestamp(first.Time)
	lastTime := output.FormatTimeFromTimestamp(last.Time)

	fmt.Printf("%s %s → %s %s\n",
		firstTime,
		first.Station,
		lastTime,
		last.Station,
	)

	fmt.Printf("\n%s\n", colors.Dim(fmt.Sprintf("%d stops total", len(stops))))
}

func (c *VehicleCmd) printAllStops(resp *api.VehicleResponse, colors *output.Colors) error {
	table := output.NewTable(os.Stdout)
	table.SetHeaders(
		colors.Bold("TIME"),
		colors.Bold("DELAY"),
		colors.Bold("STATION"),
		colors.Bold("PLATFORM"),
		colors.Bold("STATUS"),
	)
	table.WriteHeader()

	for _, stop := range resp.Stops.Stop {
		stopTime := output.FormatTimeFromTimestamp(stop.Time)
		delay := output.ParseDelay(stop.Delay)
		delayStr := colors.FormatDelay(delay)

		platform := stop.Platform
		platformChanged := stop.PlatformInfo.Normal == "0"
		platformStr := colors.FormatPlatformChange(platform, platformChanged)

		status := getStopStatus(stop, colors)

		table.AddRow(
			stopTime,
			delayStr,
			stop.Station,
			platformStr,
			status,
		)
	}

	return table.Render()
}

func getStopStatus(stop api.Stop, colors *output.Colors) string {
	if stop.Canceled == "1" {
		return colors.FormatCanceled()
	}

	if stop.Left == "1" {
		return colors.Green("departed")
	}

	if stop.Arrived == "1" {
		return colors.Blue("arrived")
	}

	if stop.IsExtraStop == "1" {
		return colors.Yellow("extra stop")
	}

	return ""
}
