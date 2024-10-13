package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/structs"
)

type cliWatchCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	// World is set in optCliGlobal
	Area string `long:"area" short:"a" description:"Area to watch for changes" default:""`

	Object string `long:"object" short:"o" description:"Object type to watch for changes" default:""`

	Id string `long:"id" short:"i" description:"ID of object(s) to watch" default:""`

	Queue string `long:"queue" description:"Queue to watch for changes, if set queue is durable" default:""`
}

func (c *cliWatchCmd) Execute(args []string) error {
	if c.World == "" {
		return fmt.Errorf("world is required")
	}
	c.World = toWorldId(c.World)[0]

	key := validKey(c.Object)
	if key == structs.Metakey_KeyNone && c.Object != "" {
		// you tried to watch an object, but we got a None key
		// presumably this isn't what you wanted
		return fmt.Errorf("invalid object type %s", c.Object)
	}

	client, err := newClient(c.Host, c.Port, c.IdleTimeout, c.ConnTimeout)
	if err != nil {
		return err
	}

	l := log.Sublogger("cli.Watch", map[string]string{
		"world":  c.World,
		"area":   c.Area,
		"object": c.Object,
		"id":     c.Id,
		"queue":  c.Queue,
	})
	l.Info().Msg("Starting watch")
	defer l.Debug().Msg("Stopped watch")
	sub, err := client.OnChange(context.Background(), &structs.OnChangeRequest{
		Data: &structs.Change{
			World: c.World,
			Area:  c.Area,
			Key:   key,
			Id:    c.Id,
		},
		Queue: c.Queue,
	})
	if err != nil {
		return err
	}

	for {
		resp, err := sub.Recv()
		l.Debug().Err(err).Msg("Received change")
		if err == io.EOF {
			l.Debug().Err(err).Msg("watch stream closed")
			return nil // stream closed
		} else if err != nil {
			l.Error().Err(err).Msg("watch stream error")
			return err // some other error
		} else if resp == nil {
			l.Warn().Msg("watch stream got nil response")
			continue // ???
		} else if resp.Error != nil {
			l.Error().Str("error", resp.Error.Message).Msg("watch stream error")
			// error sent from server
			return fmt.Errorf("error in watch: %s", resp.Error.Message)
		}

		rkey, ok := structs.Metakey_name[int32(resp.Data.Key)]
		if !ok {
			rkey = "None"
			l.Warn().Str("key", fmt.Sprintf("%d", resp.Data.Key)).Msg("Unknown key from server")
		} else {
			rkey = strings.Replace(rkey, "Key", "", 1)
		}

		fmt.Printf("World: %s | Area: %s | Object: %s | Id: %s\n", resp.Data.World, resp.Data.Area, rkey, resp.Data.Id)
	}
}
