package kind

import (
	"regexp"

	"github.com/go-playground/validator"
	"github.com/voidshard/faction/pkg/util/uuid"
)

var alphanum = regexp.MustCompile("^[a-zA-Z0-9]*$")
var alphanumsym = regexp.MustCompile("^[a-zA-Z0-9\\-_./]*$")

func isUUID4(s string) bool {
	return uuid.IsValidUUID(s)
}

func isAlphanum(s string) bool {
	return alphanum.MatchString(s)
}

// ValidateUUID4 is a custom validator function that checks if a field is a UUIDv4.
func ValidateUUID4(fl validator.FieldLevel) bool {
	return isUUID4(fl.Field().String())
}

// ValidateUUID4OrNone is a custom validator function that checks if a field is a UUIDv4 or empty.
func ValidateUUID4OrNone(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	return isUUID4(fl.Field().String())
}

// ValidateIsSet is a custom validator function that checks if a field is set.
func ValidateIsSet(fl validator.FieldLevel) bool {
	return fl.Field().String() != ""
}

// ValidateAlphanum is a custom validator function that checks if a
// field is alphanumeric.
func ValidateAlphanum(fl validator.FieldLevel) bool {
	return isAlphanum(fl.Field().String())
}

// ValidateNoOp is a custom validator function that does nothing.
func ValidateNoOp(fl validator.FieldLevel) bool {
	return true
}

// ValidateAlphanumOrNone is a custom validator function that checks if a field is alphanumeric or empty.
func ValidateAlphanumOrNone(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	return isAlphanum(fl.Field().String())
}

// ValidateAlphanumSymbol is a custom validator function that checks if a field is alphanumeric with symbols.
func ValidateAlphanumSymbol(fl validator.FieldLevel) bool {
	return alphanumsym.MatchString(fl.Field().String())
}
