package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/sim"
)

func main() {
	simulator, err := sim.New(nil)
	if err != nil {
		panic(err)
	}

	t, err := simulator.Tick()
	fmt.Printf("tick %d (%v)\n", t, err)
}
