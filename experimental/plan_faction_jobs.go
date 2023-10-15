package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

func main() {
	simulator, err := sim.New(nil)
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
