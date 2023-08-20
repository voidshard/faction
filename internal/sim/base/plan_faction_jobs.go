package base

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/sim/simutil"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	incompleteJobs = []string{
		string(structs.JobStatePending),
		string(structs.JobStateReady),
		string(structs.JobStateActive),
	}
)

func (s *Base) PlanFactionJobs(factionID string) ([]*structs.Job, error) {
	ctx, err := simutil.NewFactionContext(s.dbconn, factionID)
	if err != nil {
		return nil, err
	}
	// nb. ctx.Summary.Ranks.Total is the total number of people in the faction of various ranks

	// Get all jobs for this faction that haven't finished yet
	q := db.Q(
		db.F(db.SourceFactionID, db.Equal, factionID),
		db.F(db.JobState, db.In, incompleteJobs),
	).DisableSort()

	// TODO: implies we only check the first 1k jobs (or whatever the default limit is).
	// We .. probably shouldn't have this many jobs in progress for a single faction.
	jobs, _, err := s.dbconn.Jobs("", q)
	if err != nil {
		return nil, err
	}

	plannedMin := 0
	plannedMax := 0
	currentMin := 0
	currentMax := 0
	for _, j := range jobs {
		act, ok := s.cfg.Actions[j.Action]
		if !ok { // removed from config ??
			continue
		}

		if j.State == structs.JobStateActive {
			currentMin += act.MinPeople
			currentMax += act.MaxPeople
		} else {
			plannedMin += act.MinPeople
			plannedMax += act.MaxPeople
		}
	}
}
