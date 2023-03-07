package main

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/premade"
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

	areaIDs := []string{
		structs.NewID(),
		structs.NewID(),
		structs.NewID(),
	}

	err = simulator.Populate(
		10000*len(areaIDs),
		premade.DemographicsFantasyHuman(),
		areaIDs...,
	)
	if err != nil {
		panic(err)
	}
}
