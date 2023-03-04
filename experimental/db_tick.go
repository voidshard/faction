package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/sim"
)

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

	t, err := simulator.Tick()
	fmt.Printf("tick %d (%v)\n", t, err)
}
