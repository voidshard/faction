package economy

// Economy is a primitive way for our simulation to understand enough about the world
// economy to make hopefully not irrational decisions.
//
// We provide a simple example but users are encouraged to supply their own
// when kicking off a simulation. In light of that, we try to keep this interface fairly
// straight forward to implement.
type Economy interface {
	// CommodityValue returns the value (or forecasted value) of a (single) commodity in
	// some area at some time offset in ticks (ie. '0' is 'now').
	// Ie. the value of 1 ingot of iron
	//     the value of 1 bushel of wheat
	//     ...etc
	CommodityValue(commodity, area string, ticks int64) float64

	// CommodityYield returns the yield (or forecasted yield) of a commodity in
	// some area (1 unit squared) at some time offset in ticks (ie. '0' is 'now').
	//
	// The returned value is some number representing how good the yield is.
	// Ie. 1.0 is average, 2.0 is double, 0.5 is half, etc.
	//
	// This together with Plot (Area - Commodity - Yield) is used to determine
	// how successful a venture (harvest / craft) is.
	CommodityYield(commodity, area string, ticks int64) float64

	// LandValue returns the value of 1 unit squared of land in the area at some
	// time offset in ticks (ie. '0' is 'now').
	//
	// This might be a farm + attached building(s) in a rural area, a large complex
	// within a city area, a small apartment or .. whatever.
	LandValue(area string, ticks int64) float64
}
