package main

import (
	"math/rand"
	"time"

	"github.com/voidshard/faction/pkg/config"
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

	simulator, err := sim.New(&config.Simulation{Database: cfg}, nil)
	if err != nil {
		panic(err)
	}

	gid := structs.NewID("government1")
	rid := structs.NewID("religion1")

	faction1 := &structs.Faction{ID: structs.NewID("faction1"), GovernmentID: gid, Ethos: randomEthos()}
	faction2 := &structs.Faction{ID: structs.NewID("faction2"), GovernmentID: gid, Ethos: randomEthos(), IsReligion: true, ReligionID: rid}
	faction3 := &structs.Faction{ID: structs.NewID("faction3"), GovernmentID: gid, Ethos: randomEthos()}

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

	err = simulator.InspireFactionAffiliation(
		[]*structs.Faction{faction1, faction2, faction3},
		&config.Distribution{Min: 1000, Max: structs.MaxTuple, Mean: 2000, Deviation: 6000},
		0.2,
		0, 25000,
	)
	if err != nil {
		panic(err)
	}
}
