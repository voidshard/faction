package v1

// Culture imparts information about the culture someone lives in and castes or social
// groupings of that culture.
type Culture struct {
	Meta `json:",inline" yaml:",inline"`

	// Base set default values for all Castes, values set in Castes overwrite these.
	Base CultureCaste `yaml:"Base" json:"Base" validate:"dive"`

	// 'Castes' here are culturally significant categories. Depending on the culture
	// it may be possible for an individual to change caste, or it may be fixed at birth.
	//
	// Culture here can influence a lot of things - from what professions are available,
	// whether one can own property, marry, have children etc.
	Castes map[string]CultureCaste `yaml:"Castes" json:"Castes" validate:"min=1,max=100,dive,keys,alphanum,endkeys,dive"`
}

func (x *Culture) New(in interface{}) (Object, error) {
	i := &Culture{}
	err := unmarshalObject(in, i)
	i.Kind = "culture"
	return i, err
}

type CultureCaste struct {
	// Meta sets values on all Actors using this Caste
	Set SetMeta `yaml:"Set" json:"Set"`

	// NamingScheme is the naming scheme used for this caste.
	//
	NamingScheme string `yaml:"NamingScheme" json:"NamingScheme" validate:"alphanum-or-empty"`

	// FamilyStructure
	// Determines 'acceptable' family structures for this caste.
	Family map[string]FamilyStructure `yaml:"Family" json:"Family" validate:"min=1,max=100,dive,keys,alphanum,endkeys,dive"`
	// Friendship

	// Faith

	// Death

	// Professions

	// Actions that actors of this culture + caste prefer to take.
	//
	// In general we expect the vast majority of actions to be kicked off by factions, but this
	// allows cultures as a group to influence events. It also helps with the formation of new
	// factions where many founding members are of this culture; as the founded faction will inherit
	// views of the founders.
	Actions ActionSelection `yaml:"Actions" json:"Actions" validate:"dive"`
}

// FamilyStructure defines how a culture expects a "family" to be structured.
//
// This might imply details about the race, culture, caste or anything else really.
type FamilyStructure struct {
	// Meta sets values on all Families using this structure
	Set SetMeta `yaml:"Set" json:"Set"`

	// Probability of a family being structured this way.
	Probability float64 `yaml:"Probability" json:"Probability" validate:"gte=0,lte=1"`

	// Adults selected for the family
	// Eg.
	//  "Mother" "Father" for a human style family
	//
	// Or perhaps.
	//  "MaleParasite" "FemaleParasite" "CarrierCreature" "HostCreature"
	// To represent a parasitic species that are eaten by 'Carrier' which are in turn
	// eaten by a 'Host' where they reproduce.
	Adults map[string]FamilyUnit `yaml:"Adults" json:"Adults" validate:"min=1,max=100,dive,keys,alphanum,endkeys,dive"`

	// Lifespan for the family; ie. after this the family breaks up (or divorces).
	// This might last for the duration of the adults lives, or only long enough to
	// have children.
	Lifespan Distribution `yaml:"Lifespan" json:"Lifespan" validate:"dive"`

	// ChildrenRace is the 'adult' from which race is taken.
	// Ie. where Adults are { "Mother", "Father" } this might be "Mother".
	ChildrenRace string `yaml:"ChildrenRace" json:"ChildrenRace" validate:"alphanum"`

	// ChildrenPerCycle is the number of children produced per breeding cycle.
	ChildrenPerCycle Distribution `yaml:"ChildrenPerCycle" json:"ChildrenPerCycle" validate:"dive"`

	// BreedingCycles is the number of breeding cycles that will occur over the family lifespan.
	BreedingCycles Distribution `yaml:"BreedingCycles" json:"BreedingCycles" validate:"dive"`

	// TicksBetweenCycles is the number of ticks between breeding cycles.
	TicksBetweenCycles Distribution `yaml:"TicksBetweenCycles" json:"TicksBetweenCycles" validate:"dive"`
}

// FamilyUnit defines a sub unit of the Adult(s) in a family.
//
// Eg.
//
//	Mother, Father would each be a unit in a human style family.
//	Father (selects 1), Mother(s) (selects up to 10) might each be a unit similar to a lion pride.
type FamilyUnit struct {
	// How to pick the Actor(s) for this unit.
	Select Query `yaml:"Select" json:"Select" validate:"dive"`

	// Weight given to this unit (when making decisions in the family).
	// Ie. if a family has 2 adults, and one has a weight of 0.8 and the other 0.2 then the
	// 0.8 adult's decisions carry much more weight than the other.
	Weight float64 `yaml:"Weight" json:"Weight" validate:"gte=0,lte=1"`
}
