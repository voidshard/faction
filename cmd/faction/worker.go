package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/sim"
)

const (
	docWorker = `Run faction background worker to process events`
)

type optsWorker struct {
	optsGeneral
	optsQueue
	optsDatabase
}

func (opts *optsWorker) Execute(args []string) error {
	cfg := &config.Simulation{
		Database: config.DefaultDatabase(),
		Queue:    config.DefaultQueue(),
	}

	cfg.Queue.Driver = config.QueueDriver(opts.QueueDriver)
	cfg.Queue.Queue = opts.QueueURL
	cfg.Queue.Database = opts.QueueDatabaseURL
	cfg.Database.Driver = config.DatabaseDriver(opts.DatabaseDriver)
	cfg.Database.Location = opts.DatabaseURL

	simulator, err := sim.New(cfg)
	if err != nil {
		return err
	}

	err = simulator.StartProcessingEvents()
	if err != nil {
		return err
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit

	return nil
}
