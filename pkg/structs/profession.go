package structs

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
	Average int

	// Deviation how far from the average skill people typically are
	Deviation int

	// ValidSideProfession indicates that someone may have this as a second / third
	// profession.
	//
	// Phrased another way if `ValidSideProfession` is False it indicates that someone
	// prefers to do this profession exclusively.
	//
	// A profession might be a valid side profession if:
	// - this might be something many people have skill in (eg. farming in medieval society)
	// - or starting professions that people might be expected to graduate from
	//   > miner to merchant after someone earns enough cash to open a shop
	//   > farmer to temple acolyte after an epic conversion
	// - something someone might do in downtime because the primary profession can
	//   be a little off & on (eg. working as a guard between assassination contracts)
	ValidSideProfession bool
}
