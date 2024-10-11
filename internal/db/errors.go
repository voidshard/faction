package db

import (
	"errors"
	"fmt"
)

var (
	ErrEtagMismatch = fmt.Errorf("etag mismatch")
	ErrNotFound     = fmt.Errorf("not found")
	ErrInvalid      = fmt.Errorf("object invalid")
)

func ErrorRetryable(err error) bool {
	if errors.Is(err, ErrEtagMismatch) {
		return false
	} else if errors.Is(err, ErrNotFound) {
		return false
	} else if errors.Is(err, ErrInvalid) {
		return false
	}
	return true
}
