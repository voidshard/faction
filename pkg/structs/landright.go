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

	// Faction ID of whomever has permission to use this
	// ie. a government has granted the right to use the land
	FactionID string `db:"faction_id"`

	// Area ID
	AreaID string `db:"area_id"`
}
