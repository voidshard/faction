package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/client"
)

type cliWatchCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	// World is set in optCliGlobal
	Kind string `long:"kind" short:"k" description:"Kind of object to watch for changes" default:""`

	Id string `long:"id" short:"i" description:"ID of object(s) to watch" default:""`

	Queue string `long:"queue" description:"Queue to watch for changes, if set queue is durable" default:""`
}

func (c *cliWatchCmd) Execute(args []string) error {
	conn, err := client.New(client.NewConfig())
	if err != nil {
		return err
	}

	sub, err := conn.Watch().Kind(c.Kind).World(c.World).Id(c.Id).Queue(c.Queue).Do()
	if err != nil {
		return err
	}
	defer sub.Close()

	for event := range sub.Events() {
		fmt.Println(event)
		if c.Queue != "" {
			err = sub.Ack(event.AckId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
