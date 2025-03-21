package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/voidshard/faction/internal/db"
)

var (
	ErrShuttingDown = fmt.Errorf("shutting down")
	ErrInvalid      = fmt.Errorf("object invalid")
	ErrPrecondition = fmt.Errorf("precondition failed")
	ErrNotFound     = fmt.Errorf("object not found")
)

// errorCodeHTTP returns the HTTP status code for a given error.
// If we cannot find a specific match we'll default to 500.
func errorCodeHTTP(err error) int {
	// API errors
	if errors.Is(err, ErrShuttingDown) {
		return http.StatusServiceUnavailable
	} else if errors.Is(err, ErrInvalid) {
		return http.StatusBadRequest
	} else if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	} else if errors.Is(err, ErrPrecondition) {
		return http.StatusPreconditionFailed
	}
	// DB errors
	if errors.Is(err, db.ErrNotFound) {
		return http.StatusNotFound
	} else if errors.Is(err, db.ErrDuplicate) {
		return http.StatusConflict
	} else if errors.Is(err, db.ErrInvalid) {
		return http.StatusBadRequest
	} else if errors.Is(err, db.ErrEtagMismatch) {
		return http.StatusPreconditionFailed
	}
	return http.StatusInternalServerError
}
