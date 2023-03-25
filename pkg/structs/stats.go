package structs

// DemographicStats holds counts of some stats that appear in
// some population.
type DemographicStats struct {
	Exemplary int // 75 -> 100
	Excellent int // 50 -> 75
	Good      int // 25 -> 50
	Fine      int // 5 -> 25
	Average   int // -5 -> 5
	Poor      int // -5 -> -25
	Awful     int // -25 -> -50
	Terrible  int // -50 -> -75
	Absymal   int // -75 -> -100
}

func (d *DemographicStats) Add(v int) {
	switch {
	case v >= 75:
		d.Exemplary++
	case v >= 50:
		d.Excellent++
	case v >= 25:
		d.Good++
	case v >= 5:
		d.Fine++
	case v >= -5:
		d.Average++
	case v >= -25:
		d.Poor++
	case v >= -50:
		d.Awful++
	case v >= -75:
		d.Terrible++
	default:
		d.Absymal++
	}
}
