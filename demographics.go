package faction

// Demographics roughly describes a large population.
//
// For randomly making societies that look "sort of like this."
type Demographics struct {
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

// Profession currently serves to tie People's skillset to some Commodit(y/ies)
// useful in Harvest / Craft actions (see action.go economy.go).
//
// Note that the empty string profession "" always means "no profession" - such
// people (if any) cannot be assigned work.
type Profession struct {
	// Name of the profession (unique)
	Name string

	// Probability a given person has this profession
	Occurs float64

	// Average skill level (0-100)
	Average float64

	// Deviation how far from the average skill people typically are
	Deviation float64

	// ValidSideProfession indicates that someone may have this as a second / third
	// profession.
	// - this might be something many people have skill in (eg. farming in medieval society)
	// - or starting professions that people might be expected to graduate from
	//   > miner to merchant after someone earns enough cash to open a shop
	//   > farmer to temple acolyte after an epic conversion
	// - something someone might do in downtime because the primary profession can
	//   be a little off & on (eg. working as a guard between assassination contracts)
	//
	// Phrased another way if `ValidSideProfession` is False it indicates that someone
	// prefers to do this profession exclusively.
	ValidSideProfession bool
}
