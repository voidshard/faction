package structs

// DemographicStats holds counts of some stats that appear in
// some population.
type DemographicStats struct {
	Exemplary int // 95+
	Excellent int // 75+
	Good      int // 50+
	Fine      int // 25+

	Average int // 0+

	Poor     int // 0 -> -25
	Awful    int // -25 -> -50
	Terrible int // -50 -> -75
	Abysmal  int // -75 -> -100
}

func (d *DemographicStats) Add(in int) {
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
}
