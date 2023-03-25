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

func (d *DemographicStats) Add(v int) {
	switch {
	case v >= 95:
		d.Exemplary++
	case v >= 75:
		d.Excellent++
	case v >= 50:
		d.Good++
	case v >= 25:
		d.Fine++
	case v >= 0:
		d.Average++
	case v >= -25:
		d.Poor++
	case v >= -50:
		d.Awful++
	case v >= -75:
		d.Terrible++
	default:
		d.Abysmal++
	}
}
