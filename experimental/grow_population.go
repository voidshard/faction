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

	for _, a := range []*structs.Area{area1} {
		err = simulator.AdjustPopulation(a.ID)
		if err != nil {
			panic(err)
		}
	}
}
