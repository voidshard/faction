package main

import (
	"github.com/voidshard/faction/pkg/sim"
)

func main() {
	simulator, err := sim.New(nil)
	if err != nil {
		panic(err)
	}

	err = simulator.StartProcessingEvents()
	if err != nil {
		panic(err)
	}
}
