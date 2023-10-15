/*
random_faction.go - random faction / government generation
*/
package base

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/internal/sim/simutil"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

func (s *Base) SpawnFaction(cfg *config.Faction, areas ...string) (*structs.Faction, error) {
	tick, err := s.dbconn.Tick()
	if err != nil {
		return nil, err
	}

	// dice based on faction cfg + land yields
	dice := newFactionRand(cfg, s.tech, areas)

	// lookup which areas are run by which governments (determines faction legality / covert status)
	arf := db.Q(db.F(db.ID, db.In, areas))
	areaToGovt, govtToGovt, err := s.areaGovernments(arf)
	if err != nil {
		return nil, err
	}

	// generate factions
	govsToWrite := map[string]*structs.Government{}
	f, err := s.randFaction(dice)
	if err != nil {
		return nil, err
	}

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

	return f.Faction, err
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
func (s *Base) randFaction(fr *factionRand) (*simutil.MetaFaction, error) {
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
	professions := map[string]int{}
	actions := map[structs.ActionType]int{}
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
			actions[act]++

			actionsEthos = append(actionsEthos, &actionCfg.Ethos)
			for prof := range actionCfg.ProfessionWeights {
				professions[prof]++
			}
		}
	}

	// favoured actions influence faction's starting ethos
	mf.Faction.Ethos = *structs.EthosAverage(actionsEthos...)

	// consider a profession based guild
	seen = map[int]bool{}
	desiredArea := fr.plotSize.Int()
	landRequirements := map[string]int{"": desiredArea}
	guilds := []*config.Guild{}
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
		if guild.LandMinCommodityYield == nil {
			continue // guild doesn't need land .. technically this config is invalid
		}

		for commodity, requirement := range guild.LandMinCommodityYield {
			// record that we'd like some land
			v, _ := landRequirements[commodity]
			landRequirements[commodity] = v + requirement
		}

		guilds = append(guilds, &guild)
	}

	// take some land for our faction
	plots, resources, err := s.chooseFactionLand(fr.areas, landRequirements, mf.Faction.ID)
	if err != nil {
		return nil, err
	}

	// make land (if needed & permitted)
	morePlots, err := s.createFactionLand(landRequirements, resources, mf.Faction.ID, fr.areas, fr.cfg.AllowEmptyPlotCreation, fr.cfg.AllowCommodityPlotCreation)
	if err != nil {
		if errors.Is(err, structs.ErrNotEnoughPlots) {
			emptyLand, _ := resources[""]
			if emptyLand < desiredArea {
				// We needed to make land and could not
				return nil, err
			}
			// Nb. we allow "NotEnoughPlots" when spawning a faction as we added Guild(s) at random
			// which require specialised commodity based land without knowing if that land was actually available.
			// In such a case, we simply spawn the faction with the land we have, and drop the Guild we added.
		} else {
			return nil, err
		}
	}

	// check if our guild(s) are happy with the land we have
	imports := map[string]int{}
	exports := map[string]int{}
	for _, guild := range guilds {
		if guild.LandMinCommodityYield == nil {
			continue // guild doesn't need land .. technically this config is invalid
		}

		for commodity, requirement := range guild.LandMinCommodityYield {
			have, _ := resources[commodity]
			if have < requirement {
				continue // not enough land for this guild, we'll pretend we never wanted it
			}
			// otherwise, this faction is now a proud guild

			// update our counts (for tuples later)
			for _, prof := range guild.Professions {
				professions[prof]++
			}
			for _, c := range guild.Imports {
				imports[c]++
			}
			for _, c := range guild.Exports {
				exports[c]++
				if s.eco.IsCraftable(c) {
					// presumably we're planning on crafting this thing we export
					actions[structs.ActionTypeCraft]++
				}
			}

			// update action counts
			if len(guild.Imports) > 0 || len(guild.Exports) > 0 {
				actions[structs.ActionTypeTrade]++
			}

			if s.eco.IsHarvestable(commodity) { // I mean, this .. should be true, right? .. RIGHT?
				actions[structs.ActionTypeHarvest]++
			}
		}
	}

	// record everything
	mf.Plots = append(mf.Plots, plots...)
	mf.Plots = append(mf.Plots, morePlots...)
	if len(mf.Plots) == 0 {
		return nil, fmt.Errorf("%w no land found for faction matching requirements %v", structs.ErrNotEnoughPlots, landRequirements)
	}

	// pick a home area / HQ
	rnum := fr.rng.Intn(len(mf.Plots))
	mf.Faction.HomeAreaID = mf.Plots[rnum].AreaID
	mf.Faction.HQPlotID = mf.Plots[rnum].ID

	// now we can determine weights based on counts (more counts -> higher weight)
	for p, c := range professions {
		w := clamp(c*structs.MaxEthos/20+fr.rng.Intn(500), 0, structs.MaxEthos)
		mf.ProfWeights = append(mf.ProfWeights, &structs.Tuple{Subject: mf.Faction.ID, Object: p, Value: w})
	}
	for c, i := range imports {
		w := clamp(i*structs.MaxEthos/20+fr.rng.Intn(500), 0, structs.MaxEthos)
		mf.Imports = append(mf.Imports, &structs.Tuple{Subject: mf.Faction.ID, Object: c, Value: w})
	}
	for c, i := range exports {
		w := clamp(i*structs.MaxEthos/20+fr.rng.Intn(500), 0, structs.MaxEthos)
		mf.Exports = append(mf.Exports, &structs.Tuple{Subject: mf.Faction.ID, Object: c, Value: w})
	}
	for a, i := range actions {
		w := clamp(i*structs.MaxEthos/20+fr.rng.Intn(500), 0, structs.MaxEthos)
		mf.Actions = append(mf.Actions, a)
		mf.ActionWeights = append(mf.ActionWeights, &structs.Tuple{
			Subject: mf.Faction.ID,
			Object:  string(a),
			Value:   w,
		})
	}

	// if we need to pick research topics, we can do so now sensibly
	// (since we know where the faction is and what professions it prefers)
	mf.ResearchWeights = fr.randResearch(mf, researchCount, researchTopics)

	return mf, nil
}

func (s *Base) chooseFactionLand(areas []string, requirements map[string]int, factionID string) ([]*structs.Plot, map[string]int, error) {
	requirementMet := 0
	resourcesSoFar := map[string]int{}
	found := []*structs.Plot{}

	commodities := []string{}
	for commodity := range requirements {
		commodities = append(commodities, commodity)
	}

	q := db.Q(db.F(db.AreaID, db.In, areas), db.F(db.FactionID, db.Equal, ""), db.F(db.Commodity, db.In, commodities))

	var (
		plots []*structs.Plot
		token string
		err   error
	)
	for {
		plots, token, err = s.dbconn.Plots(token, q)
		if err != nil {
			return nil, nil, err
		}

		for _, plot := range plots {
			want, _ := requirements[plot.Commodity]
			have, _ := resourcesSoFar[plot.Commodity]

			if have >= want {
				continue
			}

			metric := plot.Size
			if plot.Commodity != "" {
				metric = plot.Yield * plot.Size
			}

			plot.FactionID = factionID
			found = append(found, plot)
			resourcesSoFar[plot.Commodity] += metric

			if have+metric >= want {
				requirementMet++
			}
			if requirementMet >= len(requirements) {
				return found, resourcesSoFar, nil
			}
		}

		if token == "" {
			break
		}
	}

	return found, resourcesSoFar, nil
}

func (s *Base) createFactionLand(want, have map[string]int, factionID string, areas []string, createEmpty, createCmd bool) ([]*structs.Plot, error) {
	newLand := []*structs.Plot{}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var err error
	for commodity, desired := range want { // check if we have to / can make stuff up
		obtained, _ := have[commodity]
		if obtained >= desired {
			continue
		}
		need := desired - obtained

		areaID := areas[0]
		if len(areas) > 1 {
			areaID = areas[rng.Intn(len(areas))]
		}

		if commodity == "" && createEmpty || commodity != "" && createCmd {
			// create land from thin air
			newLand = append(newLand, &structs.Plot{
				ID:        structs.NewID(),
				AreaID:    areaID,
				FactionID: factionID,
				Crop: structs.Crop{
					Commodity: commodity,
					Size:      need,
					Yield:     1,
				},
				Value: s.eco.LandValue(areaID, 0) * float64(need),
			})
			have[commodity] = desired
		} else {
			if err == nil {
				err = fmt.Errorf("%w when spawning faction, want %v units of %s", structs.ErrNotEnoughPlots, desired, commodity)
			} else {
				err = fmt.Errorf("%w, want %v units of %s", err, desired, commodity)
			}
		}
	}

	return newLand, err
}
