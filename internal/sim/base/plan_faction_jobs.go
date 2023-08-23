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

	ranksHigh = []structs.FactionRank{
		structs.FactionRankRuler,
		structs.FactionRankElder,
		structs.FactionRankGrandMaster,
		structs.FactionRankMaster,
	}
	ranksMed = []structs.FactionRank{
		structs.FactionRankExpert,
		structs.FactionRankAdept,
	}
	ranksLow = []structs.FactionRank{
		structs.FactionRankJourneyman,
		structs.FactionRankNovice,
		structs.FactionRankApprentice,
	}
)

/*
=> faction base action probabilities
- base probability of various actions
- apply faction ethos modifies for actions (ie. pacifist factions don't approve of violence)

=> look up favoured actions
- apply faction modifiers for various actions
- these are the actions we're aiming to do if things are going well

=> determine allies and/or enemies
- look up faction-faction-trust table

=> prioritize potential actions
order:
- survival
a. wealth: we need money to operate
    i. we define "enough money" to be operating expenses for some number of favoured actions
b. membership: (growth) we need people to work for us
    i. we define "enough people" to be the number of people we need to do our favoured actions
    ii. if we're at war (or similar) we don't have "enough people" if we're significantly below our enemies
        espionage and/or military rating(s)
c. if we have active enemies:
    i. military_defense: we need to defend ourselves
	nb. covert factions do not consider military defense an option; the aim is not to be seen, not fight openly.
    ii. espionage_defense: we need to defend our secrets
d. cohesion: (stability) we need people to be happy or we risk splitting
e. corruption: we need to keep corruption down or we hemorage money / work / people
- goals
if we're not under some existential threat, we can pursue higher level goals

=> determine list of desired jobs
    i. we should remove jobs we realistically can't do (eg. war, avoid some huge bribe)
=> get faction's current jobs
=> remove currently active jobs from desired jobs (that have the same target & action)
=> jobs that are for survival can trump non-survival jobs & cause us to cancel otherwise planned work

*/

func (s *Base) PlanFactionJobs(factionID string) ([]*structs.Job, error) {
	ctx, err := simutil.NewFactionContext(s.dbconn, factionID)
	if err != nil {
		return nil, err
	}

	// probabilities of doing various actions
	weights := simutil.NewActionWeights(s.cfg.Actions)
	if ctx.Summary.IsReligion {
		weights.SetIsReligion()
	}
	if ctx.Summary.IsGovernment {
		weights.SetIsGovernment()
	}
	weights.WeightByEthos(&ctx.Summary.Ethos)

	// weight based on the law & how tradition focused / law abiding the faction is
	tradition := float64(ctx.Summary.Ethos.Tradition/structs.MinEthos) + 1.0 // 2.0 -> 0.0
	if ctx.Summary.IsCovert {
		// factions that deliberately avoid the law tend to do so for a reason
		// Nb. a value of 1.1 means a faction with tradition of X becomes Yx more likely to break the law
		// -10k -> 2.2x
		// -5k -> 1.65x
		// -2.5k -> 1.375x
		// -1k -> 1.21x
		// -500 -> 1.122x
		weights.WeightByIllegal(1.1*tradition, ctx.AllGovernments()...)
	} else {
		// law abiding factions
		// Nb. a value of 0.85 means this value will be > 1 for factions with Tradition < -1.8k (ie. they don't respect traditions)
		// meaning factions that have low respect for traditions become slightly more likely to break the law
		weights.WeightByIllegal(0.85*tradition, ctx.AllGovernments()...)
	}

	// people we have available for work
	people := estimateAvailablePeople(ctx.Summary.Ranks)

	// Get all jobs for this faction that haven't finished yet
	q := db.Q(
		db.F(db.SourceFactionID, db.Equal, factionID),
		db.F(db.JobState, db.In, incompleteJobs),
	).DisableSort()

	// TODO: implies we only check the first 1k jobs (or whatever the default limit is).
	// We  probably shouldn't have this many jobs in progress for a single faction.
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

func estimateAvailablePeople(d *structs.DemographicRankSpread) int {
	// We anticipate that low affiliation people will spend a lot of time working for other factions too.
	// Very high affiliation people will probably spend most of their time working for us.
	total := 0
	for _, r := range ranksHigh {
		total += (d.Count(r) * 4) / 5 // 80%
	}
	for _, r := range ranksMed {
		total += d.Count(r) / 2 // 50%
	}
	for _, r := range ranksLow {
		total += d.Count(r) / 5 // 20%
	}
	return total
}
