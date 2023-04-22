package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/voidshard/faction/pkg/config"
	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	area1 = structs.NewID("area1")
	area2 = structs.NewID("area2")
	area3 = structs.NewID("area3")

	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func randomEthos() structs.Ethos {
	return structs.Ethos{
		Altruism:  int(rng.Float64()*structs.MaxEthos*2) - structs.MaxEthos,
		Ambition:  int(rng.Float64()*structs.MaxEthos*2) - structs.MaxEthos,
		Tradition: int(rng.Float64()*structs.MaxEthos*2) - structs.MaxEthos,
		Pacifism:  int(rng.Float64()*structs.MaxEthos*2) - structs.MaxEthos,
		Piety:     int(rng.Float64()*structs.MaxEthos*2) - structs.MaxEthos,
		Caution:   int(rng.Float64()*structs.MaxEthos*2) - structs.MaxEthos,
	}
}

func main() {
	// Assumes popluace already exist (see experimental/populate.go)
	cfg := &config.Database{
		Driver:   config.DatabaseSQLite3,
		Name:     "test.sqlite",
		Location: "/tmp",
	}

	simulator, err := sim.New(&config.Simulation{Database: cfg})
	if err != nil {
		panic(err)
	}

	gid := structs.NewID("government1")
	rid := structs.NewID("religion1")

	faction1 := &structs.Faction{
		ID:           structs.NewID("faction1"),
		GovernmentID: gid,
		Ethos:        randomEthos(),
		Leadership:   structs.LeaderTypeSingle, // classic pyramid ruler -> officials -> ...
		Structure:    structs.LeaderStructurePyramid,
	}
	faction2 := &structs.Faction{
		ID:           structs.NewID("faction2"),
		GovernmentID: gid,
		Ethos:        randomEthos(),
		IsReligion:   true,
		ReligionID:   rid,
		Leadership:   structs.LeaderTypeTriad,      // three equal leaders
		Structure:    structs.LeaderStructureLoose, // there isn't really a structure to speak of
	}
	faction3 := &structs.Faction{
		ID:           structs.NewID("faction3"),
		GovernmentID: gid,
		Ethos:        randomEthos(),
		Leadership:   structs.LeaderTypeCouncil,   // a group of leaders
		Structure:    structs.LeaderStructureCell, // each leader has their own cell
	}

	err = simulator.SetFactions(faction1, faction2, faction3)
	if err != nil {
		panic(err)
	}

	err = simulator.SetPlots(
		&structs.Plot{ID: structs.NewID("plot1"), FactionID: faction1.ID, AreaID: area1},
		&structs.Plot{ID: structs.NewID("plot2"), FactionID: faction1.ID, AreaID: area1},
		&structs.Plot{ID: structs.NewID("plot3"), FactionID: faction1.ID, AreaID: area2},
		&structs.Plot{ID: structs.NewID("plot4"), FactionID: faction2.ID, AreaID: area2},
	)
	if err != nil {
		panic(err)
	}

	err = simulator.SetPlots(
		&structs.Plot{ID: structs.NewID("lr1"), FactionID: faction2.ID, AreaID: area3, Commodity: "wood"},
		&structs.Plot{ID: structs.NewID("lr2"), FactionID: faction3.ID, AreaID: area3, Commodity: "wood"},
	)
	if err != nil {
		panic(err)
	}

	affil := fantasy.Affiliation()

	summaries, err := simulator.FactionSummaries(faction1.ID, faction2.ID, faction3.ID)
	for _, f := range summaries {
		err = simulator.InspireFactionAffiliation(affil, f.ID)
		if err != nil {
			panic(err)
		}

		summaries, err := simulator.FactionSummaries(f.ID)
		if err != nil {
			panic(err)
		}

		summary := summaries[0]
		fmt.Println("Faction", f.ID, "summary:", summary.Ranks)
	}
}
