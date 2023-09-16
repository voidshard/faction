package simutil

import (
	"github.com/voidshard/faction/pkg/structs"
)

type FactionScale struct {
	*structs.Faction

	// people
	Members int

	// land
	Plots int

	// plot valuation
	PlotValue int
}
