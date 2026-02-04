package errfmt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dedene/irail-cli/internal/api"
)

// Format formats an error into a user-friendly message.
func Format(err error) string {
	if err == nil {
		return ""
	}

	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		return formatAPIError(apiErr)
	}

	var notFoundErr *api.NotFoundError
	if errors.As(err, &notFoundErr) {
		return formatNotFoundError(notFoundErr)
	}

	var rateLimitErr *api.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return formatRateLimitError(rateLimitErr)
	}

	var circuitBreakerErr *api.CircuitBreakerError
	if errors.As(err, &circuitBreakerErr) {
		return formatCircuitBreakerError()
	}

	return fmt.Sprintf("Error: %v", err)
}

func formatAPIError(err *api.APIError) string {
	var sb strings.Builder

	switch err.StatusCode {
	case 400:
		sb.WriteString("Error: Invalid request (400)\n\n")
		sb.WriteString("  Check your input parameters.\n")

	case 404:
		sb.WriteString("Error: Not found (404)\n\n")

		if err.Details != "" {
			sb.WriteString("  " + err.Details + "\n\n")
		}

		sb.WriteString("  The station or resource doesn't exist.\n")
		sb.WriteString("  Tip: Use 'irail stations --search <name>' to find the correct station name.\n")

	case 429:
		sb.WriteString("Error: Rate limit exceeded (429)\n\n")
		sb.WriteString("  Too many requests. Please wait a moment and try again.\n")

	case 500, 502, 503, 504:
		sb.WriteString(fmt.Sprintf("Error: Service unavailable (%d)\n\n", err.StatusCode))
		sb.WriteString("  The iRail API is temporarily unavailable.\n")
		sb.WriteString("  Please try again in a few moments.\n")

	default:
		sb.WriteString(fmt.Sprintf("Error: %s (%d)\n", err.Message, err.StatusCode))

		if err.Details != "" {
			sb.WriteString("\n  " + err.Details + "\n")
		}
	}

	return sb.String()
}

func formatNotFoundError(err *api.NotFoundError) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Error: %s not found", err.Resource))

	if err.ID != "" {
		sb.WriteString(fmt.Sprintf(": %s", err.ID))
	}

	sb.WriteString("\n\n")

	switch err.Resource {
	case "station":
		sb.WriteString("  The station name was not recognized.\n")
		sb.WriteString("  Tip: Use 'irail stations --search <name>' to find stations.\n")
	case "vehicle":
		sb.WriteString("  The vehicle/train ID was not found.\n")
		sb.WriteString("  Make sure the train number is correct and the train is running today.\n")
	case "composition":
		sb.WriteString("  Train composition data is not available for this train.\n")
		sb.WriteString("  Composition data may not be available for all train types.\n")
	default:
		sb.WriteString("  The requested resource could not be found.\n")
	}

	return sb.String()
}

func formatRateLimitError(err *api.RateLimitError) string {
	var sb strings.Builder

	sb.WriteString("Error: Rate limit exceeded\n\n")

	if err.RetryAfter > 0 {
		sb.WriteString(fmt.Sprintf("  Please wait %d seconds before trying again.\n", err.RetryAfter))
	} else {
		sb.WriteString("  Please wait a moment and try again.\n")
	}

	return sb.String()
}

func formatCircuitBreakerError() string {
	var sb strings.Builder

	sb.WriteString("Error: Service temporarily unavailable\n\n")
	sb.WriteString("  The iRail API has been experiencing issues.\n")
	sb.WriteString("  Please wait a moment and try again.\n")

	return sb.String()
}
