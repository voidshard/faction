package structs

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
	ID string `db:"id"`

	// Government ID (if any)
	GovernmentID string `db:"government_id"`

	// something to denote the area
	// - in general an unset biome implies it's suitable for anything / catch all
	Biome string `db:"biome"`

	// Random value for getting areas are random from the DB
	Random int `db:"random"`
}

func (a *Area) ObjectID() string {
	return a.ID
}
