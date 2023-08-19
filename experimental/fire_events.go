package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/queue"
	"github.com/voidshard/faction/pkg/sim"
)

func main() {
	cfg := &config.Database{
		Driver:   config.DatabaseSQLite3,
		Name:     "test.sqlite",
		Location: "/tmp",
	}
	local := queue.NewLocalQueue(10)

	simulator, err := sim.New(&config.Simulation{Database: cfg})
	if err != nil {
		panic(err)
	}

	err = simulator.SetQueue(local)
	if err != nil {
		panic(err)
	}

	err = simulator.FireEvents()
	if err != nil {
		panic(err)
	}

	result, err := local.Await() // await all
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
