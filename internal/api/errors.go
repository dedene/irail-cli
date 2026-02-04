package api

import "fmt"

// APIError represents an error from the iRail API.
type APIError struct {
	StatusCode int
	Message    string
	Details    string
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API error %d: %s - %s", e.StatusCode, e.Message, e.Details)
	}

	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// NotFoundError indicates a resource was not found.
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s not found: %s", e.Resource, e.ID)
	}

	return fmt.Sprintf("%s not found", e.Resource)
}

// RateLimitError indicates the API rate limit was exceeded.
type RateLimitError struct {
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("rate limit exceeded, retry after %d seconds", e.RetryAfter)
	}

	return "rate limit exceeded"
}

// CircuitBreakerError indicates the circuit breaker is open.
type CircuitBreakerError struct{}

func (e *CircuitBreakerError) Error() string {
	return "service temporarily unavailable (circuit breaker open)"
}
