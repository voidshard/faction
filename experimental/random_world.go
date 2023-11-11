package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/voidshard/faction/pkg/premade/fantasy"
	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

const (
	countAreas       = 10
	countPopulace    = 30000
	countGovernments = 3
	countFactions    = 50
)

var (
	cropCommodities = []struct {
		Name     string
		Plots    int
		SizeMin  int
		SizeMax  int
		YieldMin int
		YieldMax int
	}{
		{fantasy.IRON_ORE, 10, 50, 200, 1, 3},
		{fantasy.WHEAT, 50, 50, 500, 1, 5},
		{fantasy.TIMBER, 50, 50, 300, 1, 2},
		{fantasy.WILD_GAME, 25, 200, 1000, 1, 4},
		{fantasy.FLAX, 10, 100, 600, 1, 2},
		{fantasy.FLAX, 30, 50, 500, 1, 3},
		{fantasy.OPIUM, 5, 10, 50, 1, 3},
	}
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	simulator, err := sim.New(nil)
	if err != nil {
		panic(err)
	}

	// step 1. governments (so we can assign areas as we go)
	fmt.Println("spawning governments")
	governments := []*structs.Government{}
	for i := 0; i < countGovernments; i++ {
		gov, err := simulator.SpawnGovernment(fantasy.Government())
		if err != nil {
			panic(err)
		}
		governments = append(governments, gov)

		fmt.Printf("\tgovernment %s\n", gov.ID)
	}

	// step 2. areas (broad swaths of land run by at most one government)
	fmt.Println("spawning areas")
	areas := []*structs.Area{}
	areaIDs := []string{}
	areasByGovernment := map[string][]string{}
	for i := 0; i < countAreas; i++ {
		gov := governments[i%len(governments)]
		area := &structs.Area{ID: structs.NewID(), GovernmentID: gov.ID}

		fmt.Printf("\tarea %s [Government: %s]\n", area.ID, gov.ID)

		areas = append(areas, area)
		areaIDs = append(areaIDs, area.ID)

		byGov, ok := areasByGovernment[gov.ID]
		if !ok {
			byGov = []string{}
		}
		areasByGovernment[gov.ID] = append(byGov, area.ID)
	}
	err = simulator.SetAreas(areas...)
	if err != nil {
		panic(err)
	}

	// step 3. insert usable plots with crops into areas at random
	fmt.Println("spawning plots")
	plots := []*structs.Plot{}
	for _, crop := range cropCommodities {
		for i := 0; i < crop.Plots; i++ {
			randArea := areas[rng.Intn(len(areas))]
			p := &structs.Plot{
				ID:     structs.NewID(),
				AreaID: randArea.ID,
				Crop: structs.Crop{
					Commodity: crop.Name,
					Size:      rng.Intn(crop.SizeMax-crop.SizeMin) + crop.SizeMin,
					Yield:     rng.Intn(int(crop.YieldMax-crop.YieldMin)) + crop.YieldMin,
				},
			}
			plots = append(plots, p)

			fmt.Println("\tplot", p.ID, "in", p.AreaID, "of size", p.Crop.Size, "with", p.Crop.Commodity, "with average yield", p.Crop.Yield)
		}
	}
	err = simulator.SetPlots(plots...)
	if err != nil {
		panic(err)
	}

	// step 4. create some people spread fairly evenly (this is easier here, realisitically distribution should be uneven)
	fmt.Println("spawning populace")
	err = simulator.SpawnPopulace(countPopulace, "human", "human", areaIDs...)
	if err != nil {
		panic(err)
	}

	// step 5. create some factions (by government is easy here, but not required)
	fmt.Println("spawning factions")
	factionsByArea := countFactions / countGovernments
	governmentFactions := []*structs.Faction{}
	factions := []*structs.Faction{}
	for govID, govAreaIDs := range areasByGovernment {
		fmt.Printf("\tspawning %d factions for government %s\n", factionsByArea, govID)

		localFactions := []*structs.Faction{}
		for i := 0; i < factionsByArea; i++ {
			f, err := simulator.SpawnFaction(fantasy.Faction(), govAreaIDs...)
			if err != nil {
				panic(err)
			}
			localFactions = append(localFactions, f)
		}

		// select someone as government
		govFaction := localFactions[rng.Intn(len(localFactions))]
		govFaction.GovernmentID = govID
		govFaction.IsGovernment = true

		fmt.Println("\t\t", govFaction.ID, "is government of", govFaction.GovernmentID)

		// record
		governmentFactions = append(governmentFactions, govFaction)
		factions = append(factions, localFactions...)

		for _, f := range localFactions {
			fmt.Println("\t\t", f.ID, "faction under", f.GovernmentID)
		}
	}
	err = simulator.SetFactions(governmentFactions...) // record who we marked as the government
	if err != nil {
		panic(err)
	}

	// step 6. inspire affiliation (ie. add people to factions)
	fmt.Println("inspiring faction affiliation")
	for _, f := range factions {
		fmt.Println("\tinspiring", f.ID, "faction affiliation")
		err = simulator.InspireFactionAffiliation(fantasy.Affiliation(), f.ID)
		if err != nil {
			panic(err)
		}
	}

	for _, f := range factions {
		fmt.Println("\tplanning jobs for", f.ID)
		jobs, err := simulator.PlanFactionJobs(f.ID)
		if err != nil {
			panic(err)
		}
		for _, j := range jobs {
			fmt.Println(
				"\t\t",
				j.ID, j.Action, "job for", j.SourceFactionID,
				"targeting", j.TargetFactionID, "in", j.TargetAreaID,
				"(", j.TargetMetaKey, j.TargetMetaVal, ")",
			)
		}
	}

	// step 7. fire events (handles all the post processing for all the changes we enacted)
	//fmt.Println("firing events")
	//err = simulator.FireEvents()
	//if err != nil {
	//		panic(err)
	//	}
}
