package structs

type FactionRelation int

// FactionRelation denotes how a faction relates to it's parent faction (if any)
const (
	FactionRelationMember    FactionRelation = iota // equal member of a super faction (ie. members of the EU)
	FactionRelationDivision                         // a specialization ie. a government spawns a specialist spy faction
	FactionRelationTributary                        // pays tax to super faction, but does it's own thing
	FactionRelationVassal                           // strictly subordinate, but with regin over it's own affairs
	FactionRelationPuppet                           // strictly subordinate, very little self rule
)
