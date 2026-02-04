package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dedene/irail-cli/internal/api"
	"github.com/dedene/irail-cli/internal/output"
)

type DisturbancesCmd struct {
	Type string `help:"Filter by type: planned, disturbance" short:"t"`
}

func (c *DisturbancesCmd) Run(root *RootFlags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := api.NewClient(root.DetectLang())

	resp, err := client.GetDisturbances(ctx)
	if err != nil {
		return wrapAPIError(err)
	}

	disturbances := resp.Disturbance

	// Filter by type if specified
	if c.Type != "" {
		disturbances = filterDisturbances(disturbances, c.Type)
	}

	if root.JSON {
		return output.WriteJSON(os.Stdout, disturbances)
	}

	if len(disturbances) == 0 {
		fmt.Println("No disturbances found")

		return nil
	}

	colors := output.NewColors(os.Stdout, root.NoColor)

	for _, d := range disturbances {
		printDisturbance(d, colors)
	}

	return nil
}

func filterDisturbances(disturbances []api.Disturbance, typeFilter string) []api.Disturbance {
	var result []api.Disturbance

	for _, d := range disturbances {
		if strings.EqualFold(d.Type, typeFilter) {
			result = append(result, d)
		}
	}

	return result
}

func printDisturbance(d api.Disturbance, colors *output.Colors) {
	// Type badge
	typeBadge := formatDisturbanceType(d.Type, colors)

	fmt.Printf("%s %s\n", typeBadge, colors.Bold(d.Title))

	if d.Description != "" {
		// Wrap long descriptions
		desc := strings.TrimSpace(d.Description)
		fmt.Printf("  %s\n", desc)
	}

	if d.Link != "" {
		fmt.Printf("  %s %s\n", colors.Dim("Link:"), d.Link)
	}

	fmt.Println()
}

func formatDisturbanceType(t string, colors *output.Colors) string {
	switch strings.ToLower(t) {
	case "planned":
		return colors.Blue("[PLANNED]")
	case "disturbance":
		return colors.Red("[DISRUPTION]")
	default:
		return colors.Yellow("[" + strings.ToUpper(t) + "]")
	}
}
