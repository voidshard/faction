package main

import (
	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

func main() {
	simulator, err := sim.New(nil)
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
