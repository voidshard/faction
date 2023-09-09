package main

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

func main() {
	cfg := &config.Database{
		Driver:   config.DatabaseSQLite3,
		Name:     "test.sqlite",
		Location: "/tmp",
	}

	simulator, err := sim.New(&config.Simulation{Database: cfg})
	if err != nil {
		panic(err)
	}

	area1 := &structs.Area{ID: structs.NewID("area1")}
	area2 := &structs.Area{ID: structs.NewID("area2")}
	area3 := &structs.Area{ID: structs.NewID("area3")}

	err = simulator.SetAreas(area1, area2, area3)
	if err != nil {
		panic(err)
	}

	err = simulator.SpawnPopulace(
		30000, // nb. this is approximate & doesn't include people spawned dead
		"human",
		"human",
		[]string{
			area1.ID,
			area2.ID,
			area3.ID,
		},
	)
	if err != nil {
		panic(err)
	}
}
