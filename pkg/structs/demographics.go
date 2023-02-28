package structs

// Demographics roughly describes a large population.
//
// For randomly making societies that look "sort of like this."
type Demographics struct {
	// Decides how large to make family units
	FamilySizeAverage   int
	FamilySizeDeviation float64

	// Ethos represents the average outlook of members of the population
	Ethos *Ethos

	// EthosDeviation is the standard deviation of the populace from the average (above)
	EthosDeviation *EthosDeviation

	// Professions to allocate to people
	Professions []*Profession

	// SideProfessions is a list of probabilities such that a given person has the number
	// of professions indicated by the index (any yes, index 0 means "" or "no profession").
	//
	// - a person cannot have the same profession multiple times
	// - a person may have only one profession with `ValidSideProfession` False
	// - a person can have any number of (unique) professions with `ValidSideProfession` True
	//
	// Ie. given
	//   {Name: farmer, ValidSideProfession: True}
	//   {Name: miner, ValidSideProfession: True}
	//   {Name: scribe, ValidSideProfession: False}
	//   {Name: priest, ValidSideProfession: False}
	// A person may not have "scribe" and "priest" at the same time, but could have any other
	// combination (including no profession).
	//
	// Eg. given SideProfessions = []float64{0.2, 0.4, 0.3, 0.1}
	// 20% chance of "no profession"
	// 40% chance of 1 profession
	// 30% chance of 2 professions
	// 10% chance of 3 professions
	SideProfessions []float64
}
