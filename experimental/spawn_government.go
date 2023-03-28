package main

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/premade"
	"github.com/voidshard/faction/pkg/sim"
)

func main() {
	cfg := &config.Database{
		Driver:   config.DatabaseSQLite3,
		Name:     "test.sqlite",
		Location: "/tmp",
	}

	simulator, err := sim.New(&config.Simulation{Database: cfg}, premade.NewFantasyEconomy())
	if err != nil {
		panic(err)
	}

	faction, govt, err := simulator.SpawnGovernment(
		premade.GovernmentFantasy(),
	)
	if err != nil {
		panic(err)
	}

}
