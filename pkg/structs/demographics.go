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
	Exemplary int // 9500 -> 10000
	Excellent int // 7500 -> 9500
	Good      int // 5000 -> 7500
	Fine      int // 2500 -> 5000
	Average   int // 0 -> 2500
	Poor      int // 0 -> -2500
	Awful     int // -2500 -> -5000
	Terrible  int // -5000 -> -7500
	Abysmal   int // -7500 -> -10000

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

// Count returns the number of people of a given rank.
func (d *DemographicRankSpread) Count(r FactionRank) int {
	// TODO: I'm sure there's a smarter way
	switch r {
	case FactionRankRuler:
		return d.Ruler
	case FactionRankElder:
		return d.Elder
	case FactionRankGrandMaster:
		return d.GrandMaster
	case FactionRankMaster:
		return d.Master
	case FactionRankExpert:
		return d.Expert
	case FactionRankAdept:
		return d.Adept
	case FactionRankJourneyman:
		return d.Journeyman
	case FactionRankNovice:
		return d.Novice
	case FactionRankApprentice:
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
	d.Rank[in].Add(val, 1)
}

func (d *DemographicRankSpread) Add(val FactionRank, i int) {
	d.Total += i
	switch val {
	case FactionRankRuler:
		d.Ruler += i
	case FactionRankElder:
		d.Elder += i
	case FactionRankGrandMaster:
		d.GrandMaster += i
	case FactionRankMaster:
		d.Master += i
	case FactionRankExpert:
		d.Expert += i
	case FactionRankAdept:
		d.Adept += i
	case FactionRankJourneyman:
		d.Journeyman += i
	case FactionRankNovice:
		d.Novice += i
	case FactionRankApprentice:
		d.Apprentice += i
	case FactionRankAssociate:
		d.Associate += i
	}
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
