package db

import (
	"fmt"
)

var (
	ErrEtagMismatch = fmt.Errorf("etag mismatch")
	ErrDuplicate    = fmt.Errorf("duplicate object")
	ErrNotFound     = fmt.Errorf("not found")
	ErrInvalid      = fmt.Errorf("object invalid")
)
