package structs

const (
	RandomMax = 1000000
)

// relationshipGender returns if
// - the relationship is gendered
// - if the relationshipGender isMale
// Ie.
//   - returns (true, true) for PersonalRelation_Husband
//     returns (true, false) for PersonalRelation_Wife
//     returns (false, false) for PersonalRelation_Friend
func relationshipGender(r PersonalRelation) (bool, bool) {
	if r == PersonalRelation_ExWife || r == PersonalRelation_Wife || r == PersonalRelation_Mother || r == PersonalRelation_Daughter || r == PersonalRelation_Grandmother {
		return true, false
	}
	if r == PersonalRelation_ExHusband || r == PersonalRelation_Husband || r == PersonalRelation_Father || r == PersonalRelation_Son || r == PersonalRelation_Grandfather {
		return true, true
	}
	return false, false
}

// flipRelationship returns the relationship from the other persons perspective.
// This can be the same as flipRelationshipGender, but is not always.
// Ie. grandson -> grandfather
func flipRelationship(r PersonalRelation) PersonalRelation {
	switch r {
	case PersonalRelation_ExWife:
		return PersonalRelation_ExHusband
	case PersonalRelation_ExHusband:
		return PersonalRelation_ExWife
	case PersonalRelation_Wife:
		return PersonalRelation_Husband
	case PersonalRelation_Husband:
		return PersonalRelation_Wife
	case PersonalRelation_Father:
		return PersonalRelation_Son
	case PersonalRelation_Mother:
		return PersonalRelation_Daughter
	case PersonalRelation_Son:
		return PersonalRelation_Father
	case PersonalRelation_Daughter:
		return PersonalRelation_Mother
	case PersonalRelation_Grandmother:
		return PersonalRelation_Granddaughter
	case PersonalRelation_Grandfather:
		return PersonalRelation_Grandson
	case PersonalRelation_Grandson:
		return PersonalRelation_Grandfather
	case PersonalRelation_Granddaughter:
		return PersonalRelation_Grandmother
	}
	return r
}

// flipRelationshipGender returns the opposite gendered relationship.
// Ie. Son -> Daughter, Husband -> Wife, etc.
// That is we flip male -> female, female -> male
func flipRelationshipGender(r PersonalRelation) PersonalRelation {
	switch r {
	case PersonalRelation_ExWife:
		return PersonalRelation_ExHusband
	case PersonalRelation_ExHusband:
		return PersonalRelation_ExWife
	case PersonalRelation_Wife:
		return PersonalRelation_Husband
	case PersonalRelation_Husband:
		return PersonalRelation_Wife
	case PersonalRelation_Father:
		return PersonalRelation_Mother
	case PersonalRelation_Mother:
		return PersonalRelation_Father
	case PersonalRelation_Son:
		return PersonalRelation_Daughter
	case PersonalRelation_Daughter:
		return PersonalRelation_Son
	case PersonalRelation_Grandmother:
		return PersonalRelation_Grandfather
	case PersonalRelation_Grandfather:
		return PersonalRelation_Grandmother
	case PersonalRelation_Grandson:
		return PersonalRelation_Granddaughter
	case PersonalRelation_Granddaughter:
		return PersonalRelation_Grandson
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
//     PersonalRelation_Wife,
//     PersonalRelation_Husband,
//     PersonalRelation_Lover,
//     PersonalRelation_Colleague,
//     } -> PersonalRelation_Husband
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
