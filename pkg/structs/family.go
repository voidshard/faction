package structs

type PersonalRelation int

// PersonalRelation represents how two people know each other.
const (
	PersonalRelationSpouse      PersonalRelation = iota // married
	PersonalRelationFiance                              // about to be married
	PersonalRelationCloseFriend                         // long time friends
	PersonalRelationFriend
	PersonalRelationExtendedFamily // could indicate inlaws, cousins, uncles, aunts etc

	PersonalRelationLover

	PersonalRelationChild
	PersonalRelationParent

	PersonalRelationColleague // people that work together

	PersonalRelationEnemy      // implies some bad history
	PersonalRelationHatedEnemy // implies bad history and dedication to ruining each other
)

// Family is created when one of
// - simulation.go 'Populate' creates a family
// - a male / female of child bearing age marry
// - a male / female of child bearing age who are lovers have a child (ie. outside of wedlock)
//
// For as long as both people are lovers / married & within child-bearing age a family
// has a chance of bearing children every so often.
type Family struct {
	ID string

	// Area where the family is based (where children will be placed)
	AreaID string

	// Faction ID (if any) if this family is simulated as a major player.
	//
	// This implies the family is fairly wealthy and/or influential, probably 95% of families
	// will not have this set; which is probably a good thing and saves us lots of calculations
	// for families which don't really have the resources to act on the national / international
	// stage.
	FactionID string

	// True while;
	// - both people are capable of bearing children
	// - both people are married or lovers (ie. willing to bear children)
	IsChildBearing bool

	// A family consists of a male & female and can bear children.
	// Nb. this does not imply that the couple are married ..
	Male   string
	Female string
}
