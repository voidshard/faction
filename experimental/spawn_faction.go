package main

import (
	"fmt"

	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

func main() {
	simulator, err := sim.New(nil)
	if err != nil {
		panic(err)
	}

	govt, err := simulator.SpawnGovernment(fantasy.Government())
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

	// nb. numbers here are just picked from thin air. In actual use where your units mean something this is due
	// some consideration. Size is land units squared, yield is the expected yield (units of commodity) per unit squared of land.
	// What does that mean? Whatever you want really.
	land1 := &structs.Plot{ID: structs.NewID(), AreaID: area1.ID, Crop: structs.Crop{Size: 100, Commodity: fantasy.WHEAT, Yield: 120}}
	land2 := &structs.Plot{ID: structs.NewID(), AreaID: area2.ID, Crop: structs.Crop{Size: 20, Commodity: fantasy.WHEAT, Yield: 200}}
	land3 := &structs.Plot{ID: structs.NewID(), AreaID: area2.ID, Crop: structs.Crop{Size: 300, Commodity: fantasy.IRON_ORE, Yield: 20}}
	land4 := &structs.Plot{ID: structs.NewID(), AreaID: area3.ID, Crop: structs.Crop{Size: 500, Commodity: fantasy.WHEAT, Yield: 50}}
	land5 := &structs.Plot{ID: structs.NewID(), AreaID: area3.ID, Crop: structs.Crop{Size: 30, Commodity: fantasy.IRON_ORE, Yield: 500}}
	err = simulator.SetPlots(land1, land2, land3, land4, land5)
	if err != nil {
		panic(err)
	}

	fcfg := fantasy.Faction()

	for i := 0; i < 10; i++ {
		f, err := simulator.SpawnFaction(fcfg, area1.ID, area2.ID, area3.ID)
		if err != nil {
			panic(err)
		}

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
