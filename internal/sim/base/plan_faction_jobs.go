package base

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/sim/simutil"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

type actionTypeCategory int

const (
	Friendly actionTypeCategory = iota
	Neutral
	Unfriendly
)

var (
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

	rnggen = rand.New(rand.NewSource(time.Now().UnixNano()))
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
	tick, err := s.dbconn.Tick()
	if err != nil {
		return nil, err
	}

	// check if we're at war with other factions
	attacking, defending, err := s.determineFactionAtWar(tick, factionID)
	if err != nil {
		return nil, err
	}

	if len(attacking) > 0 || len(defending) > 0 {
		// if we're at war, not dying consumes most of our concerns
		return s.planFactionJobsWartime(tick, factionID, attacking, defending)
	} else {
		// if we're at peace, we have lots of competing concerns
		return s.planFactionJobsPeacetime(tick, factionID)
	}
}

func (s *Base) planFactionJobsPeacetime(tick int, factionID string) ([]*structs.Job, error) {
	ctx, err := simutil.NewFactionContext(s.dbconn, factionID)
	if err != nil {
		return nil, err
	}

	// people we have available for work
	people := estimateAvailablePeople(ctx.Summary.Ranks)
	availablePeople := people

	// Get all jobs for this faction that are active & wont finish next tick
	q := db.Q(
		db.F(db.SourceFactionID, db.Equal, factionID),
		db.F(db.JobState, db.In, []string{
			string(structs.JobStatePending),
			string(structs.JobStateReady),
			string(structs.JobStateActive),
		}),
		db.F(db.TickEnds, db.Greater, tick),
	).DisableSort()
	jobs, _, err := s.dbconn.Jobs("", q)
	if err != nil {
		return nil, err
	}
	inflight := map[string]bool{}
	for _, j := range jobs {
		if j.State == structs.JobStateActive {
			availablePeople -= j.PeopleNow
		}
		inflight[jobKey(j.TargetFactionID, j.Action)] = true
	}

	// summary of the faction's land
	land, err := s.dbconn.LandSummary(nil, []string{factionID})
	if err != nil {
		return nil, err
	}

	// how likely we are to do various actions
	weights, survivalConcerns, err := s.actionWeightsForFaction(ctx, land, people, availablePeople)
	if err != nil {
		return nil, err
	}

	if !survivalConcerns {
		// we'll only apply higher level goals if imminent death isn't a concern
		goals := determineFactionGoals(ctx)
		weights.WeightByGoal(s.cfg.Settings.GoalWeight, goals...)
	}

	// decide which actions we should do
	if weights.Normalise() <= 0 {
		// there are no actions we wish to do / can do
		return nil, nil
	}

	// pick actions to do, if we can't find valid targets we'll move on
	plans := []*structs.Job{}
	setTargets := []*structs.Job{}
	for {
		// pick out an action
		act := weights.Choose()
		cfg, ok := s.cfg.Actions[act]
		if !ok {
			continue
		}

		target := s.chooseJobTargetFaction(ctx, inflight, act, cfg)
		if target == "" {
			weights.WeightAction(0.0, act) // we can't do this action, remove it from the list
			if weights.Normalise() <= 0 {
				break // no more choices :(
			}
			continue
		}

		job := simutil.NewJob(tick, act, cfg)
		job.SourceFactionID = factionID
		job.TargetFactionID = target
		job.IsIllegal = ctx.Summary.IsCovert || simutil.IsIllegalAction(act, ctx.AllGovernments()...)
		if job.IsIllegal {
			act, _ := s.cfg.Actions[job.Action]
			job.Secrecy = int(float64(rnggen.Intn(structs.MaxTuple)+ctx.Summary.EspionageDefense) * act.SecrecyWeight)
		}

		if job.SourceFactionID == job.TargetFactionID {
			// set here as we already have all the faction info we need
			job.TargetAreaID = ctx.RandomArea(act)
			if job.Action == structs.ActionTypeResearch {
				job.TargetMetaKey = structs.MetaKeyResearch
				job.TargetMetaVal = ctx.RandomResearch()
			}
		}

		_, ok = structs.ActionTarget[act]
		if ok && job.TargetMetaKey == "" { // we need to set a more specific target
			setTargets = append(setTargets, job)
		}

		inflight[jobKey(job.TargetFactionID, act)] = true

		// reduce the number of available people / money by the actions min
		availablePeople -= cfg.MinPeople
		ctx.Summary.Wealth -= int(cfg.Cost.Min)
		if availablePeople <= 0 {
			break
		} else {
			weights.WeightByMinPeople(0.0, availablePeople)
			weights.WeightByCost(0.0, float64(ctx.Summary.Wealth))
			if weights.Normalise() <= 0 {
				break
			}
		}

		plans = append(plans, job)
	}

	if len(setTargets) > 0 {
		err = s.setSpecificJobTargets(setTargets)
		if err != nil {
			return nil, err
		}
	}

	// nb. save jobs, faction

	return plans, err
}

func (s *Base) setSpecificJobTargetsPerson(jobs []*structs.Job) error {
	factionIDs := []string{}
	for _, j := range jobs {
		factionIDs = append(factionIDs, j.TargetFactionID)
	}

	// put out top 10 people of each faction
	people, err := s.dbconn.FactionLeadership(10, factionIDs...)
	if err != nil {
		return err
	}

	for _, j := range jobs {
		targets, ok := people[j.TargetFactionID]
		if !ok {
			return fmt.Errorf("no people found for faction %s", j.TargetFactionID)
		}

		chosen := targets.Get(rnggen.Intn(targets.Total))
		if chosen == "" {
			return fmt.Errorf("no target person found for faction %s", j.TargetFactionID)
		}

		j.TargetMetaKey = structs.MetaKeyPerson
		j.TargetMetaVal = chosen
	}

	return nil
}

func (s *Base) setSpecificJobTargetsPlot(jobs []*structs.Job) error {
	factionIDs := []string{}

	for _, j := range jobs {
		switch j.Action {
		case structs.ActionTypeDownsize:
			// targeting our own land(s)
			factionIDs = append(factionIDs, j.SourceFactionID)
		default:
			// targeting lands of the target
			factionIDs = append(factionIDs, j.TargetFactionID)
		}
	}

	return nil
}

func (s *Base) setSpecificJobTargets(plans []*structs.Job) error {
	// bucket jobs based on what kind of target they need
	needed := map[structs.MetaKey][]*structs.Job{}
	for _, j := range plans {
		metakey, ok := structs.ActionTarget[j.Action]
		if !ok {
			return nil
		}

		sofar, ok := needed[metakey]
		if !ok {
			sofar = []*structs.Job{}
		}
		sofar = append(sofar, j)
		needed[metakey] = sofar
	}

	// for each target type, select targets (ie. batched by ActionType)
	for metakey, jobs := range needed {
		switch metakey {
		case structs.MetaKeyPlot:

		case structs.MetaKeyPerson:
			err := s.setSpecificJobTargetsPerson(jobs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func jobKey(targetID string, t structs.ActionType) string {
	return fmt.Sprintf("%s:%s", targetID, string(t))
}

func actionCategory(cfg *config.Action) actionTypeCategory {
	if cfg.TargetMinTrust < 0 {
		if cfg.TargetMaxTrust > 0 {
			return Neutral
		} else {
			return Unfriendly // min & max are both < 0
		}
	} else {
		return Friendly // min & max are both > 0
	}
}

// chooseJobTargetFaction picks a target for a job.
func (s *Base) chooseJobTargetFaction(ctx *simutil.FactionContext, inflight map[string]bool, act structs.ActionType, cfg *config.Action) string {
	relations := ctx.Relations()

	if cfg.TargetMinTrust == 0 && cfg.TargetMaxTrust == 0 {
		return ctx.Summary.ID // target ourselves
	}

	// for friendly actions, we order potential targets by how much we like them
	// for unfriendly actions, we order potential targets by how much we dislike them
	reverse := false
	if actionCategory(cfg) != Friendly {
		reverse = true
	}

	// choices here is sorted by Trust
	choices := relations.TrustBetween(cfg.TargetMinTrust, cfg.TargetMaxTrust, reverse)
	if len(choices) == 0 {
		return ""
	}

	// pick the first valid target that for which we don't already have a job
	for _, targetFactionID := range choices {
		_, ok := inflight[jobKey(targetFactionID, act)]
		if ok { // this job & target is already inflight
			continue
		}
		return targetFactionID
	}

	return ""
}

func (s *Base) planFactionJobsWartime(tick int, factionID string, attacking, defending []string) ([]*structs.Job, error) {
	return nil, nil
}

func (s *Base) determineFactionAtWar(tick int, factionID string) ([]string, []string, error) {
	q := db.Q(
		// people attacking us
		db.F(db.TargetFactionID, db.Equal, factionID),
		db.F(db.JobState, db.Equal, structs.JobStateActive),
		db.F(db.ActionType, db.In, []string{
			string(structs.ActionTypeWar),
			string(structs.ActionTypeCrusade),
			string(structs.ActionTypeShadowWar),
		}),
		db.F(db.TickEnds, db.Greater, tick),
	).Or(
		// people planning to attack us
		db.F(db.TargetFactionID, db.Equal, factionID),
		db.F(db.JobState, db.In, []string{
			string(structs.JobStatePending),
			string(structs.JobStateReady),
		}),
		db.F(db.ActionType, db.In, []string{
			string(structs.ActionTypeWar),
			string(structs.ActionTypeCrusade),
			// nb. shadow wars are covert, we don't know about them until they start
		}),
		db.F(db.TickEnds, db.Greater, tick),
	).Or(
		// people we're attacking
		db.F(db.SourceFactionID, db.Equal, factionID),
		db.F(db.ActionType, db.In, []string{
			string(structs.ActionTypeWar),
			string(structs.ActionTypeCrusade),
			string(structs.ActionTypeShadowWar),
		}),
		db.F(db.TickEnds, db.Greater, tick),
	).DisableSort()

	jobs, _, err := s.dbconn.Jobs("", q)
	if err != nil {
		return nil, nil, err
	}

	attacking := []string{}
	defending := []string{}
	for _, j := range jobs {
		if j.SourceFactionID == factionID {
			attacking = append(attacking, j.TargetFactionID)
		} else {
			defending = append(defending, j.SourceFactionID)
		}
	}

	return attacking, defending, nil
}

func determineFactionGoals(ctx *simutil.FactionContext) []structs.Goal {
	goals := []structs.Goal{}

	// if we're a research faction, GoalResearch
	researchTopics := []string{}
	researchProb := []float64{}
	for topic, weight := range ctx.Summary.Research {
		researchTopics = append(researchTopics, topic)
		researchProb = append(researchProb, float64(weight))
	}
	if len(researchTopics) > 0 {
		goals = append(goals, structs.GoalResearch)
	}

	// if we're a religion or highly religious, GoalPiety
	if ctx.Summary.IsReligion || ctx.Summary.Piety > (structs.MaxEthos*9)/10 {
		goals = append(goals, structs.GoalPiety)
	} else if ctx.Summary.Piety > structs.MaxEthos/10 {
		if rnggen.Float64() < toProbability(ctx.Summary.Piety) {
			goals = append(goals, structs.GoalPiety)
		}
	}

	relations := ctx.Relations()

	// decide if we need a military or military intelligence
	if len(relations.Nemesis)+len(relations.Hostile) > 0 {
		goals = append(goals, structs.GoalMilitary)
	} else if len(relations.Hostile)+len(relations.Rival) > 0 {
		goals = append(goals, structs.GoalEspionage)
	} else if ctx.Summary.IsCovert && rnggen.Float64() < toProbability(ctx.Summary.Ethos.Caution) {
		goals = append(goals, structs.GoalEspionage)
	}

	// maybe make friends
	if len(relations.Unfriendly)+len(relations.Neutral) > 0 && (ctx.Summary.Altruism > structs.MaxEthos/4 || ctx.Summary.Pacifism > structs.MaxEthos/4) {
		goals = append(goals, structs.GoalDiplomacy)
	} else if len(relations.Unfriendly) > 0 && rnggen.Float64()*2 < (toProbability(ctx.Summary.Altruism)+toProbability(ctx.Summary.Pacifism)) {
		goals = append(goals, structs.GoalDiplomacy)
	}

	// improve stability
	if ctx.Summary.Corruption > structs.MaxEthos/4 && rnggen.Float64() < toProbability(ctx.Summary.Altruism) {
		goals = append(goals, structs.GoalStability)
	} else if ctx.Summary.Cohesion < structs.MaxEthos/4 && rnggen.Float64() < toProbability(ctx.Summary.Tradition) {
		goals = append(goals, structs.GoalStability)
	} else if rnggen.Float64() < toProbability(ctx.Summary.Ethos.Caution) && ctx.Summary.Cohesion > structs.MaxEthos/5 {
		goals = append(goals, structs.GoalStability)
	}

	// increase our power & reach
	if rnggen.Float64() < toProbability(ctx.Summary.Ethos.Ambition) {
		goals = append(goals, structs.GoalPower)
	} else if rnggen.Float64() < toProbability(ctx.Summary.Ethos.Ambition) {
		goals = append(goals, structs.GoalTerritory)
	}

	// increase our wealth / influence
	if rnggen.Float64() < toProbability(ctx.Summary.Ethos.Altruism) {
		goals = append(goals, structs.GoalWealth)
	} else if rnggen.Float64() < toProbability(ctx.Summary.Ethos.Altruism) {
		goals = append(goals, structs.GoalGrowth)
	}

	// increase our military power / knowledge
	if rnggen.Float64() > toProbability(ctx.Summary.Ethos.Pacifism) {
		goals = append(goals, structs.GoalMilitary)
	} else if rnggen.Float64() > toProbability(ctx.Summary.Ethos.Pacifism) {
		goals = append(goals, structs.GoalEspionage)
	}

	// finally, randomly ensure one of these is included (they're always good background goals)
	rvalue := rnggen.Float64()
	if rvalue < 0.20 {
		goals = append(goals, structs.GoalWealth)
	} else if rvalue < 0.40 {
		goals = append(goals, structs.GoalGrowth)
	} else if rvalue < 0.60 {
		goals = append(goals, structs.GoalTerritory)
	} // 40% nothing added

	return goals
}

func toProbability(x int) float64 {
	return (float64(x) + structs.MaxEthos) / float64(structs.MaxEthos*2)
}

func (s *Base) actionWeightsForFaction(ctx *simutil.FactionContext, land *structs.LandSummary, people, peopleAvail int) (*simutil.ActionWeights, bool, error) {
	survivalConcerns := false
	weights := simutil.NewActionWeights(s.cfg.Actions)

	// allow actions we otherwise can't use if applicable
	if ctx.Summary.IsReligion {
		weights.SetIsReligion()
	}
	if ctx.Summary.IsGovernment {
		weights.SetIsGovernment()
	}

	// nullify actions we can't afford
	wealth := float64(ctx.Summary.Wealth)
	weights.WeightByCost(0.0, wealth)
	weights.WeightByCost(0.5, wealth*0.75) // we don't want to spend all our money on one action

	// nullify actions we don't have the people for
	weights.WeightByMinPeople(0.0, peopleAvail)

	// actions we aglin with we're more likely to do & vice versa
	weights.WeightByEthos(&ctx.Summary.Ethos)

	// weight based on the law & how tradition focused / law abiding the faction is
	tradition := float64(ctx.Summary.Ethos.Tradition/structs.MinEthos) + 1.0 // 2.0 -> 0.0
	if ctx.Summary.IsCovert {
		// Nb. a value of 1.1 means a faction with tradition of X becomes Yx more likely to break the law
		// -10k -> 2.2x
		// -5k -> 1.65x
		// -2.5k -> 1.375x
		// -1k -> 1.21x
		// -500 -> 1.122x
		weights.WeightByIllegal(1.1*tradition, ctx.AllGovernments()...)
	} else {
		// Nb. a value of 0.85 means this value will be > 1 for factions with Tradition < -1.8k (ie. they don't respect traditions)
		// meaning factions that have low respect for traditions become slightly more likely to break the law
		weights.WeightByIllegal(0.85*tradition, ctx.AllGovernments()...)
	}

	// survival: we need money to operate
	valuation := 0.0
	for areaID, arealand := range land.Areas {
		// valuation of the base land
		valuation += s.eco.LandValue(areaID, 0) * float64(arealand.TotalSize)
	}

	desired := valuation / 10
	if !ctx.Summary.IsCovert && ctx.LocalGovernment != nil {
		// if we pay tax (not covert) and there is a local government, try to keep enough money to pay tax
		if ctx.LocalGovernment.TaxRate > 0 && ctx.LocalGovernment.TaxFrequency > 0 {
			desired = math.Max(desired, valuation*ctx.LocalGovernment.TaxRate)
		}
	}

	if wealth < desired { // we believe we're at risk of running out of money
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight, structs.GoalWealth)
		survivalConcerns = true
	} else if wealth < desired*1.2 { // we're starting to run low on money
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight/2, structs.GoalWealth)
	}

	// survival: we need people to work for us
	if people < s.cfg.Settings.SurvivalMinPeople {
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight, structs.GoalGrowth)
		survivalConcerns = true
	} else if people < (s.cfg.Settings.SurvivalMinPeople*6)/5 { // 1.2x
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight/2, structs.GoalGrowth)
	}

	// survival: we need to keep cohesion up
	if ctx.Summary.Cohesion < structs.MaxEthos/10 {
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight, structs.GoalStability)
		survivalConcerns = true
	} else if ctx.Summary.Cohesion < structs.MaxEthos/5 {
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight/2, structs.GoalStability)
	}

	// survival: we need to keep corruption down
	if ctx.Summary.Corruption > structs.MaxEthos-(structs.MaxEthos/10) {
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight, structs.GoalStability)
		survivalConcerns = true
	} else if ctx.Summary.Corruption > structs.MaxEthos-(structs.MaxEthos/5) {
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight/2, structs.GoalStability)
	}

	// survival: we need to have place(s) of work
	if land.Count < 2 && ctx.Summary.IsCovert || land.Count < 1 {
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight, structs.GoalTerritory)
	}

	return weights, survivalConcerns, nil
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
