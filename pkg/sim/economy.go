package sim

import (
	"github.com/voidshard/faction/pkg/structs"
)

// Economy is a primitive way for our simulation to understand enough about the world
// economy to make hopefully not irrational decisions.
//
// We provide a simple example but users are encouraged to supply their own
// when kicking off a simulation. In light of that, we try to keep this interface fairly
// straight forward to implement.
type Economy interface {
	// Commodity looks up a commodity by name & returns it
	Commodity(commodity string) *structs.Commodity

	// CommodityValue returns the value (or forecasted value) of a commodity in
	// some area at some time offset in ticks (ie. '0' is 'now').
	CommodityValue(commodity, area string, ticks int) float64

	// CommodityYield returns the yield (or forecasted yield) of a commodity in
	// some area at some time offset in ticks (ie. '0' is 'now').
	CommodityYield(commodity, area string, ticks, professionSkill int) float64

	// LandValue returns the value of 1 unit squared of land in the area
	LandValue(area string, ticks int) float64
}
