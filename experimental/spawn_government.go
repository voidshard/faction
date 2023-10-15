package main

import (
	"fmt"

	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
	"github.com/voidshard/faction/pkg/sim"
)

func main() {
	simulator, err := sim.New(nil)
	if err != nil {
		panic(err)
	}

	govt, err := simulator.SpawnGovernment(fantasy.Government())
	if err != nil {
		panic(err)
	}
	fmt.Println(govt.ID, "\n\tTax:", govt.TaxRate, govt.TaxFrequency)
	fmt.Println("\tOutlawed:")
	fmt.Println("\t\tFactions:", govt.Outlawed.Factions)
	fmt.Println("\t\tActions:", govt.Outlawed.Actions)
	fmt.Println("\t\tCommodities:", govt.Outlawed.Commodities)
	fmt.Println("\t\tResearch:", govt.Outlawed.Research)
	fmt.Println("\t\tReligions:", govt.Outlawed.Religions)
}
