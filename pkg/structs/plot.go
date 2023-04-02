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
}
