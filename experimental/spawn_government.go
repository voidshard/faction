package main

import (
	"fmt"

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

	govt, err := simulator.SpawnGovernment(
		premade.GovernmentFantasy(),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(govt.ID, "\n\tTax:", govt.TaxRate, govt.TaxFrequency)
	fmt.Println("\tOutlawed:")
	fmt.Println("\t\tFactions:", govt.Outlawed.Factions)
	fmt.Println("\t\tActions:", govt.Outlawed.Actions)
	fmt.Println("\t\tCommodities:", govt.Outlawed.Commodities)
}
