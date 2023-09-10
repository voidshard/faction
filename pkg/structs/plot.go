package structs

// Plot is some package of land, buildings and attachments that a faction
// can buy, sell or use as a place of work / trade / whatever.
//
// Plots may yield some resource (eg. a Commodity) or simply be a plot of
// land / building(s) that can be used for some purpose.
//
// It might be a farm + land, a castle complete with moat, a high rise building,
// a small sea-side house + jetty .. for the purposes of the simulation all that
// matters is that it can be used as a place of work, for whatever work means
// for that faction.
type Plot struct {
	ID string `db:"id"`

	AreaID    string `db:"area_id"`
	FactionID string `db:"faction_id"` // the owner

	// secrecy value if it isn't widely known who owns this plot
	Hidden int `db:"hidden"`

	// average value for this plot (to the owner)
	// - size * land value + commodity value * yield
	// Nb. rough estimation at last write given current market values
	Value int `db:"value"`

	Crop
}

// Crop holds extra data about the land & it's usage
type Crop struct {
	// Commodity that can be harvested from this land (if any)
	Commodity string `db:"commodity"`

	// Size in land units squared
	Size int `db:"size"`

	// Yield of the resource, ie how many "units" of `resource` are produced
	// (or expected to be produced) from this per unit squared of land.
	//
	// This is an average, the actual value is dictated by the Economy interface
	// for the given tick(s) when needed.
	//
	// (If Commodity is set, otherwise this is 0).
	//
	// Nb. this land could be a small but super productive area, or a
	// massive expanse. It doesn't really matter .. all we mean here is that
	// this land is productive for a given purpose, and can be owned & run
	// by a faction.
	Yield int `db:"yield"`
}

type LandSummary struct {
	// commodity -> crop
	Commodities map[string]*Crop

	// area -> LandSummary
	Areas map[string]*LandSummary

	TotalSize int
	Count     int
}

func (ls *LandSummary) Add(p *Plot) {
	// summary by commodity
	crop := ls.Commodities[p.Commodity]
	if crop == nil {
		crop = &Crop{}
	}

	crop.Size += p.Size
	crop.Yield += p.Yield

	ls.TotalSize += p.Size
	ls.Count += 1
	ls.Commodities[p.Commodity] = crop

	// summary by area
	areaSum, ok := ls.Areas[p.AreaID]
	if !ok {
		areaSum = NewLandSummary()
	}
	areaSum.Add(p)
	ls.Areas[p.AreaID] = areaSum
}

func NewLandSummary() *LandSummary {
	return &LandSummary{
		Commodities: map[string]*Crop{},
		Areas:       map[string]*LandSummary{},
	}
}
