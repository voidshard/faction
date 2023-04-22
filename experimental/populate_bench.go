package main

import (
	"github.com/voidshard/faction/pkg/config"
	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

func main() {
	popPerArea := 100000
	areas := 10

	cfg := &config.Database{
		Driver:   config.DatabaseSQLite3,
		Name:     "test.sqlite",
		Location: "/tmp",
	}

	simulator, err := sim.New(&config.Simulation{Database: cfg})
	if err != nil {
		panic(err)
	}

	allAreas := []*structs.Area{}
	ids := []string{}
	for i := 0; i < areas; i++ {
		allAreas = append(allAreas, &structs.Area{ID: structs.NewID("area.perf.%d", i)})
		ids = append(ids, allAreas[i].ID)
	}

	err = simulator.SetAreas(allAreas...)
	if err != nil {
		panic(err)
	}

	err = simulator.SpawnPopulace(
		popPerArea,
		fantasy.DemographicsHuman(),
		ids...,
	)
	if err != nil {
		panic(err)
	}
}
