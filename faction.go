package faction

// Faction represents some group we would like to simulate.
// Nb. we don't assume these are the *only* factions, just that they're the
// most notable / influential / interesting.
type Faction struct {
	ID string

	Ethos Ethos // average ethics of faction members (traditionally & leaders)

	Focus []ActionType // actions the faction favours

	Leadership LeaderType // how faction is run

	Wealth     int // money/liquid wealth the faction has available to spend
	Cohesion   int // a metric for how united a faction is
	Corruption int // 'asset' corruption or how much money tends to go inexplicably missing
	Property   []*Land

	IsCovert     bool // the faction actively avoids notice, discourages public action
	IsIllegal    bool // explicitly illegal (by design or running afoul of the government
	IsGovernment bool // governments come with some extra powers

	IsReligion bool   // organised religion, not simply *has* a religion (ie. church vs. order of knights)
	ReligionID string // if any

	MemberByBirth bool // if having parent(s) in the faction allows children to join

	EspionageOffense int // represent stored good(s), member training & general preparation
	EspionageDefense int
	MilitaryOffense  int
	MilitaryDefense  int

	ParentFactionID       string       // ID of parent faction (if any)
	ParentFactionRelation RelationType // relation to parent faction (if any)

	ResearchScience  int // the various kinds of research a faction can do
	ResearchTheology int
	ResearchMagic    int
	ResearchOccult   int
}
