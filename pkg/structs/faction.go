package structs

// Faction represents some group we would like to simulate.
// Nb. we don't assume these are the *only* factions, just that they're the
// most notable / influential / interesting.
type Faction struct {
	Ethos // average ethics of faction members (by tradition / current leadership)

	ID   string `db:"id"`
	Name string `db:"name"`

	HomeAreaID string `db:"home_area_id"` // the area where the faction is based
	HQPlotID   string `db:"hq_plot_id"`   // faction headquarters

	ActionFrequencyTicks int `db:"action_frequency_ticks"` // faction offers new jobs every X ticks

	Leadership LeaderType      `db:"leadership"` // how faction is run
	Structure  LeaderStructure `db:"structure"`  // how faction is structured

	Wealth     int `db:"wealth"`     // money/liquid wealth the faction has available to spend. Min:0
	Cohesion   int `db:"cohesion"`   // a metric for how united a faction is. Min:0, Max:MaxTuple
	Corruption int `db:"corruption"` // 'asset' corruption or how much money tends to go inexplicably missing. Min:0, Max:MaxTuple

	IsCovert bool `db:"is_covert"` // the faction actively avoids notice, discourages public action

	GovernmentID string `db:"government_id"` // ID of Government that this faction is under (if only nominally)
	IsGovernment bool   `db:"is_government"` // a government can modify the Government struct

	ReligionID string `db:"religion_id"` // ID of Religion (TODO)
	IsReligion bool   `db:"is_religion"` // organised religion, not simply *has* a religion (ie. church vs. order of knights)

	IsMemberByBirth bool `db:"is_member_by_birth"` // if having parent(s) in the faction auto joins children

	EspionageOffense int `db:"espionage_offense"` // how good the faction is at spying on others, Min:0
	EspionageDefense int `db:"espionage_defense"` // how good the faction is at defending against spying, Min:0
	MilitaryOffense  int `db:"military_offense"`  // how good the faction is at attacking others. Min:0
	MilitaryDefense  int `db:"military_defense"`  // how good the faction is at defending against attacks. Min:0

	ParentFactionID       string          `db:"parent_faction_id"`       // ID of parent faction (if any)
	ParentFactionRelation FactionRelation `db:"parent_faction_relation"` // relation to parent faction (if any)

	// Numbers are best-effort calculated as part of a faction's Job decision process
	Members int `db:"members"` // number of members in the faction
	Plots   int `db:"plots"`   // number of plots owned by the faction
	Areas   int `db:"areas"`   // number of areas the faction is present in
}

// FactionSummary is a high level overview of a faction, including related tuples
// (weights), information on current faction leadership (Ranks) and any research (ResearchProgress).
type FactionSummary struct {
	Faction

	// ammassed research (per topic)
	ResearchProgress map[string]int

	// weights
	Professions map[string]int
	Actions     map[ActionType]int
	Research    map[string]int
	Trust       map[string]int

	// counts of people in each rank
	Ranks *DemographicRankSpread
}

func (fs *FactionSummary) ToFaction() *Faction {
	return &(fs).Faction
}

func NewFactionSummary(f *Faction) *FactionSummary {
	return &FactionSummary{
		Faction:          *f,
		ResearchProgress: map[string]int{},
		Professions:      map[string]int{},
		Actions:          map[ActionType]int{},
		Research:         map[string]int{},
		Trust:            map[string]int{},
		Ranks:            &DemographicRankSpread{},
	}
}
