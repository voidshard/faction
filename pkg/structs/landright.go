package structs

// LandRights are (Area, Commodity) pairs that the government has the right to.
// Ie.
//
//	Area FooTown - Iron Ore Mine
//	Area BobVillage - Wheat Field
//
// Normally government doesn't work these itself but sells them off to faction(s)
// to work them (or spawns divisions of the government to handle them).
type LandRight struct {
	ID string `db:"id"`

	// Faction ID of whomever has ultimate say over this
	GoverningFactionID string `db:"governing_faction_id"`

	// Faction ID of whomever has permission to use this
	// ie. a government has granted the right to use the land
	ControllingFactionID string `db:"controlling_faction_id"`

	// Area ID
	AreaID string `db:"area_id"`

	// Resource (commodity name)
	Resource string `db:"resource"`

	// Yield of the resource, ie how many "units" of `resource` are produced
	// (or expected to be produced) from this land.
	//
	// Nb. this land could be a small but super productive area, or a
	// massive expanse. It doesn't really matter .. all we mean here is that
	// this land is productive for a given purpose, and can be owned & run
	// by a faction.
	Yield int `db:"yield"`
}
