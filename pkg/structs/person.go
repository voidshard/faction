package structs

// Person roughly outlines someone that can belong to / work for a faction.
type Person struct {
	Ethos *Ethos // rough outlook

	Job string // ie. current job id (if any)

	Age int // current age, in ticks
}
