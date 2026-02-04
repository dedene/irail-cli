package cmd

import (
	"errors"

	"github.com/dedene/irail-cli/internal/api"
)

// Exit codes
const (
	ExitSuccess      = 0
	ExitGeneralError = 1
	ExitInvalidArgs  = 2
	ExitAPIError     = 3
	ExitNotFound     = 4
	ExitRateLimited  = 5
)

type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}

	return e.Err.Error()
}

func (e *ExitError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func ExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}

	var ee *ExitError
	if errors.As(err, &ee) && ee != nil {
		if ee.Code < 0 {
			return ExitGeneralError
		}

		return ee.Code
	}

	return ExitGeneralError
}

// wrapAPIError wraps an API error with the appropriate exit code.
func wrapAPIError(err error) *ExitError {
	if err == nil {
		return nil
	}

	var notFound *api.NotFoundError
	if errors.As(err, &notFound) {
		return &ExitError{Code: ExitNotFound, Err: err}
	}

	var rateLimit *api.RateLimitError
	if errors.As(err, &rateLimit) {
		return &ExitError{Code: ExitRateLimited, Err: err}
	}

	return &ExitError{Code: ExitAPIError, Err: err}
}
