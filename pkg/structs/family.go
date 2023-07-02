package structs

// Family is created when one of
// - simulation.go 'Populate' creates a family
// - a male / female of child bearing age marry
// - a male / female of child bearing age who are lovers have a child (ie. outside of wedlock)
//
// For as long as both people are lovers / married & within child-bearing age a family
// has a chance of bearing children every so often.
type Family struct {
	Ethos // rough outlook

	ID string `db:"id"`

	// Race demographic assigned to children
	Race string `db:"race"`

	// Area where the family is based (where children will be placed)
	AreaID string `db:"area_id"`

	// Faction ID (if any) if this family is simulated as a major player.
	//
	// This implies the family is fairly wealthy and/or influential, probably 95% of families
	// will not have this set; which is probably a good thing and saves us lots of calculations
	// for families which don't really have the resources to act on the national / international
	// stage.
	FactionID string `db:"faction_id"`

	// True while;
	// - both people are capable of bearing children
	// - both people are married or lovers (ie. willing to bear children)
	IsChildBearing bool `db:"is_child_bearing"`

	// Represents the tick when one of the two potential parents becomes too old to bear children.
	MaxChildBearingTick int `db:"max_child_bearing_tick"`

	// If mother is pregnant then this is the tick when she will give birth.
	// Nb.
	// - we always set a Family to *not* child bearing if either partner dies
	//   (as in, it cannot produce more children)
	// - we set PregnancyEnd to 0 when the child is born or the mother dies
	//   (this means the child can be born if the father is dead)
	// This saves us having to query the parents when doing calculations
	PregnancyEnd int `db:"pregnancy_end"`

	// A family consists of a male & female and can bear children.
	// Nb. this does not imply that the couple are married ..
	MaleID   string `db:"male_id"`
	FemaleID string `db:"female_id"`

	// Save us looking up families (the info doesn't change anyways)
	MaGrandmaID string `db:"ma_grandma_id"` // mother's mother
	MaGrandpaID string `db:"ma_grandpa_id"` // mother's father
	PaGrandmaID string `db:"pa_grandma_id"` // father's mother
	PaGrandpaID string `db:"pa_grandpa_id"` // father's father

	NumberOfChildren int `db:"number_of_children"` // number of children this family has had

	MarriageTick int `db:"marriage_tick"` // tick when the couple married
	DivorceTick  int `db:"divorce_tick"`  // tick when the couple divorced
	WidowedTick  int `db:"widowed_tick"`  // tick when the couple were widowed
}

// relationshipGender returns if
// - the relationship is gendered
// - if the relationshipGender isMale
// Ie.
//   - returns (true, true) for PersonalRelationHusband
//     returns (true, false) for PersonalRelationWife
//     returns (false, false) for PersonalRelationFriend
func relationshipGender(r PersonalRelation) (bool, bool) {
	if r == PersonalRelationExWife || r == PersonalRelationWife || r == PersonalRelationMother || r == PersonalRelationDaughter || r == PersonalRelationGrandmother {
		return true, false
	}
	if r == PersonalRelationExHusband || r == PersonalRelationHusband || r == PersonalRelationFather || r == PersonalRelationSon || r == PersonalRelationGrandfather {
		return true, true
	}
	return false, false
}

// flipRelationship returns the relationship from the other persons perspective.
// This can be the same as flipRelationshipGender, but is not always.
// Ie. grandson -> grandfather
func flipRelationship(r PersonalRelation) PersonalRelation {
	switch r {
	case PersonalRelationExWife:
		return PersonalRelationExHusband
	case PersonalRelationExHusband:
		return PersonalRelationExWife
	case PersonalRelationWife:
		return PersonalRelationHusband
	case PersonalRelationHusband:
		return PersonalRelationWife
	case PersonalRelationFather:
		return PersonalRelationSon
	case PersonalRelationMother:
		return PersonalRelationDaughter
	case PersonalRelationSon:
		return PersonalRelationFather
	case PersonalRelationDaughter:
		return PersonalRelationMother
	case PersonalRelationGrandmother:
		return PersonalRelationGranddaughter
	case PersonalRelationGrandfather:
		return PersonalRelationGrandson
	case PersonalRelationGrandson:
		return PersonalRelationGrandfather
	case PersonalRelationGranddaughter:
		return PersonalRelationGrandmother
	}
	return r
}

// flipRelationshipGender returns the opposite gendered relationship.
// Ie. Son -> Daughter, Husband -> Wife, etc.
// That is we flip male -> female, female -> male
func flipRelationshipGender(r PersonalRelation) PersonalRelation {
	switch r {
	case PersonalRelationExWife:
		return PersonalRelationExHusband
	case PersonalRelationExHusband:
		return PersonalRelationExWife
	case PersonalRelationWife:
		return PersonalRelationHusband
	case PersonalRelationHusband:
		return PersonalRelationWife
	case PersonalRelationFather:
		return PersonalRelationMother
	case PersonalRelationMother:
		return PersonalRelationFather
	case PersonalRelationSon:
		return PersonalRelationDaughter
	case PersonalRelationDaughter:
		return PersonalRelationSon
	case PersonalRelationGrandmother:
		return PersonalRelationGrandfather
	case PersonalRelationGrandfather:
		return PersonalRelationGrandmother
	case PersonalRelationGrandson:
		return PersonalRelationGranddaughter
	case PersonalRelationGranddaughter:
		return PersonalRelationGrandson
	}
	return r
}

// decideRelationship returns the strongest relationship between two people given
// a list of relationships to consider and the target person `b`.
//
// We also return the correct gendered term (if applicable).
//
// Ie.
//   - b {IsMale: true}
//     in {
//     PersonalRelationWife,
//     PersonalRelationHusband,
//     PersonalRelationLover,
//     PersonalRelationColleague,
//     } -> PersonalRelationHusband
//
// Because b isMale and Husband is the strongest relationship in the given set.
func decideRelationship(b *Person, in ...PersonalRelation) PersonalRelation {
	lowest := in[0]
	for i := 1; i < len(in); i++ {
		if in[i] < lowest {
			lowest = in[i]
		}
	}

	impliesGender, isMale := relationshipGender(lowest)
	if impliesGender && b.IsMale != isMale {
		return flipRelationshipGender(lowest)
	}

	return lowest
}
