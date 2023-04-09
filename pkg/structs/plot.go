package structs

// Plot is some package of land, buildings and attachments that a faction
// can buy, sell or use as a place of work / trade / whatever.
//
// It might be a farm + land, a castle complete with moat, a high rise building,
// a small sea-side house + jetty .. for the purposes of the simulation all that
// matters is that it can be used as a place of work, for whatever work means
// for that faction.
type Plot struct {
	ID string `db:"id"`

	AreaID    string `db:"area_id"`
	FactionID string `db:"faction_id"` // the owner

	Size int `db:"size"` // in land units

	// Commodity that can be harvested from this land (if any)
	Commodity string `db:"commodity"`

	// Yield of the resource, ie how many "units" of `resource` are produced
	// (or expected to be produced) from this land.
	//
	// (If Commodity is set, otherwise this is 0).
	//
	// Nb. this land could be a small but super productive area, or a
	// massive expanse. It doesn't really matter .. all we mean here is that
	// this land is productive for a given purpose, and can be owned & run
	// by a faction.
	Yield int `db:"yield"`
}
