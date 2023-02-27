package dbutils

import (
	"regexp"
)

var (
	// Name accepts a superset of UUIDs - we allow any alpha / number and
	// '-' (hypen) symbols to appear anywhere but the start / end.
	//
	// must start with alpha / number
	// can contain any number of alpha, numbers or '-'
	// must end with alpha / number
	validName = regexp.MustCompile("^[0-9a-zA-Z]{1}[0-9a-zA-Z\\-]*[0-9a-zA-Z]{1}$")
)

// IsValidName returns if we consider a name valid
func IsValidName(in string) bool {
	return validName.MatchString(in)
}
