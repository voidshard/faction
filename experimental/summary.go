package main

import (
	"fmt"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

func main() {
	faction1 := structs.NewID("faction1")
	//faction2 := structs.NewID("faction2")
	//faction3 := structs.NewID("faction3")
	//

	cfg := &config.Database{
		Driver:   config.DatabaseSQLite3,
		Name:     "test.sqlite",
		Location: "/tmp",
	}

	dbconn, err := db.New(cfg)
	if err != nil {
		panic(err)
	}

	simulator, err := sim.New(&config.Simulation{Database: cfg})
	if err != nil {
		panic(err)
	}

	summaries, err := simulator.FactionSummaries(faction1) // , faction2, faction3)
	if err != nil {
		panic(err)
	}

	areas, err := dbconn.FactionAreas(true, faction1)
	if err != nil {
		panic(err)
	}

	for _, s := range summaries {
		fmt.Println("Faction", s.ID, "summary:", s.Ranks)

		adata, ok := areas[s.ID]
		if !ok {
			continue
		}
		for areaID, area := range adata {
			fmt.Println("\tArea", areaID, area)
		}
	}
}
