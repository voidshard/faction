package structs

type FactionRelation int

// FactionRelation denotes how a faction relates to it's parent faction (if any)
const (
	FactionRelationTributary FactionRelation = iota // pays tax to super faction, but does it's own thing
	FactionRelationPuppet                           // strictly subordinate, very little self rule
	FactionRelationVassal                           // strictly subordinate, but with regin over it's own affairs
	FactionRelationMember                           // equal member of a super faction (ie. members of the EU)
)

var (
	// All known faction relations
	FactionRelationsAll = []FactionRelation{
		FactionRelationTributary,
		FactionRelationPuppet,
		FactionRelationVassal,
		FactionRelationMember,
	}
)
