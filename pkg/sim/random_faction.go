/*
random_faction.go - random faction / government generation
*/
package sim

import (
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/stats"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// factionRand is a helper struct to generate random factions
// using dice / distributions provided in configs
type factionRand struct {
	yieldRand      *yieldRand
	cfg            *config.Faction
	rng            *rand.Rand
	ethosAltruism  *stats.Rand
	ethosAmbition  *stats.Rand
	ethosTradition *stats.Rand
	ethosPacifism  *stats.Rand
	ethosPiety     *stats.Rand
	ethosCaution   *stats.Rand
	leaderOccur    stats.Normalised
	leaderList     []structs.LeaderType
	structOccur    stats.Normalised
	structList     []structs.LeaderStructure
	wealth         *stats.Rand
	cohesion       *stats.Rand
	corruption     *stats.Rand
	espOffense     *stats.Rand
	espDefense     *stats.Rand
	milOffense     *stats.Rand
	milDefense     *stats.Rand
	focusOccur     stats.Normalised
	focusCount     stats.Normalised
	focusWeights   []*stats.Rand
	guildOccur     stats.Normalised
	guildCount     stats.Normalised
	propertyCount  stats.Normalised
	areas          []string
	plotSize       *stats.Rand
}

// faction + associated metadata used during creation
type metaFaction struct {
	faction       *structs.Faction
	actions       []structs.ActionType
	actionWeights []*structs.Tuple
	land          []*structs.LandRight
	plots         []*structs.Plot
	profWeights   []*structs.Tuple
}

func (fr *factionRand) randLandForGuild(g *config.Guild) []*structs.LandRight {
	options, ok := fr.yieldRand.professionLand[g.Profession]
	if !ok {
		return nil
	}

	index := 0
	total := 0
	for _, l := range options {
		index++
		total += l.Yield
		if total >= g.MinYield {
			break
		}
	}

	if total < g.MinYield {
		return nil
	}

	defer fr.recalcGuildProb()

	if index >= len(options) {
		// we've assigned all available land
		delete(fr.yieldRand.professionLand, g.Profession)
		delete(fr.yieldRand.professionYield, g.Profession)
		return options
	}

	// we've assigned some subset of available land
	fr.yieldRand.professionLand[g.Profession] = options[index+1:]
	ftotal := fr.yieldRand.professionYield[g.Profession]
	fr.yieldRand.professionYield[g.Profession] = ftotal - total

	return options[:index+1]
}

func (fr *factionRand) recalcGuildProb() {
	guildProb := []float64{}
	if fr.yieldRand != nil {
		for _, guild := range fr.cfg.Guilds {
			total, _ := fr.yieldRand.professionYield[guild.Profession]
			prob := 0.0
			if total >= guild.MinYield {
				prob = guild.Probability
			}
			guildProb = append(guildProb, prob)
		}
	}
	fr.guildOccur = stats.NewNormalised(guildProb)
}

func (s *simulationImpl) SpawnFactions(count int, cfg *config.Faction, areas ...string) ([]*structs.Faction, error) {
	// prep some filters
	arf := []*db.AreaFilter{}
	lrf := []*db.LandRightFilter{}
	for _, area := range areas {
		arf = append(arf, &db.AreaFilter{ID: area})
		lrf = append(lrf, &db.LandRightFilter{AreaID: area})
	}

	// what land is good for what / who
	yields, err := s.areaYields(lrf, false)
	if err != nil {
		return nil, err
	}

	// dice based on faction cfg + land yields
	dice := newFactionRand(cfg, yields, areas)

	// lookup which areas are run by which governments (determines faction legality / covert status)
	areaToGovt, govtToGovt, err := s.areaGovernments(arf)
	if err != nil {
		return nil, err
	}

	// generate factions
	factions := []*structs.Faction{}
	govsToWrite := map[string]*structs.Government{}
	for i := 0; i < count; i++ {
		f := s.randFaction(dice)
		govtID, _ := areaToGovt[f.faction.HomeAreaID]
		govt, _ := govtToGovt[govtID]

		factionOutlawed := false
		for _, act := range f.actions {
			illegal, _ := govt.Outlawed.Actions[act]
			if illegal {
				factionOutlawed = true
				break
			}
		}

		f.faction.GovernmentID = govtID
		f.faction.IsCovert = factionOutlawed

		if factionOutlawed {
			govt.Outlawed.Factions[f.faction.ID] = true
			govsToWrite[govtID] = govt
		}

		err = writeMetaFaction(s.dbconn, f)
		if err != nil {
			return nil, err
		}

		factions = append(factions, f.faction)
	}

	// finally, flush law change(s) to govt.
	govs := []*structs.Government{}
	for _, govt := range govsToWrite {
		govs = append(govs, govt)
	}
	if len(govs) > 0 {
		err = s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
			return tx.SetGovernments(govs...)
		})
	}

	return factions, err
}

func writeMetaFaction(conn *db.FactionDB, f *metaFaction) error {
	return conn.InTransaction(func(tx db.ReaderWriter) error {
		err := tx.SetFactions(f.faction)
		if err != nil {
			return err
		}
		err = tx.SetLandRights(f.land...)
		if err != nil {
			return err
		}
		err = tx.SetPlots(f.plots...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationFactionActionTypeWeight, f.actionWeights...)
		if err != nil {
			return err
		}
		return tx.SetTuples(db.RelationFactionProfessionWeight, f.profWeights...)
	})
}

// randFaction spits out a random faction.
func (s *simulationImpl) randFaction(fr *factionRand) *metaFaction {
	// start with a lot of randomly inserted fields
	mf := &metaFaction{
		faction: &structs.Faction{
			Ethos: structs.Ethos{ // random ethos to start with
				Altruism:  fr.ethosAltruism.Int(),
				Ambition:  fr.ethosAmbition.Int(),
				Tradition: fr.ethosTradition.Int(),
				Pacifism:  fr.ethosPacifism.Int(),
				Piety:     fr.ethosPiety.Int(),
				Caution:   fr.ethosCaution.Int(),
			},
			ID:               structs.NewID(),
			Leadership:       fr.leaderList[fr.leaderOccur.Int()],
			Structure:        fr.structList[fr.structOccur.Int()],
			Wealth:           fr.wealth.Int(),
			Cohesion:         fr.cohesion.Int(),
			Corruption:       fr.corruption.Int(),
			EspionageOffense: fr.espOffense.Int(),
			EspionageDefense: fr.espDefense.Int(),
			MilitaryOffense:  fr.milOffense.Int(),
			MilitaryDefense:  fr.milDefense.Int(),
		},
		actions:       []structs.ActionType{},
		actionWeights: []*structs.Tuple{},
		land:          []*structs.LandRight{},
		plots:         []*structs.Plot{},
		profWeights:   []*structs.Tuple{},
	}

	// consider action focuses
	professions := []string{}
	actions := []structs.ActionType{}
	seen := map[int]bool{}
	actionsEthos := []*structs.Ethos{&mf.faction.Ethos}
	for i := 0; i < fr.focusCount.Int(); i++ {
		choice := fr.focusOccur.Int()
		_, ok := seen[choice]
		if ok {
			continue
		}
		seen[choice] = true

		focus := fr.cfg.Focuses[choice]
		weight := fr.focusWeights[choice]

		mf.faction.EspionageOffense += int(focus.EspionageOffenseBonus * float64(mf.faction.EspionageOffense))
		mf.faction.EspionageDefense += int(focus.EspionageDefenseBonus * float64(mf.faction.EspionageDefense))
		mf.faction.MilitaryOffense += int(focus.MilitaryOffenseBonus * float64(mf.faction.MilitaryOffense))
		mf.faction.MilitaryDefense += int(focus.MilitaryDefenseBonus * float64(mf.faction.MilitaryDefense))

		for _, act := range focus.Actions {
			mf.actionWeights = append(mf.actionWeights, &structs.Tuple{
				Subject: mf.faction.ID,
				Object:  string(act),
				Value:   weight.Int(),
			})

			actionCfg, ok := s.cfg.Actions[act]
			if !ok {
				continue
			}

			actions = append(actions, act)
			actionsEthos = append(actionsEthos, &actionCfg.Ethos)
			for prof := range actionCfg.ProfessionWeights {
				professions = append(professions, prof)
			}
		}
	}

	// favoured actions influence faction's starting ethos
	mf.faction.Ethos = *structs.EthosAverage(actionsEthos...)

	// consider a profession based guild
	seen = map[int]bool{}
	countHarvest := 0
	countCraft := 0
	for i := 0; i < fr.guildCount.Int(); i++ {
		choice := fr.guildOccur.Int()
		if choice < 0 {
			break // there is no choice
		}

		_, ok := seen[choice]
		if ok {
			continue
		}

		guild := fr.cfg.Guilds[choice]
		landrights := fr.randLandForGuild(&guild)

		if landrights == nil {
			continue
		}

		for _, land := range landrights {
			land.FactionID = mf.faction.ID

			if s.eco.IsCraftable(land.Commodity) {
				countCraft++
			}
			if s.eco.IsHarvestable(land.Commodity) {
				countHarvest++
			}
		}
		mf.land = append(mf.land, landrights...)
		professions = append(professions, guild.Profession)
	}

	// now we can determine profession weights
	professionCounts := map[string]int{}
	for _, p := range professions {
		professionCounts[p]++
	}
	for p, c := range professionCounts {
		w := c * structs.MaxEthos / 10
		if w > structs.MaxEthos {
			w = structs.MaxEthos
		}
		mf.profWeights = append(mf.profWeights, &structs.Tuple{
			Subject: mf.faction.ID,
			Object:  p,
			Value:   w,
		})
	}

	// and action weights...
	actionCounts := map[structs.ActionType]int{}
	for _, a := range actions {
		actionCounts[a]++
	}
	for a, c := range actionCounts {
		w := c * structs.MaxEthos / 10
		if w > structs.MaxEthos {
			w = structs.MaxEthos
		}
		mf.actions = append(mf.actions, a)
		mf.actionWeights = append(mf.actionWeights, &structs.Tuple{
			Subject: mf.faction.ID,
			Object:  string(a),
			Value:   w,
		})
	}

	// give faction land if we're still below the min
	for i := len(mf.land); i < fr.propertyCount.Int()+1; i++ {
		mf.plots = append(mf.plots, &structs.Plot{
			ID:        structs.NewID(),
			AreaID:    fr.areas[fr.rng.Intn(len(fr.areas))],
			FactionID: mf.faction.ID,
			Size:      fr.plotSize.Int(),
		})
	}

	// pick a headquarters
	switch len(mf.plots) {
	case 0:
		mf.faction.HomeAreaID = mf.land[fr.rng.Intn(len(mf.land))].AreaID
	case 1:
		mf.faction.HomeAreaID = mf.plots[0].AreaID
	default:
		mf.faction.HomeAreaID = mf.plots[fr.rng.Intn(len(mf.plots))].AreaID
	}

	return mf
}

// newFactionRand creates a new dice roller for creationg factions based on faction config
// and the available land rights in some area(s).
func newFactionRand(f *config.Faction, yields *yieldRand, areas []string) *factionRand {
	focusOccurProb := []float64{}
	focusWeights := []*stats.Rand{}
	for _, focus := range f.Focuses {
		focusOccurProb = append(focusOccurProb, focus.Probability)
		focusWeights = append(focusWeights, stats.NewRand(focus.Weight.Min, focus.Weight.Max, focus.Weight.Mean, focus.Weight.Deviation))
	}

	leaderProb := []float64{}
	llist := []structs.LeaderType{}
	for leader, prob := range f.LeadershipProbability {
		leaderProb = append(leaderProb, prob)
		llist = append(llist, leader)
	}

	structureProb := []float64{}
	slist := []structs.LeaderStructure{}
	for structure, prob := range f.LeadershipStructureProbability {
		structureProb = append(structureProb, prob)
		slist = append(slist, structure)
	}

	fr := &factionRand{
		yieldRand:      yields,
		cfg:            f,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
		ethosAltruism:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Altruism), float64(f.EthosDeviation.Altruism)),
		ethosAmbition:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Ambition), float64(f.EthosDeviation.Ambition)),
		ethosTradition: stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Tradition), float64(f.EthosDeviation.Tradition)),
		ethosPacifism:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Pacifism), float64(f.EthosDeviation.Pacifism)),
		ethosPiety:     stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Piety), float64(f.EthosDeviation.Piety)),
		ethosCaution:   stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Caution), float64(f.EthosDeviation.Caution)),
		leaderOccur:    stats.NewNormalised(leaderProb),
		leaderList:     llist,
		structOccur:    stats.NewNormalised(structureProb),
		structList:     slist,
		wealth:         stats.NewRand(f.Wealth.Min, f.Wealth.Max, f.Wealth.Mean, f.Wealth.Deviation),
		cohesion:       stats.NewRand(f.Cohesion.Min, f.Cohesion.Max, f.Cohesion.Mean, f.Cohesion.Deviation),
		corruption:     stats.NewRand(f.Corruption.Min, f.Corruption.Max, f.Corruption.Mean, f.Corruption.Deviation),
		espOffense:     stats.NewRand(f.EspionageOffense.Min, f.EspionageOffense.Max, f.EspionageOffense.Mean, f.EspionageOffense.Deviation),
		espDefense:     stats.NewRand(f.EspionageDefense.Min, f.EspionageDefense.Max, f.EspionageDefense.Mean, f.EspionageDefense.Deviation),
		milOffense:     stats.NewRand(f.MilitaryOffense.Min, f.MilitaryOffense.Max, f.MilitaryOffense.Mean, f.MilitaryOffense.Deviation),
		milDefense:     stats.NewRand(f.MilitaryDefense.Min, f.MilitaryDefense.Max, f.MilitaryDefense.Mean, f.MilitaryDefense.Deviation),
		focusOccur:     stats.NewNormalised(focusOccurProb),
		focusCount:     stats.NewNormalised(f.FocusProbability),
		focusWeights:   focusWeights,
		guildCount:     stats.NewNormalised(f.GuildProbability),
		propertyCount:  stats.NewNormalised(f.PropertyProbability),
		areas:          areas,
		plotSize:       stats.NewRand(f.PlotSize.Min, f.PlotSize.Max, f.PlotSize.Mean, f.PlotSize.Deviation),
	}

	fr.recalcGuildProb()
	return fr
}
