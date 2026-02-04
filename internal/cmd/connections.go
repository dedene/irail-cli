package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dedene/irail-cli/internal/api"
	"github.com/dedene/irail-cli/internal/output"
)

type ConnectionsCmd struct {
	From     string `arg:"" help:"Departure station"`
	To       string `arg:"" help:"Arrival station"`
	Time     string `help:"Departure time (HH:MM)" short:"t"`
	Date     string `help:"Date (YYYY-MM-DD)" short:"d"`
	ArriveBy bool   `help:"Time is arrival time" short:"a"`
	Results  int    `help:"Number of results" default:"6" short:"n"`
}

func (c *ConnectionsCmd) Run(root *RootFlags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := api.NewClient(root.DetectLang())

	apiDate := output.ConvertDateForAPI(c.Date)
	apiTime := output.ConvertTimeForAPI(c.Time)

	resp, err := client.GetConnections(ctx, c.From, c.To, apiDate, apiTime, c.ArriveBy, c.Results)
	if err != nil {
		return wrapAPIError(err)
	}

	if root.JSON {
		return output.WriteJSON(os.Stdout, resp)
	}

	colors := output.NewColors(os.Stdout, root.NoColor)
	hyperlinks := output.DefaultHyperlinks()

	for i, conn := range resp.Connection {
		printConnection(i+1, conn, colors, hyperlinks)
	}

	return nil
}

func printConnection(index int, conn api.Connection, colors *output.Colors, hyperlinks *output.Hyperlinks) {
	depTime := output.FormatTimeFromTimestamp(conn.Departure.Time)
	arrTime := output.FormatTimeFromTimestamp(conn.Arrival.Time)

	duration, _ := strconv.Atoi(conn.Duration)
	durationStr := output.FormatDuration(duration)

	// Count transfers
	numVias := 0
	if conn.Vias.Number != "" {
		numVias, _ = strconv.Atoi(conn.Vias.Number)
	}

	// Build train list
	trains := []string{hyperlinks.TrainLink(conn.Departure.Vehicle)}
	for _, via := range conn.Vias.Via {
		trains = append(trains, hyperlinks.TrainLink(via.Departure.Vehicle))
	}
	trainsStr := strings.Join(trains, " → ")

	// Transfer count text
	var transferStr string
	switch numVias {
	case 0:
		transferStr = colors.Green("direct")
	case 1:
		transferStr = "1 transfer"
	default:
		transferStr = fmt.Sprintf("%d transfers", numVias)
	}

	// Header line
	fmt.Printf("[%s] %s → %s (%s) | %s | %s\n",
		colors.Bold(fmt.Sprintf("%d", index)),
		colors.Bold(depTime),
		colors.Bold(arrTime),
		colors.Dim(durationStr),
		trainsStr,
		transferStr,
	)

	// Departure info
	depDelay := output.ParseDelay(conn.Departure.Delay)
	depDelayStr := colors.FormatDelay(depDelay)
	depPlatform := conn.Departure.Platform
	depPlatformChanged := conn.Departure.PlatformInfo.Normal == "0"

	fmt.Printf("    %s %s P%s %s\n",
		colors.Dim("From:"),
		conn.Departure.Station,
		colors.FormatPlatformChange(depPlatform, depPlatformChanged),
		depDelayStr,
	)

	// Via info (transfers)
	for _, via := range conn.Vias.Via {
		waitTime, _ := strconv.Atoi(via.TimeBetween)
		waitStr := output.FormatDuration(waitTime)

		arrPlatform := via.Arrival.Platform
		depPlatform := via.Departure.Platform

		// Check for walking transfer
		if via.Arrival.Walking == "1" || via.Departure.Walking == "1" {
			fmt.Printf("    %s %s (%s wait, 🚶 walking) P%s → P%s\n",
				colors.Yellow("Via:"),
				via.Station,
				waitStr,
				arrPlatform,
				depPlatform,
			)
		} else {
			fmt.Printf("    %s %s (%s wait) P%s → P%s\n",
				colors.Yellow("Via:"),
				via.Station,
				waitStr,
				arrPlatform,
				depPlatform,
			)
		}
	}

	// Arrival info
	arrDelay := output.ParseDelay(conn.Arrival.Delay)
	arrDelayStr := colors.FormatDelay(arrDelay)
	arrPlatform := conn.Arrival.Platform
	arrPlatformChanged := conn.Arrival.PlatformInfo.Normal == "0"

	fmt.Printf("    %s %s P%s %s\n",
		colors.Dim("To:"),
		conn.Arrival.Station,
		colors.FormatPlatformChange(arrPlatform, arrPlatformChanged),
		arrDelayStr,
	)

	fmt.Println()
}
