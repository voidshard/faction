package structs

// Count returns the number of people of a given rank.
func (d *DemographicRankSpread) Count(r FactionRank) int64 {
	// TODO: I'm sure there's a smarter way
	switch r {
	case FactionRank_Ruler:
		return d.Ruler
	case FactionRank_Elder:
		return d.Elder
	case FactionRank_GrandMaster:
		return d.GrandMaster
	case FactionRank_Master:
		return d.Master
	case FactionRank_Expert:
		return d.Expert
	case FactionRank_Adept:
		return d.Adept
	case FactionRank_Journeyman:
		return d.Journeyman
	case FactionRank_Novice:
		return d.Novice
	case FactionRank_Apprentice:
		return d.Apprentice
	}
	return d.Associate
}

// NewDemographics
func NewDemographics() *Demographics {
	return &Demographics{
		Faith:       map[string]*DemographicStatSpread{}, // ReligionID -> DemographicStatSpread
		Profession:  map[string]*DemographicStatSpread{}, // Profession -> DemographicStatSpread
		Affiliation: map[string]*DemographicStatSpread{}, // FactionID -> DemographicStatSpread
		Rank:        map[string]*DemographicRankSpread{}, // FactionID -> DemographicRankSpread
	}
}

func (d *Demographics) AddFaith(in string, val int64) {
	if _, ok := d.Faith[in]; !ok {
		d.Faith[in] = &DemographicStatSpread{}
	}
	d.Faith[in].Add(val)
}

func (d *Demographics) AddProfession(in string, val int64) {
	if _, ok := d.Profession[in]; !ok {
		d.Profession[in] = &DemographicStatSpread{}
	}
	d.Profession[in].Add(val)
}

func (d *Demographics) AddAffiliation(in string, val int64) {
	if _, ok := d.Affiliation[in]; !ok {
		d.Affiliation[in] = &DemographicStatSpread{}
	}
	d.Affiliation[in].Add(val)
}

func (d *Demographics) AddRank(in string, val FactionRank) {
	if _, ok := d.Rank[in]; !ok {
		d.Rank[in] = &DemographicRankSpread{}
	}
	d.Rank[in].Add(val, 1)
}

func (d *DemographicRankSpread) Add(val FactionRank, i int64) {
	d.Total += i
	switch val {
	case FactionRank_Ruler:
		d.Ruler += i
	case FactionRank_Elder:
		d.Elder += i
	case FactionRank_GrandMaster:
		d.GrandMaster += i
	case FactionRank_Master:
		d.Master += i
	case FactionRank_Expert:
		d.Expert += i
	case FactionRank_Adept:
		d.Adept += i
	case FactionRank_Journeyman:
		d.Journeyman += i
	case FactionRank_Novice:
		d.Novice += i
	case FactionRank_Apprentice:
		d.Apprentice += i
	case FactionRank_Associate:
		d.Associate += i
	}
}

func (d *DemographicStatSpread) Add(in int64) {
	v := float64(in)
	switch {
	case v >= float64(MaxTuple)*0.95:
		d.Exemplary++
	case v >= float64(MaxTuple)*0.75:
		d.Excellent++
	case v >= float64(MaxTuple)*0.5:
		d.Good++
	case v >= float64(MaxTuple)*0.25:
		d.Fine++
	case v >= 0:
		d.Average++
	case v >= float64(MinTuple)*0.25:
		d.Poor++
	case v >= float64(MinTuple)*0.5:
		d.Awful++
	case v >= float64(MinTuple)*0.75:
		d.Terrible++
	default:
		d.Abysmal++
	}
	d.Total++
}
