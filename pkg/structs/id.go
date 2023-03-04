package structs

import (
	"github.com/voidshard/faction/internal/dbutils"
)

// NewID returns a 36 char UUID that is
// - random, if no args are passed in
// - deterministic, if args are passed in
// For a given struct type, you should stick to one or the other.
//
// Deterministic is recommended as it allows you to translate your IDs
// (whatever they are) to UUIDs valid with this library whenever you
// want.
func NewID(args ...interface{}) string {
	if len(args) == 0 {
		return dbutils.RandomID()
	}
	return dbutils.NewID(args...)
}

// IsValidID returns true if the given ID is a valid UUID
func IsValidID(id string) bool {
	return dbutils.IsValidID(id)
}
