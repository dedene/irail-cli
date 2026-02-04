package output

import (
	"fmt"
	"io"
	"os"

	"github.com/muesli/termenv"
)

// Colors provides color formatting for terminal output.
type Colors struct {
	profile termenv.Profile
	output  *termenv.Output
	enabled bool
}

// NewColors creates a new Colors instance.
func NewColors(w io.Writer, noColor bool) *Colors {
	output := termenv.NewOutput(w, termenv.WithProfile(termenv.EnvColorProfile()))

	enabled := !noColor && !termenv.EnvNoColor() && output.Profile != termenv.Ascii

	return &Colors{
		profile: output.Profile,
		output:  output,
		enabled: enabled,
	}
}

// DefaultColors creates colors for stdout with auto-detection.
func DefaultColors() *Colors {
	return NewColors(os.Stdout, false)
}

// IsEnabled returns whether colors are enabled.
func (c *Colors) IsEnabled() bool {
	return c.enabled
}

// Red returns text in red.
func (c *Colors) Red(s string) string {
	if !c.enabled {
		return s
	}

	return termenv.String(s).Foreground(c.profile.Color("#ef4444")).String()
}

// BoldRed returns text in bold red.
func (c *Colors) BoldRed(s string) string {
	if !c.enabled {
		return s
	}

	return termenv.String(s).Foreground(c.profile.Color("#ef4444")).Bold().String()
}

// Green returns text in green.
func (c *Colors) Green(s string) string {
	if !c.enabled {
		return s
	}

	return termenv.String(s).Foreground(c.profile.Color("#22c55e")).String()
}

// Yellow returns text in yellow.
func (c *Colors) Yellow(s string) string {
	if !c.enabled {
		return s
	}

	return termenv.String(s).Foreground(c.profile.Color("#eab308")).String()
}

// Blue returns text in blue.
func (c *Colors) Blue(s string) string {
	if !c.enabled {
		return s
	}

	return termenv.String(s).Foreground(c.profile.Color("#3b82f6")).String()
}

// Dim returns text in a dimmed color.
func (c *Colors) Dim(s string) string {
	if !c.enabled {
		return s
	}

	return termenv.String(s).Foreground(c.profile.Color("#9ca3af")).String()
}

// Bold returns text in bold.
func (c *Colors) Bold(s string) string {
	if !c.enabled {
		return s
	}

	return termenv.String(s).Bold().String()
}

// FormatDelay formats a delay value with appropriate color.
// Positive delays are red, negative (early) are green, zero is unchanged.
func (c *Colors) FormatDelay(delaySeconds int) string {
	if delaySeconds == 0 {
		return ""
	}

	delayMins := delaySeconds / 60

	// Small delays (< 1 min) are negligible
	if delayMins == 0 {
		return ""
	}

	if delayMins > 0 {
		return c.Red(fmt.Sprintf("+%dm", delayMins))
	}

	return c.Green(fmt.Sprintf("%dm", delayMins))
}

// FormatPlatformChange formats a platform change warning.
func (c *Colors) FormatPlatformChange(platform string, isChanged bool) string {
	if !isChanged {
		return platform
	}

	return c.Yellow(fmt.Sprintf("%s ⚠", platform))
}

// FormatCanceled formats a canceled indicator.
func (c *Colors) FormatCanceled() string {
	return c.BoldRed("CANCELED")
}

// FormatOccupancy formats an occupancy indicator.
func (c *Colors) FormatOccupancy(occupancy string) string {
	switch occupancy {
	case "low":
		return c.Green("○○○")
	case "medium":
		return c.Yellow("●○○")
	case "high":
		return c.Red("●●○")
	default:
		return c.Dim("···")
	}
}
