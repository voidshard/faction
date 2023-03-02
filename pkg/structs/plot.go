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

	IsHeadQuarters bool `db:"is_headquarters"` // area HQ usually the main place of business, trade, whatever

	AreaID         string `db:"area_id"`
	OwnerFactionID string `db:"owner_faction_id"`

	Size int `db:"size"` // in land units
}
