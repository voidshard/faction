package structs

// Route represents a path between two areas.
//
// Generally we accept that there will be more than one path between two
// areas, but this denotes the fastest / most economical for the transport
// of people and/or goods. It's possible that for example the wealthy
// can afford some more expensive but much faster mode of transport .. but
// we only need to understand the average / run-of-the-mill route here.
//
// Routes are one-way, so you'll want to have two to denote a two-way path
// (which is normally the case). We do this because travel time can vary
// based on direction; eg. sailing with vs against prevailing winds,
// travelling up vs down from the mountains.
type Route struct {
	SourceAreaID string `db:"source_area_id"` // area ID
	TargetAreaID string `db:"target_area_id"` // area ID
	TravelTime   int    `db:"travel_time"`    // travel time in ticks
}
