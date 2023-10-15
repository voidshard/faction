package main

import (
	"github.com/voidshard/faction/pkg/sim"
	"github.com/voidshard/faction/pkg/structs"
)

func main() {
	popPerArea := 100000
	areas := 10

	simulator, err := sim.New(nil)
	if err != nil {
		panic(err)
	}

	allAreas := []*structs.Area{}
	ids := []string{}
	for i := 0; i < areas; i++ {
		allAreas = append(allAreas, &structs.Area{ID: structs.NewID("area.perf.%d", i)})
		ids = append(ids, allAreas[i].ID)
	}

	err = simulator.SetAreas(allAreas...)
	if err != nil {
		panic(err)
	}

	err = simulator.SpawnPopulace(
		popPerArea,
		"human",
		"human",
		ids...,
	)
	if err != nil {
		panic(err)
	}
}
