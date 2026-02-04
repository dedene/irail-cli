package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dedene/irail-cli/internal/api"
	"github.com/dedene/irail-cli/internal/output"
)

type CompositionCmd struct {
	ID string `arg:"" help:"Vehicle/train ID (e.g., S51507)"`
}

func (c *CompositionCmd) Run(root *RootFlags) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := api.NewClient(root.DetectLang())

	// Normalize ID
	vehicleID := c.ID
	if !isFullVehicleID(vehicleID) {
		vehicleID = "BE.NMBS." + vehicleID
	}

	resp, err := client.GetComposition(ctx, vehicleID)
	if err != nil {
		var notFound *api.NotFoundError
		if errors.As(err, &notFound) {
			fmt.Fprintf(os.Stderr, "Composition data not available for %s\n", c.ID)
			fmt.Fprintln(os.Stderr, "Note: Composition data may not be available for all train types.")
		}

		return wrapAPIError(err)
	}

	if root.JSON {
		return output.WriteJSON(os.Stdout, resp)
	}

	colors := output.NewColors(os.Stdout, root.NoColor)

	fmt.Printf("%s Composition\n\n", colors.Bold(c.ID))

	// Aggregate stats
	var totalFirstClass, totalSecondClass int
	var amenities []string
	amenitiesMap := make(map[string]bool)

	for _, segment := range resp.Composition.Segments {
		for _, unit := range segment.Composition.Units.Unit {
			first, _ := strconv.Atoi(unit.SeatsFirstClass)
			second, _ := strconv.Atoi(unit.SeatsSecondClass)
			totalFirstClass += first
			totalSecondClass += second

			if unit.HasToilets == "1" && !amenitiesMap["toilet"] {
				amenitiesMap["toilet"] = true
				amenities = append(amenities, "🚻 Toilet")
			}

			if unit.HasAirco == "1" && !amenitiesMap["airco"] {
				amenitiesMap["airco"] = true
				amenities = append(amenities, "❄️ Air conditioning")
			}

			if unit.HasBikeSection == "1" && !amenitiesMap["bike"] {
				amenitiesMap["bike"] = true
				amenities = append(amenities, "🚲 Bike section")
			}

			if unit.HasPrmSection == "1" && !amenitiesMap["wheelchair"] {
				amenitiesMap["wheelchair"] = true
				amenities = append(amenities, "♿ Wheelchair accessible")
			}

			if unit.HasTables == "1" && !amenitiesMap["tables"] {
				amenitiesMap["tables"] = true
				amenities = append(amenities, "🪑 Tables")
			}

			if unit.HasFirstClassOutlets == "1" || unit.HasSecondClassOutlets == "1" {
				if !amenitiesMap["power"] {
					amenitiesMap["power"] = true
					amenities = append(amenities, "🔌 Power outlets")
				}
			}
		}
	}

	// Print summary
	fmt.Printf("%s\n", colors.Bold("Seating:"))
	fmt.Printf("  1st class: %d seats\n", totalFirstClass)
	fmt.Printf("  2nd class: %d seats\n", totalSecondClass)
	fmt.Printf("  Total: %d seats\n\n", totalFirstClass+totalSecondClass)

	if len(amenities) > 0 {
		fmt.Printf("%s\n", colors.Bold("Amenities:"))

		for _, a := range amenities {
			fmt.Printf("  %s\n", a)
		}
	}

	return nil
}
