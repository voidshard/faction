package base

func (s *Base) AssignJobs(factionID string) error {

	// - find all jobs that need filling, either by this faction or it's superiors (if any)
	// - prioritise jobs
	//   - war, crusade take priority; not dying is key
	//   - people prefer to work for the given faction (not a superior if possible)
	//   - people prefer to work in their home area(s)
	// - find unemployed people who would like to work for this faction in the jobs target area
	// - expand search to other areas

	return nil
}
