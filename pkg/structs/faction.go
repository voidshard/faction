package structs

func (f *Faction) ObjectID() string {
	return f.ID
}

func NewFactionSummary(f *Faction) *FactionSummary {
	return &FactionSummary{
		Faction:          f,
		ResearchProgress: map[string]int64{},
		Professions:      map[string]int64{},
		Actions:          map[string]int64{},
		Research:         map[string]int64{},
		Trust:            map[string]int64{},
		Ranks:            &DemographicRankSpread{},
	}
}
