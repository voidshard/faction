package faction

// Area represents a logical area of land .. since our simulation cares about people
// and factions, it should probably represent somewhere people can live.
//
// The 'id' here is used with Commodity, Harvest and Craft (economy.go) in order
// to derive prices & yields for a given area.
//
// It's also the centre of all job calculations, which are open to people with
// affiliation to a given faction within a given area.
//
// All this is to say; make sure you recall the Area ids you give us, because we'll
// be using it a lot and asking for it back.
type Area struct {
	// a unique reference for this area
	ID string

	// Faction ID
	GoverningFaction string

	// Resources at this location that can (in theory) be harvested (ActionTypeHarvest).
	//
	// By default the rights to collect these fall to the area government, who then
	// sells / allocates the right to factions to utilize (ActionTypeGrantLand).
	Resources []string // Commodity name(s)
}
