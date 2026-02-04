package output

import (
	"io"
	"os"
	"strings"

	"github.com/muesli/termenv"
)

// Hyperlinks provides terminal hyperlink support.
type Hyperlinks struct {
	output  *termenv.Output
	enabled bool
}

// NewHyperlinks creates a new Hyperlinks instance.
func NewHyperlinks(w io.Writer) *Hyperlinks {
	output := termenv.NewOutput(w, termenv.WithProfile(termenv.EnvColorProfile()))

	// Enable hyperlinks if terminal supports colors (reasonable heuristic)
	// Most modern terminals that support colors also support OSC 8 hyperlinks
	enabled := !termenv.EnvNoColor() && output.Profile != termenv.Ascii

	return &Hyperlinks{
		output:  output,
		enabled: enabled,
	}
}

// DefaultHyperlinks creates hyperlinks for stdout.
func DefaultHyperlinks() *Hyperlinks {
	return NewHyperlinks(os.Stdout)
}

// IsEnabled returns whether hyperlinks are supported.
func (h *Hyperlinks) IsEnabled() bool {
	return h.enabled
}

// Link creates a clickable hyperlink if supported.
// Falls back to just the text if not supported.
func (h *Hyperlinks) Link(text, url string) string {
	if !h.enabled || url == "" {
		return text
	}

	return h.output.Hyperlink(url, text)
}

// TrainLink formats the train name for display.
// Note: belgiantrain.be doesn't support direct train URLs, so we just return the short name.
func (h *Hyperlinks) TrainLink(trainName string) string {
	if trainName == "" {
		return ""
	}

	// Extract train number from name like "BE.NMBS.IC1832" or "IC1832"
	return extractTrainNumber(trainName)
}

// extractTrainNumber extracts the train number from various formats.
func extractTrainNumber(vehicle string) string {
	// Handle BE.NMBS.IC1832 format
	if short, found := strings.CutPrefix(vehicle, "BE.NMBS."); found {
		return short
	}

	// Handle http://irail.be/vehicle/IC1832 format
	if _, after, found := strings.Cut(vehicle, "/vehicle/"); found {
		return after
	}

	return vehicle
}
