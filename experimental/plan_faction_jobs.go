package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
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
		Actions: fantasy.Actions(),
	}

	simulator, err := sim.New(cfg)
	if err != nil {
		panic(err)
	}

	f01 := structs.NewID("faction1")
	f02 := structs.NewID("faction2")
	f03 := structs.NewID("faction3")

	factions, err := simulator.Factions(f01, f02, f03)
	if err != nil {
		panic(err)
	}

	for _, f := range factions {
		jobs, err := simulator.PlanFactionJobs(f.ID)
		if err != nil {
			fmt.Println("error planning faction jobs:", err)
			return
		}

		fmt.Println("Faction", f.ID, "jobs:")
		for _, j := range jobs {
			fmt.Println("  ", j)
		}
	}
}
