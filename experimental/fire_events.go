package main

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/queue"
	"github.com/voidshard/faction/pkg/sim"
)

func main() {
	local, err := queue.NewAsynqQueue(config.DefaultQueue())
	if err != nil {
		panic(err)
	}

	simulator, err := sim.New(nil)
	if err != nil {
		panic(err)
	}

	err = simulator.SetQueue(local)
	if err != nil {
		panic(err)
	}
	err = simulator.StartProcessingEvents()
	if err != nil {
		panic(err)
	}

	defer simulator.StopProcessingEvents()

	err = simulator.FireEvents()
	if err != nil {
		panic(err)
	}
}
