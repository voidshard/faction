package base

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/sim/simutil"
	"github.com/voidshard/faction/pkg/structs"
)

func (s *Base) AssignJobs(factionID string) error {
	tick, err := s.dbconn.Tick()
	if err != nil {
		return err
	}

	conscription := factionID == ""
	jq := s.getAssignableJobsQ(factionID, conscription)

	var (
		jobs   []*structs.Job
		jtoken string
		fareas map[string]map[string]*structs.Area
	)

	for {
		// fetch the next set of jobs
		// TODO. strictly speaking we should sort by Priority here
		jobs, jtoken, err = s.dbconn.Jobs(jtoken, jq)
		if err != nil {
			return err
		}

		if conscription {
			// if conscription, re-get areas, since we don't know what factions we're assigning to
			factions := []string{}
			for _, j := range jobs {
				factions = append(factions, j.SourceFactionID)
			}
			newAreaData, err := s.dbconn.FactionAreas(false, factions...)
			if err != nil {
				return err
			}
			fareas = newAreaData
		}

		// order the jobs by priority
		jlist := simutil.NewJobPriorityList(jobs, s.cfg.Actions)
		pq := s.getAssignablePeopleQuery(tick, factionID, fareas)

		var (
			people []*structs.Person
			ptoken string
		)
		for { // alright, time to find our people
			people, ptoken, err = s.dbconn.People(ptoken, pq)
			if err != nil {
				return err
			}

			assign := map[string][]string{}
			toJob := map[string]*structs.Job{}

			for _, p := range people {
				var job *structs.Job
				if conscription {
					factionAreas, _ := fareas[p.PreferredFactionID]
					job = jlist.FindJobWithin(p, factionAreas)
				} else {
					job = jlist.FindJob(p)
				}

				assignments, ok := assign[job.ID]
				if !ok {
					assignments = []string{}
				}
				assign[job.ID] = append(assignments, p.ID)
				toJob[p.ID] = job
			}

			for jobID, assignees := range assign {
				count, err := s.dbconn.AssignJob(jobID, assignees)
				if err != nil {
					return err
				}
				j, ok := toJob[jobID]
				if ok {
					j.PeopleNow += count
				}
			}

			if ptoken == "" {
				break
			}
		}

		if jtoken == "" {
			break
		}
	}

	return nil
}

func (s *Base) getAssignablePeopleQuery(tick int, factionID string, fareas map[string]map[string]*structs.Area) *db.Query {
	if factionID == "" {
		return db.Q(
			db.F(db.AdulthoodTick, db.Less, tick), // is an adult
			db.F(db.DeathTick, db.Equal, 0),       // is alive
			db.F(db.JobID, db.Equal, ""),          // has no job
			// we can't iterate again for each faction (+ it's areas) (well, we can, but we don't want to)
		)
	}
	ids := []string{}
	areas, ok := fareas[factionID]
	if ok {
		for id := range areas {
			ids = append(ids, id)
		}
	}
	if len(ids) > 0 {
		return db.Q(
			db.F(db.PreferredFactionID, db.Equal, factionID),
			db.F(db.AdulthoodTick, db.Less, tick), // is an adult
			db.F(db.DeathTick, db.Equal, 0),       // is alive
			db.F(db.JobID, db.Equal, ""),          // has no job
			db.F(db.AreaID, db.In, ids),           // in one of the given areas
		)
	}
	return db.Q( // this really shouldn't happen
		db.F(db.PreferredFactionID, db.Equal, factionID),
		db.F(db.AdulthoodTick, db.Less, tick), // is an adult
		db.F(db.DeathTick, db.Equal, 0),       // is alive
		db.F(db.JobID, db.Equal, ""),          // has no job
	)
}

func (s *Base) getAssignableJobsQ(factionID string, conscription bool) *db.Query {
	if conscription {
		return db.Q(
			db.F(db.JobState, db.Equal, structs.JobStatePending),
			db.F(db.Conscription, db.Equal, true),
			db.F(db.PeopleNow, db.Less, db.PeopleMax),
		)
	}
	return db.Q(
		db.F(db.SourceFactionID, db.Equal, factionID),
		db.F(db.JobState, db.Equal, structs.JobStatePending),
		db.F(db.Conscription, db.Equal, false),
		db.F(db.PeopleNow, db.Less, db.PeopleMax),
	)
}
