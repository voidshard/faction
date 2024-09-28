package db

import (
	"fmt"
)

var (
	ErrEtagMismatch = fmt.Errorf("etag mismatch")
	ErrNotFound     = fmt.Errorf("not found")
)
