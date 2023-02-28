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
	// Faction ID of whomever has ultimate say over this
	GoverningFaction string

	// Faction ID of whomever has permission to use this
	// ie. a government has granted the right to use the land
	ControllingFaction string

	// Area ID
	Area string

	// Resource (commodity name)
	Resource string
}
