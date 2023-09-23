/*
random_faction.go - random faction / government generation
*/
package base

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/internal/sim/simutil"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

func (s *Base) SpawnFactions(count int, cfg *config.Faction, areas ...string) ([]*structs.Faction, error) {
	tick, err := s.dbconn.Tick()
	if err != nil {
		return nil, err
	}

	// prep some filters
	arf := db.Q(db.F(db.ID, db.In, areas))
	lrf := db.Q(
		db.F(db.AreaID, db.In, areas),
		db.F(db.Commodity, db.NotEqual, ""),
	)

	// what land is good for what / who
	yields, err := s.areaYields(lrf, false)
	if err != nil {
		return nil, err
	}

	// dice based on faction cfg + land yields
	dice := newFactionRand(cfg, s.tech, yields, areas)

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

		govtID, _ := areaToGovt[f.Faction.HomeAreaID]
		govt, _ := govtToGovt[govtID]

		if govt != nil {
			f.Faction.GovernmentID = govtID
			if isFactionIllegal(govt, f) {
				govt.Outlawed.Factions[f.Faction.ID] = true
				govsToWrite[govtID] = govt
				makeFactionCovert(dice, f)
			}

		}

		f.Events = append(f.Events, simutil.NewFactionCreatedEvent(f.Faction, tick))
		err = simutil.WriteMetaFaction(s.dbconn, f)
		if err != nil {
			return nil, err
		}

		factions = append(factions, f.Faction)
	}

	// finally, flush law change(s) to govt (if needed)
	if len(govsToWrite) > 0 {
		govs := []*structs.Government{}
		for _, govt := range govsToWrite {
			govs = append(govs, govt)
		}
		err = s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
			return tx.SetGovernments(govs...)
		})
	}

	return factions, err
}

func determineBaseTrust(f1, f2 *simutil.MetaFaction) int {
	// In general we give larger negatives than positives. Trust is much harder to build than suspicion.

	// base value starts from Ethos distance
	// y = -0.4x + 0.2 where x is the distance between the two factions' ethos (0-1)
	// Yields
	// ->  1/5 where x = 0
	// -> -1/5 where x = 1
	// ->  0   where x = 0.5
	dist := structs.EthosDistance(&f1.Faction.Ethos, &f2.Faction.Ethos)
	base := int(float64(structs.MaxEthos)*(-0.4*dist+0.2)) - structs.MaxEthos/20

	// apply some basic "are we alike" modifiers
	if f1.Faction.IsGovernment && f2.Faction.IsCovert {
		base -= structs.MaxEthos / 5
	}
	if f1.Faction.IsCovert && f2.Faction.IsCovert {
		base += structs.MaxEthos / 20
	}
	if f1.Faction.ReligionID != "" {
		if f1.Faction.ReligionID == f2.Faction.ReligionID {
			base += structs.MaxEthos / 10
		} else {
			if f1.Faction.IsReligion {
				base -= structs.MaxEthos / 5
			} else {
				base -= structs.MaxEthos / 10
			}
		}
	}

	// more specific modifiers for favoured actions / resources
	research := map[string]bool{}
	for _, topic := range f1.ResearchWeights { // my research topics
		research[topic.Object] = true
	}
	sharedResearch := false
	for _, topic := range f2.ResearchWeights { // their research topics
		_, ok := research[topic.Object]
		if ok {
			base -= structs.MaxEthos / 10 // research rivalry
			sharedResearch = true
		}
	}
	if !sharedResearch {
		base += structs.MaxEthos / 10 // research cooperation
	}

	// favoured actions
	actions := map[structs.ActionType]bool{}
	for _, act := range f2.Actions { // their actions
		actions[act] = true
		switch act {
		case structs.ActionTypeFestival, structs.ActionTypeCharity:
			base += structs.MaxEthos / 20 // they're nice / charitable
		case structs.ActionTypeWar, structs.ActionTypeCrusade, structs.ActionTypeShadowWar:
			base -= structs.MaxEthos / 5 // they're warlike
		case structs.ActionTypeRaid, structs.ActionTypePillage, structs.ActionTypeEnslave:
			base -= structs.MaxEthos / 10 // they're raiders
		}
	}
	for _, act := range f1.Actions { // my actions
		switch act {
		case structs.ActionTypeTrade, structs.ActionTypeCraft, structs.ActionTypeHarvest:
			_, ok := actions[structs.ActionTypeSteal]
			if ok {
				base -= structs.MaxEthos / 10 // dirty theif!
			}
		case structs.ActionTypeSteal, structs.ActionTypeRaid, structs.ActionTypePillage, structs.ActionTypeEnslave:
			_, trade := actions[structs.ActionTypeTrade]
			_, craft := actions[structs.ActionTypeCraft]
			_, harvest := actions[structs.ActionTypeHarvest]
			if trade || craft || harvest {
				base -= structs.MaxEthos / 5 // looks like a tempting target ;)
			}
		case structs.ActionTypeConcealSecrets:
			_, investigate := actions[structs.ActionTypeGatherSecrets]
			_, spy := actions[structs.ActionTypeHireSpies]
			if investigate || spy {
				base -= structs.MaxEthos / 5 // we don't like people who snoop on us
			}
		case structs.ActionTypeConsolidate, structs.ActionTypePropoganda:
			_, frame := actions[structs.ActionTypeFrame]
			_, blackmail := actions[structs.ActionTypeBlackmail]
			_, propoganda := actions[structs.ActionTypePropoganda]
			if frame || blackmail || propoganda {
				base -= structs.MaxEthos / 5 // we don't like people who talk about us
			}
		}
	}

	return base
}

func isFactionIllegal(govt *structs.Government, f *simutil.MetaFaction) bool {
	// look for any reason to mark the faction as illegal
	if f.Faction.ReligionID != "" {
		illegal, _ := govt.Outlawed.Religions[f.Faction.ReligionID]
		if illegal {
			return true
		}
	}
	for _, act := range f.Actions {
		illegal, _ := govt.Outlawed.Actions[act]
		if illegal {
			return true
		}
	}
	for _, research := range f.ResearchWeights {
		illegal, _ := govt.Outlawed.Research[research.Object]
		if illegal {
			return true
		}
	}
	return false
}

func makeFactionCovert(fr *factionRand, f *simutil.MetaFaction) {
	f.Faction.IsCovert = true

	tenth := structs.MaxTuple / 10
	for _, p := range f.Plots {
		if f.Faction.HQPlotID == p.ID {
			// Headquarters plots are always hidden & harder to find
			p.Hidden = fr.rng.Intn(structs.MaxTuple/2) + structs.MaxTuple/2
		} else if fr.rng.Float64() < 0.90 {
			// 90% of plots of "covert" are hidden (or fronts that masquerade as something else)
			p.Hidden = fr.rng.Intn(structs.MaxTuple-tenth) + tenth
		}
	}
}

// areaGovernments returns
// 1. a map of area id to government id
// 2. a map of government id to government
func (s *Base) areaGovernments(in *db.Query) (map[string]string, map[string]*structs.Government, error) {
	areaToGovt := map[string]string{}
	govtToGovt := map[string]*structs.Government{}

	var (
		areas []*structs.Area
		token string
		err   error
	)

	for {
		areas, token, err = s.dbconn.Areas(token, in)
		if err != nil {
			return nil, nil, err
		}

		for _, area := range areas {
			if dbutils.IsValidID(area.GovernmentID) {
				areaToGovt[area.ID] = area.GovernmentID
				govtToGovt[area.GovernmentID] = nil
			} else {
				areaToGovt[area.ID] = ""
			}
		}

		if token == "" {
			break
		}
	}

	if len(govtToGovt) == 0 {
		return areaToGovt, govtToGovt, nil
	}

	gids := []string{}
	for gid := range govtToGovt {
		gids = append(gids, gid)
	}
	gf := db.Q(db.F(db.ID, db.In, gids))

	var governments []*structs.Government

	for {
		governments, token, err = s.dbconn.Governments(token, gf)
		if err != nil {
			return nil, nil, err
		}

		for _, govt := range governments {
			govtToGovt[govt.ID] = govt
		}

		if token == "" {
			break
		}
	}

	return areaToGovt, govtToGovt, nil
}

// randFaction spits out a random faction.
func (s *Base) randFaction(fr *factionRand) *simutil.MetaFaction {
	// start with a lot of randomly inserted fields
	mf := simutil.NewMetaFaction()
	mf.Faction = &structs.Faction{
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
	}

	// consider action focuses
	professions := []string{}
	actions := []structs.ActionType{}
	seen := map[int]bool{}
	actionsEthos := []*structs.Ethos{&mf.Faction.Ethos}
	researchCount := 0
	researchTopics := []string{}
	for i := 0; i < fr.focusCount.Int(); i++ {
		choice := fr.focusOccur.Int()
		_, ok := seen[choice]
		if ok {
			continue
		}
		seen[choice] = true

		focus := fr.cfg.Focuses[choice]
		weight := fr.focusWeights[choice]

		mf.Faction.EspionageOffense += int(focus.EspionageOffenseBonus * float64(mf.Faction.EspionageOffense))
		mf.Faction.EspionageDefense += int(focus.EspionageDefenseBonus * float64(mf.Faction.EspionageDefense))
		mf.Faction.MilitaryOffense += int(focus.MilitaryOffenseBonus * float64(mf.Faction.MilitaryOffense))
		mf.Faction.MilitaryDefense += int(focus.MilitaryDefenseBonus * float64(mf.Faction.MilitaryDefense))

		for _, act := range focus.Actions {
			actionCfg, ok := s.cfg.Actions[act]
			if !ok {
				continue
			}

			if act == structs.ActionTypeResearch {
				researchCount++
				if focus.ResearchTopics != nil && len(focus.ResearchTopics) > 0 {
					researchTopics = append(researchTopics, focus.ResearchTopics...)
				}
			}

			mf.ActionWeights = append(mf.ActionWeights, &structs.Tuple{
				Subject: mf.Faction.ID,
				Object:  string(act),
				Value:   weight.Int(),
			})
			actions = append(actions, act)

			actionsEthos = append(actionsEthos, &actionCfg.Ethos)
			for prof := range actionCfg.ProfessionWeights {
				professions = append(professions, prof)
			}
		}
	}

	// favoured actions influence faction's starting ethos
	mf.Faction.Ethos = *structs.EthosAverage(actionsEthos...)

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
		seen[choice] = true

		guild := fr.cfg.Guilds[choice]
		landrights := fr.randLandForGuild(&guild)

		if landrights == nil {
			continue
		}

		actions = append(actions, structs.ActionTypeTrade)

		for _, land := range landrights {
			land.FactionID = mf.Faction.ID

			if s.eco.IsCraftable(land.Commodity) {
				countCraft++
				actions = append(actions, structs.ActionTypeCraft)
			}
			if s.eco.IsHarvestable(land.Commodity) {
				countHarvest++
				actions = append(actions, structs.ActionTypeHarvest)
			}
			mf.Areas[land.AreaID] = true
		}
		mf.Plots = append(mf.Plots, landrights...)
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
		mf.ProfWeights = append(mf.ProfWeights, &structs.Tuple{
			Subject: mf.Faction.ID,
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
		mf.Actions = append(mf.Actions, a)
		mf.ActionWeights = append(mf.ActionWeights, &structs.Tuple{
			Subject: mf.Faction.ID,
			Object:  string(a),
			Value:   w,
		})
	}

	// give faction land if we're still below the min
	for i := len(mf.Plots); i < fr.propertyCount.Int()+1; i++ {
		area := fr.areas[fr.rng.Intn(len(fr.areas))]
		mf.Plots = append(mf.Plots, &structs.Plot{
			ID:        structs.NewID(),
			AreaID:    area,
			FactionID: mf.Faction.ID,
			Crop: structs.Crop{
				Size: fr.plotSize.Int(),
			},
		})
		mf.Areas[area] = true
	}
	mf.Faction.Plots = len(mf.Plots)
	mf.Faction.Areas = len(mf.Areas)

	// if we need to pick research topics, we can do so now sensibly
	// (since we know where the faction is and what professions it prefers)
	mf.ResearchWeights = fr.randResearch(mf, researchCount, researchTopics)

	// pick a headquarters
	switch len(mf.Plots) {
	case 0:
		mf.Faction.HomeAreaID = mf.Plots[0].AreaID
		mf.Faction.HQPlotID = mf.Plots[0].ID
	default:
		rnum := fr.rng.Intn(len(mf.Plots))
		mf.Faction.HomeAreaID = mf.Plots[rnum].AreaID
		mf.Faction.HQPlotID = mf.Plots[rnum].ID
	}

	return mf
}
