package structs

// Demographics includes information about various stats
// in some areas.
// eg. Faith, Professions etc
type Demographics struct {
	Faith       map[string]*DemographicStatSpread
	Profession  map[string]*DemographicStatSpread
	Affiliation map[string]*DemographicStatSpread
	Rank        map[string]*DemographicRankSpread
}

// DemographicStatSpread holds counts of some stats that appear in
// some population.
type DemographicStatSpread struct {
	// See: MaxTuple, MinTuple in tuple.go
	Exemplary int // 95+
	Excellent int // 75+
	Good      int // 50+
	Fine      int // 25+
	Average   int // 0+
	Poor      int // 0 -> -25
	Awful     int // -25 -> -50
	Terrible  int // -50 -> -75
	Abysmal   int // -75 -> -100

	Total int
}

// DemographicRankSpread holds counts of faction ranks among a population
// in a given faction.
type DemographicRankSpread struct {
	// See: FactionRank in faction_relations.go
	Ruler       int
	Elder       int
	GrandMaster int
	Master      int
	Expert      int
	Adept       int
	Journeyman  int
	Novice      int
	Apprentice  int
	Associate   int

	Total int
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

func (d *Demographics) AddFaith(in string, val int) {
	if _, ok := d.Faith[in]; !ok {
		d.Faith[in] = &DemographicStatSpread{}
	}
	d.Faith[in].Add(val)
}

func (d *Demographics) AddProfession(in string, val int) {
	if _, ok := d.Profession[in]; !ok {
		d.Profession[in] = &DemographicStatSpread{}
	}
	d.Profession[in].Add(val)
}

func (d *Demographics) AddAffiliation(in string, val int) {
	if _, ok := d.Affiliation[in]; !ok {
		d.Affiliation[in] = &DemographicStatSpread{}
	}
	d.Affiliation[in].Add(val)
}

func (d *Demographics) AddRank(in string, val FactionRank) {
	if _, ok := d.Rank[in]; !ok {
		d.Rank[in] = &DemographicRankSpread{}
	}
	d.Rank[in].Add(val)

}

func (d *DemographicRankSpread) Add(val FactionRank) {
	switch val {
	case FactionRankRuler:
		d.Ruler++
	case FactionRankElder:
		d.Elder++
	case FactionRankGrandMaster:
		d.GrandMaster++
	case FactionRankMaster:
		d.Master++
	case FactionRankExpert:
		d.Expert++
	case FactionRankAdept:
		d.Adept++
	case FactionRankJourneyman:
		d.Journeyman++
	case FactionRankNovice:
		d.Novice++
	case FactionRankApprentice:
		d.Apprentice++
	case FactionRankAssociate:
		d.Associate++
	}
	d.Total++
}

func (d *DemographicStatSpread) Add(in int) {
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
