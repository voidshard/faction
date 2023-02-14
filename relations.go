package faction

type RelationType int

// RelationType denotes how a faction relates to it's parent faction (if any)
const (
	RelationTypeMember    RelationType = iota // equal member of a super faction (ie. members of the EU)
	RelationTypeDivision                      // a specialization ie. a government spawns a specialist spy faction
	RelationTypeTributary                     // pays super faction, but otherwise does it's own thing
	RelationTypeVassal                        // strictly subordinate, but with regin over it's own affairs
	RelationTypePuppet                        // strictly subordinate, very little self rule
)
