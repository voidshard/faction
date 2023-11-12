package simutil

import (
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/internal/random/rng"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

func NewJob(tick int, action string, cfg *config.Action) *structs.Job {
	preptime := rng.NewRand(
		cfg.TimeToPrepare.Min,
		cfg.TimeToPrepare.Max,
		cfg.TimeToPrepare.Mean,
		cfg.TimeToPrepare.Deviation,
	).Int()
	exetime := rng.NewRand(
		cfg.TimeToExecute.Min,
		cfg.TimeToExecute.Max,
		cfg.TimeToExecute.Mean,
		cfg.TimeToExecute.Deviation,
	).Int()
	return &structs.Job{
		ID: dbutils.RandomID(),
		// 		ParentJobID: "",
		// 		SourceFactionID: "",
		// 		SourceAreaID:    "",
		Action:   action,
		Priority: cfg.JobPriority,
		// 		TargetFactionID: "",
		// 		TargetAreaID:    "",
		// 		TargetMetaKey:   "",
		// 		TargetMetaVal:   "",
		PeopleMin: cfg.MinPeople,
		PeopleMax: cfg.MaxPeople,
		//		PeopleNow:  0,
		TickCreated: tick,
		TickStarts:  tick + preptime,
		TickEnds:    tick + preptime + exetime,
		//		Secrecy:     0,
		//		IsIllegal:   false,
		State: structs.JobStatePending,
	}
}
