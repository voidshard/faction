package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

func printDemo(demo *structs.DemographicStats) {
	fmt.Println("\t\tExemplary", demo.Exemplary)
	fmt.Println("\t\tExcellent", demo.Excellent)
	fmt.Println("\t\tGood", demo.Good)
	fmt.Println("\t\tFine", demo.Fine)
	fmt.Println("\t\tAverage", demo.Average)
	fmt.Println("\t\tPoor", demo.Poor)
	fmt.Println("\t\tAwful", demo.Awful)
	fmt.Println("\t\tTerrible", demo.Terrible)
	fmt.Println("\t\tAbsymal", demo.Absymal)
}

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

	faith, err := simulator.FaithDemographics()
	if err != nil {
		panic(err)
	}
	fmt.Println("Faith")
	for obj, demo := range faith {
		fmt.Println("\t", obj)
		printDemo(demo)
	}

	professions, err := simulator.ProfessionDemographics()
	if err != nil {
		panic(err)
	}
	fmt.Println("Professions")
	for obj, demo := range professions {
		fmt.Println("\t", obj)
		printDemo(demo)
	}

	affiliations, err := simulator.AffiliationDemographics()
	if err != nil {
		panic(err)
	}
	fmt.Println("Affiliation")
	for obj, demo := range affiliations {
		fmt.Println("\t", obj)
		printDemo(demo)
	}
}
