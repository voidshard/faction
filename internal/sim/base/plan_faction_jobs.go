package base

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
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

const (
	estimateVassalPeopleDegradeByStep = 0.4 // estimate for people of vassal factions who might work for us
	maxWarTargetSizeDifference        = 1.3 // we wont make war on factions too much larger than us
	estimateParentWarFactor           = 0.6 // we'll add this % of the parent faction size to the child faction size
	warArroganceFactor                = 0.3 // if an opponent(s) are less that this % of our membership we assume victory
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

	ctx, _, availablePeople, err := s.prepFactionData(factionID)
	if err != nil {
		return nil, err
	}

	var plans []*structs.Job
	if len(attacking) > 0 || len(defending) > 0 {
		// if we're at war, not dying consumes most of our concerns
		plans, err = s.planFactionJobsWartime(tick, ctx, attacking, defending, availablePeople)
	} else {
		// if we're at peace, we have lots of competing concerns
		plans, err = s.planFactionJobsPeacetime(tick, ctx, availablePeople)
	}
	if err != nil || plans == nil || len(plans) == 0 {
		// whatever happened, we're done
		return nil, err
	}

	// final checks & event creation
	jobs := []*structs.Job{}
	events := []*structs.Event{}
	for _, j := range plans {
		cfg, ok := s.cfg.Actions[j.Action]
		if !ok {
			// we have no config for the Job, we can't do it (how did it even get here?!)
			continue
		}

		extraTarget, ok := structs.ActionTarget[j.Action]
		if ok && j.TargetMetaKey != extraTarget {
			// we need extra target info, but it's not set (probably we couldn't find a target)
			continue
		}

		ctx.Summary.Wealth -= int(cfg.Cost.Min)
		events = append(events, simutil.NewJobPendingEvent(j))
		jobs = append(jobs, j)
	}

	// push everything into the DB
	err = s.dbconn.SetFactions(ctx.Summary.ToFaction()) // updated Wealth & Members values
	if err != nil {
		return jobs, err
	}
	err = s.dbconn.SetJobs(jobs...)
	if err != nil {
		return jobs, err
	}
	return jobs, s.dbconn.SetEvents(events...)
}

func (s *Base) prepFactionData(factionID string) (*simutil.FactionContext, []*structs.Faction, int, error) {
	// reads most data relevant to the faction
	ctx, err := simutil.NewFactionContext(s.dbconn, factionID)
	if err != nil {
		return nil, nil, 0, err
	}

	// look up all child factions that might lend us people
	children, _, err := s.dbconn.Factions(
		"",
		db.Q(
			db.F(db.ParentFactionID, db.Equal, factionID),
			db.F(db.ParentFactionRelation, db.Greater, structs.FactionRelationPuppet),
		),
	)
	if err != nil {
		return nil, nil, 0, err
	}

	// people we have available for work
	vassals := estimateVassalPeople(factionID, children)
	people := estimateAvailablePeople(ctx.Summary.Ranks) + vassals/2

	ctx.Summary.Members = ctx.Summary.Ranks.Total // nb. this is people with ranks, not available people for work
	ctx.Summary.Vassals = vassals                 // estimation of the number of people from child factions we could use
	ctx.Summary.Plots = ctx.Land.Count
	ctx.Summary.Areas = len(ctx.Areas)

	return ctx, children, people, nil
}

func (s *Base) planFactionJobsPeacetime(tick int, ctx *simutil.FactionContext, availablePeople int) ([]*structs.Job, error) {
	availableWealth := float64(ctx.Summary.Wealth)

	// Get all jobs for this faction that are active & wont finish next tick
	q := db.Q(
		db.F(db.SourceFactionID, db.Equal, ctx.Summary.ID),
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

	// how likely we are to do various actions
	weights, survivalConcerns, err := s.actionWeightsForFaction(ctx, availablePeople)
	if err != nil {
		return nil, err
	}

	if !survivalConcerns {
		// we'll only apply higher level goals if imminent death isn't a concern
		goals := determineFactionGoals(ctx)
		weights.WeightByGoal(s.cfg.Settings.GoalWeight, goals...)
	}

	// hirable services in our area(s) (we'll fill out if we need)
	services := map[structs.ActionType][]string{}

	// decide which actions we should do
	if weights.Normalise() <= 0 {
		// there are no actions we wish to do / can do
		return nil, nil
	}

	// pick actions to do, if we can't find valid targets we'll move on
	limitCount := 100
	plans := []*structs.Job{}
	for {
		// we should remove jobs as we decide we can't do them, but this adds a hard limit
		limitCount -= 1
		if limitCount <= 0 {
			break
		}

		// pick out an action
		act := weights.Choose()
		cfg, ok := s.cfg.Actions[act]
		if !ok {
			continue
		}

		target, err := s.chooseJobTargetFaction(ctx, inflight, act, cfg)
		if err != nil {
			return nil, err
		}
		if target == "" {
			weights.WeightAction(0.0, act) // we can't do this action, remove it from the list
			if weights.Normalise() <= 0 {
				break // no more choices :(
			}
			continue
		}

		job := simutil.NewJob(tick, act, cfg)
		job.SourceFactionID = ctx.Summary.ID
		job.SourceAreaID = ctx.RandomArea(act)
		job.TargetFactionID = target
		job.Conscription = simutil.SourceMeetsActionConditions(ctx.Summary, cfg.Conscription)
		job.IsIllegal = ctx.Summary.IsCovert || simutil.IsIllegalAction(act, ctx.AllGovernments()...)
		if job.IsIllegal {
			job.Secrecy = int(float64(rnggen.Intn(structs.MaxTuple)+ctx.Summary.EspionageDefense) * cfg.SecrecyWeight)
		}

		if job.Action == structs.ActionTypeHireMercenaries || job.Action == structs.ActionTypeHireSpies {
			if services == nil { // services cached at this level so we don't re-fetch data
				services, err = simutil.ServicesForHire(s.cfg.Actions, s.dbconn, ctx.AreaIDs())
				if err != nil {
					return nil, err
				}
			}

			hirelingJob, err := s.buildMercenaryJob(services, job)
			if err != nil {
				return nil, err
			}
			if hirelingJob == nil {
				weights.WeightAction(0.0, act) // we can't do this action, remove it from the list
				if weights.Normalise() <= 0 {
					break // no more choices :(
				}
				continue
			}

			hirelingJob.IsIllegal = ctx.Summary.IsCovert || simutil.IsIllegalAction(hirelingJob.Action, ctx.AllGovernments()...)
			if hirelingJob.IsIllegal {
				hirelingJob.Secrecy = int(float64(rnggen.Intn(structs.MaxTuple)+ctx.Summary.EspionageDefense) * cfg.SecrecyWeight)
			}
			job.TargetAreaID = ctx.Summary.HomeAreaID // we sit at home while someone else works
			plans = append(plans, hirelingJob)
		}

		inflight[jobKey(job.TargetFactionID, act)] = true

		// reduce the number of available people / money by the actions min
		availablePeople -= cfg.MinPeople
		availableWealth -= cfg.Cost.Min
		if availablePeople <= 0 {
			break
		} else {
			weights.WeightByMinPeople(0.0, availablePeople)
			weights.WeightByCost(0.0, availableWealth)
			if weights.Normalise() <= 0 {
				break
			}
		}

		plans = append(plans, job)
	}

	err = s.setSpecificJobTargets(ctx, plans)
	if err != nil {
		return nil, err
	}

	return plans, err
}

func (s *Base) setSpecificJobTargetsPerson(jobs []*structs.Job) error {
	factionIDs := []string{}
	for _, j := range jobs {
		factionIDs = append(factionIDs, j.TargetFactionID)
	}

	// pull out top 10 people of each faction
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

		personObj, ok := targets.People[chosen]
		if ok {
			j.TargetAreaID = personObj.AreaID
		}

		j.TargetMetaKey = structs.MetaKeyPerson
		j.TargetMetaVal = chosen
	}

	return nil
}

func (s *Base) buildMercenaryJob(services map[structs.ActionType][]string, job *structs.Job) (*structs.Job, error) {
	// choose who we will hire and for what action
	var vendorAction structs.ActionType
	var vendorID string

	for act, actcfg := range s.cfg.Actions {
		if job.Action == structs.ActionTypeHireMercenaries && !actcfg.ValidServiceMercenary {
			continue
		} else if job.Action == structs.ActionTypeHireSpies && !actcfg.ValidServiceSpy {
			continue
		}

		factionsThatOfferService, ok := services[act]
		if !ok {
			continue
		}
		for _, fid := range factionsThatOfferService {
			if fid != job.TargetFactionID { // don't hire target to attack itself
				vendorID = fid
				vendorAction = act
				break
			}
		}
		if vendorID != "" {
			break
		}
	}
	if vendorID == "" {
		return nil, nil
	}

	// build a new child job
	cfg, _ := s.cfg.Actions[vendorAction] // we checked this is valid above
	hirejob := simutil.NewJob(job.TickCreated, vendorAction, cfg)
	hirejob.ParentJobID = job.ID
	hirejob.SourceFactionID = vendorID
	hirejob.TargetFactionID = job.TargetFactionID

	// set the meta information on the parent job
	job.TargetMetaKey = structs.MetaKeyJob
	job.TargetMetaVal = hirejob.ID

	return hirejob, nil
}

func (s *Base) setSpecificJobTargetsPlot(jobs []*structs.Job) error {
	factionIDs := []string{}

	for _, j := range jobs {
		// targeting lands of the target
		factionIDs = append(factionIDs, j.TargetFactionID)
	}

	// pull out the most expensive plots
	plots, err := s.dbconn.FactionPlots(20, factionIDs...)
	if err != nil {
		return err
	}

	// read out what we know about the targets
	intel, err := s.dbconn.TuplesSumModsBySubject(
		db.RelationFactionFactionIntelligence,
		jobs[0].SourceFactionID,
		factionIDs...,
	)
	if err != nil {
		return err
	}

	for _, j := range jobs {
		intelligence, _ := intel[j.TargetFactionID]

		fplots, ok := plots[j.TargetFactionID]
		if !ok || len(fplots) == 0 {
			continue
		}

		potentialTargets := []*structs.Plot{}
		for _, p := range fplots {
			// plot isn't hidden or we know it belongs to the target
			if p.Hidden <= 0 || intelligence > p.Hidden {
				potentialTargets = append(potentialTargets, p)
			}
		}

		if len(potentialTargets) == 0 {
			continue
		}

		target := potentialTargets[rnggen.Intn(len(potentialTargets))]
		j.TargetAreaID = target.AreaID
		j.TargetMetaKey = structs.MetaKeyPlot
		j.TargetMetaVal = target.ID
	}

	return nil
}

func (s *Base) setSpecificJobTargetsArea(ctx *simutil.FactionContext, jobs []*structs.Job) error {
	needAreas := []string{}
	for _, j := range jobs {

		if j.SourceFactionID != j.TargetFactionID {
			// we need areas these factions are based in
			needAreas = append(needAreas, j.TargetFactionID)
		}
	}

	factionToArea, err := s.dbconn.FactionAreas(false, needAreas...)
	if err != nil {
		return err
	}

	for _, j := range jobs {
		if j.SourceFactionID == j.TargetFactionID {
			j.TargetAreaID = ctx.RandomArea(j.Action)
			if j.Action == structs.ActionTypeHarvest {
				// Pick some area in which we can harvest resources.
				for areaID, summary := range ctx.Land.Areas { // because order of iteration of maps is undefined
					if len(summary.Commodities) > 0 {
						j.TargetAreaID = areaID
						break
					}
				}
			}
		} else {
			fareas, ok := factionToArea[j.TargetFactionID]
			if !ok {
				continue
			}
			for areaID := range fareas { // because order of iteration of maps is undefined
				j.TargetAreaID = areaID
				break
			}
		}
	}

	return nil
}

func (s *Base) setSpecificJobTargets(ctx *simutil.FactionContext, plans []*structs.Job) error {
	// bucket jobs based on what kind of target they need
	needed := map[structs.MetaKey][]*structs.Job{}
	for _, j := range plans {
		// a few actions need special handling
		if j.Action == structs.ActionTypeResearch {
			j.TargetAreaID = ctx.RandomArea(j.Action)
			j.TargetMetaKey = structs.MetaKeyResearch
			j.TargetMetaVal = ctx.RandomResearch()
			continue
		}

		metakey, ok := structs.ActionTarget[j.Action]
		if !ok {
			// Area -> everything needs this if nothing else
			metakey = structs.MetaKeyArea
			if j.SourceFactionID == j.TargetFactionID {
				j.SourceAreaID = ctx.RandomArea(j.Action)
				j.TargetAreaID = j.SourceAreaID
				continue
			}
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
			err := s.setSpecificJobTargetsPlot(jobs)
			if err != nil {
				return err
			}
		case structs.MetaKeyPerson:
			err := s.setSpecificJobTargetsPerson(jobs)
			if err != nil {
				return err
			}
		case structs.MetaKeyJob:
			// only HireMercenaries & HireSpies need this and it's Area / Job targets are set above.
			// Any child job(s) they spawn target Plots / Areas / People the same
			// as everything else.
		default:
			err := s.setSpecificJobTargetsArea(ctx, jobs)
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
func (s *Base) chooseJobTargetFaction(ctx *simutil.FactionContext, inflight map[string]bool, act structs.ActionType, cfg *config.Action) (string, error) {
	relations := ctx.Relations()

	if cfg.TargetMinTrust == 0 && cfg.TargetMaxTrust == 0 {
		_, ok := inflight[jobKey(ctx.Summary.ID, act)]
		if ok { // this job & target is already inflight
			return "", nil
		}
		return ctx.Summary.ID, nil // target ourselves
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
		return "", nil
	}

	// before starting a war we have extra things to consider
	if act == structs.ActionTypeWar || act == structs.ActionTypeCrusade || act == structs.ActionTypeShadowWar {
		factions, _, err := s.dbconn.Factions("", db.Q(db.F(db.ID, db.In, choices)).DisableSort())
		if err != nil {
			return "", err
		}
		// if they have parent factions, we need to grab those too
		parentIds := []string{}
		for _, f := range factions {
			if f.ParentFactionID != "" {
				parentIds = append(parentIds, f.ParentFactionID)
			}
		}
		parents := map[string]*structs.Faction{}
		if len(parentIds) > 0 {
			parentFactions, _, err := s.dbconn.Factions("", db.Q(db.F(db.ID, db.In, parentIds)).DisableSort())
			if err != nil {
				return "", err
			}
			for _, f := range parentFactions {
				parents[f.ID] = f
			}
		}
		// work out how much larger a faction (might) be compared to us
		me := float64(ctx.Summary.Members+ctx.Summary.Vassals) * maxWarTargetSizeDifference
		choices = []string{}
		for _, f := range factions {
			size := float64(f.Members + f.Vassals)
			pf, ok := parents[f.ParentFactionID] // nb. "" here ('no parent') will yield !ok
			if ok {
				size += (float64(pf.Members+pf.Vassals) - float64(f.Members)*estimateVassalPeopleDegradeByStep) * estimateParentWarFactor
			}
			if size <= me {
				choices = append(choices, f.ID)
			}
		}
	}

	// pick the first valid target that for which we don't already have a job
	for _, targetFactionID := range choices {
		_, ok := inflight[jobKey(targetFactionID, act)]
		if ok { // this job & target is already inflight
			continue
		}
		return targetFactionID, nil
	}

	return "", nil
}

func (s *Base) planFactionJobsWartime(tick int, ctx *simutil.FactionContext, attacking, defending []*structs.Job, people int) ([]*structs.Job, error) {
	// get factions we're at war with
	atkSet := mapset.NewSet[string]()
	defSet := mapset.NewSet[string]()
	factionIDs := []string{}
	for _, j := range attacking {
		factionIDs = append(factionIDs, j.TargetFactionID)
		atkSet.Add(j.TargetFactionID)
	}
	needsMoreSoldiers := false
	for _, j := range defending {
		factionIDs = append(factionIDs, j.SourceFactionID)
		defSet.Add(j.SourceFactionID)
		needsMoreSoldiers = needsMoreSoldiers || j.State == structs.JobStatePending
	}
	if needsMoreSoldiers {
		// we're already fighting a war we don't have enough people signed up to fight -> no need to queue any more jobs
		return nil, nil
	}

	counterAtk := defSet.Difference(atkSet).ToSlice() // people we're defending against that we're not attacking
	if len(counterAtk) > 0 && len(defending) > 0 {
		act := structs.ActionTypeWar
		if ctx.Summary.IsCovert {
			act = structs.ActionTypeShadowWar
		} else if ctx.Summary.IsReligion {
			act = structs.ActionTypeCrusade
		}

		cfg, ok := s.cfg.Actions[act]
		if !ok {
			act = defending[0].Action
			cfg, ok = s.cfg.Actions[act] // this should certainly have a config ..!
			if !ok {
				return nil, fmt.Errorf("no config for action %s", act)
			}
		}

		j := simutil.NewJob(tick, act, cfg)
		j.SourceFactionID = ctx.Summary.ID
		j.SourceAreaID = ctx.RandomArea(j.Action)
		j.TargetFactionID = counterAtk[0]

		plans := []*structs.Job{j}
		err := s.setSpecificJobTargets(ctx, plans)
		return plans, err
	}

	// pull out the factions we're fighting with
	factions, _, err := s.dbconn.Factions("", db.Q(db.F(db.ID, db.In, factionIDs)).DisableSort())
	if err != nil {
		return nil, err
	}

	// sum up how large the enemy faction(s) are
	enemyScale := 0.0
	for _, f := range factions {
		enemyScale += float64(f.Members + f.Vassals)
	}

	// check how worried we are (if we're not worried, we'll queue jobs like we're not at war)
	// high caution yields closer to 50%, low caution closer to 30% (the base)
	arrogance := warArroganceFactor + (toProbability(ctx.Summary.Caution) * 0.2)
	if enemyScale <= arrogance*float64(ctx.Summary.Members+ctx.Summary.Vassals) {
		return s.planFactionJobsPeacetime(tick, ctx, people) // HA IS THAT ALL YOU GOT
	}

	// nb. we hope in doing this that more people join the fight, alternatively, we hope to save money for the war
	return nil, nil
}

func (s *Base) determineFactionAtWar(tick int, factionID string) ([]*structs.Job, []*structs.Job, error) {
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

	attacking := []*structs.Job{}
	defending := []*structs.Job{}
	for _, j := range jobs {
		if j.SourceFactionID == factionID {
			attacking = append(attacking, j)
		} else {
			defending = append(defending, j)
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

func (s *Base) actionWeightsForFaction(ctx *simutil.FactionContext, peopleAvail int) (*simutil.ActionWeights, bool, error) {
	survivalConcerns := false
	weights := simutil.NewActionWeights(s.cfg.Actions)

	// forbid actions our faction isn't allowed to consider
	weights.ApplyActionConditions(ctx.Summary)

	// we have nothing to harvest :(
	if len(ctx.Land.Commodities) == 0 {
		weights.WeightAction(0.0, structs.ActionTypeHarvest)
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
	for areaID, arealand := range ctx.Land.Areas {
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
	if ctx.Summary.Ranks.Total < s.cfg.Settings.SurvivalMinPeople {
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight, structs.GoalGrowth)
		survivalConcerns = true
	} else if ctx.Summary.Ranks.Total < (s.cfg.Settings.SurvivalMinPeople*6)/5 { // 1.2x
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
	if ctx.Land.Count < 2 && ctx.Summary.IsCovert || ctx.Land.Count < 1 {
		weights.WeightByGoal(s.cfg.Settings.SurvivalGoalWeight, structs.GoalTerritory)
	}

	return weights, survivalConcerns, nil
}

func estimateVassalPeople(id string, vassals []*structs.Faction) int {
	total := 0
	for _, v := range vassals {
		total += int(float64(v.Members) * estimateVassalPeopleDegradeByStep)
		total += int(float64(v.Vassals) * estimateVassalPeopleDegradeByStep * estimateVassalPeopleDegradeByStep)
	}
	return total
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
