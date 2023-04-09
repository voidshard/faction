package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/premade"
	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

func main() {
	cfg := &config.Simulation{
		Database: &config.Database{
			Driver:   config.DatabaseSQLite3,
			Name:     "test.sqlite",
			Location: "/tmp",
		},
		Actions: premade.ActionsFantasy(),
	}

	simulator, err := sim.New(cfg)
	if err != nil {
		panic(err)
	}

	govt, err := simulator.SpawnGovernment(
		premade.GovernmentFantasy(),
	)
	if err != nil {
		panic(err)
	}

	area1 := &structs.Area{ID: structs.NewID(), GovernmentID: govt.ID}
	area2 := &structs.Area{ID: structs.NewID(), GovernmentID: govt.ID}
	area3 := &structs.Area{ID: structs.NewID()} // no govt
	err = simulator.SetAreas(area1, area2, area3)
	if err != nil {
		panic(err)
	}

	land1 := &structs.Plot{ID: structs.NewID(), AreaID: area1.ID, Commodity: premade.WHEAT, Yield: 120}
	land2 := &structs.Plot{ID: structs.NewID(), AreaID: area2.ID, Commodity: premade.WHEAT, Yield: 200}
	land3 := &structs.Plot{ID: structs.NewID(), AreaID: area2.ID, Commodity: premade.IRON_ORE, Yield: 20}
	land4 := &structs.Plot{ID: structs.NewID(), AreaID: area3.ID, Commodity: premade.WHEAT, Yield: 50}
	land5 := &structs.Plot{ID: structs.NewID(), AreaID: area3.ID, Commodity: premade.IRON_ORE, Yield: 500}
	err = simulator.SetPlots(land1, land2, land3, land4, land5)
	if err != nil {
		panic(err)
	}

	fcfg := premade.FactionFantasy()

	factions, err := simulator.SpawnFactions(10, fcfg, area1.ID, area2.ID, area3.ID)
	if err != nil {
		panic(err)
	}

	for _, f := range factions {
		fmt.Println("ID:", f.ID, "\n\tHomeAreaID:", f.HomeAreaID, "\n\tIsCovert:", f.IsCovert, "\n\tGovernmentID:", f.GovernmentID)
		fmt.Println("\tAltruism:", f.Altruism, "Ambition:", f.Ambition, "Tradition:", f.Tradition, "Pacifism:", f.Pacifism, "Piety:", f.Piety, "Caution:", f.Caution)
		fmt.Println("\tEspionageOffense", f.EspionageOffense, "EspionageDefense", f.EspionageDefense)
		fmt.Println("\tMilitaryOffense", f.MilitaryOffense, "MilitaryDefense", f.MilitaryDefense)
		fmt.Println("\tWealth", f.Wealth, "Cohesion", f.Cohesion, "Corruption", f.Corruption)

		fs, err := simulator.FactionSummaries(f.ID)
		if err != nil {
			panic(err)
		}
		if len(fs) != 1 {
			panic("expected 1 faction summary")
		}

		sum := fs[0]
		fmt.Println("\t\tProfessions:", sum.Professions)
		fmt.Println("\t\tActions:", sum.Actions)
		fmt.Println("\t\tResearch:", sum.Research)
		fmt.Println("\t\tTrust:", sum.Trust)
		fmt.Println("\t\tResearchProgress:", sum.ResearchProgress)
		fmt.Println("\t\tRanks:", sum.Ranks)
	}
}
