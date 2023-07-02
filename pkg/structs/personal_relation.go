package structs

type PersonalRelation int

// PersonalRelation represents how two people know each other.
//
// In general a person is linked with the strongest relevant link.
// Ie.
// - if a person is a lover and a friend, they are considered a lover.
// - if a person is a lover and a wife, they are considered a wife.
// Nb. these are in order, where the lower number(s) trump higher ones.
const (
	// family
	PersonalRelationExWife PersonalRelation = iota
	PersonalRelationExHusband
	PersonalRelationWife
	PersonalRelationHusband
	PersonalRelationFather
	PersonalRelationMother
	PersonalRelationSon
	PersonalRelationDaughter
	PersonalRelationBrother
	PersonalRelationSister
	PersonalRelationGrandmother
	PersonalRelationGrandfather
	PersonalRelationGrandson
	PersonalRelationGranddaughter
	PersonalRelationExtendedFamily // could indicate inlaws, cousins, uncles, aunts etc

	// near family
	PersonalRelationFiance
	PersonalRelationLover

	// enemies
	PersonalRelationHatedEnemy // implies bad history and dedication to ruining each other
	PersonalRelationEnemy      // implies some bad history

	// friends
	PersonalRelationCloseFriend // long time friends
	PersonalRelationFriend

	// acquaintances
	PersonalRelationMentor
	PersonalRelationColleague // people that work together
)
