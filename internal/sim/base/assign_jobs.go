package base

func (s *Base) AssignJobs(factionID string) error {
	// find all jobs this faction has pending
	// sort by priority
	// read people who do not have jobs and want to work for this faction
	// - assign them to the highest priority job in their area
	// - assign them to the highest priority job in their faction
	// if we run out of people, stop
	// if we still have people, but run out of jobs, read all jobs of factions we are
	// subordinate to and assign people to them in the same way

	// TODO: What about a person's profession?
	//

	/*
		var (
			jobs   []*structs.Job
			jtoken string
			people []*structs.Person
			ptoken string
			err    error
		)

		tick, err := s.dbconn.Tick()
		jq := db.Q(
			db.F(db.SourceFactionID, db.Equal, factionID),
			db.F(db.JobState, db.Equal, structs.JobStatePending),
		)
		pq := db.Q(
			db.F(db.PreferredFactionID, db.Equal, factionID),
			db.F(db.AdulthoodTick, db.Less, tick), // is an adult
			db.F(db.DeathTick, db.Equal, 0),       // is alive
			db.F(db.JobID, db.Equal, ""),          // has no job
		)

		for {
			jobs, jtoken, err = s.dbconn.Jobs(jtoken, jq)
		}
	*/

	return nil
}
