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
	// Commodities returns a list of all commodities known to the economy
	// in the given area.
	Commodities(area string) []*structs.Commodity

	// Commodity looks up a commodity by name & returns it
	Commodity(commodity string) *structs.Commodity

	// CommodityValue returns the value (or forecasted value) of a commodity in
	// some area at some time offset in ticks (ie. '0' is 'now').
	CommodityValue(commodity, area string, ticks int) float64

	// CommodityYield returns the yield (or forecasted yield) of a commodity in
	// some area at some time offset in ticks (ie. '0' is 'now').
	//
	// The returned value is some number representing how good the yield is.
	// Ie. 1.0 is average, 2.0 is double, 0.5 is half, etc.
	//
	// This together with LandRight (Area - Commodity - Yield) is used to determine
	// how successful a venture (harvest / craft) is.
	CommodityYield(commodity, area string, ticks, professionSkill int) float64

	// LandValue returns the value of 1 unit squared of land in the area.
	// This might be a farm + attached building(s) in a rural area, a large complex
	// within a city area, a small apartment or .. whatever.
	LandValue(area string, ticks int) float64

	// True if the commodity is produced by crafting, generally
	// this requires input commodities.
	// Ie. iron ingots
	IsCraftable(commodity string) bool

	// True if the commodity is produced by harvesting, generally
	// base resources.
	// Ie. iron ore
	IsHarvestable(commodity string) bool
}

func commodityToProfession(eco Economy, areas ...string) map[string]string {
	prof := map[string]string{}
	for _, a := range areas {
		for _, c := range eco.Commodities(a) {
			prof[c.Name] = c.Profession
		}
	}
	return prof
}
