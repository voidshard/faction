package structs

import (
	"fmt"
)

var (
	// ErrNotEnoughPlots is returned when there are not enough plots to satisfy a request
	// ie. when you ask for 100 units of plots with the "timber" commodity but only 20 are found.
	ErrNotEnoughPlots = fmt.Errorf("not enough plots were found to satisty the request")
)
