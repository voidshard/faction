package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/voidshard/faction/pkg/client"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"
	"google.golang.org/grpc"
)

type cliWatchCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	// World is set in optCliGlobal
	Type string `long:"type" short:"t" description:"Type to watch for changes" default:""`

	Object string `long:"object" short:"o" description:"Object type to watch for changes" default:""`

	Id string `long:"id" short:"i" description:"ID of object(s) to watch" default:""`

	Queue string `long:"queue" description:"Queue to watch for changes, if set queue is durable" default:""`
}

func (c *cliWatchCmd) Execute(args []string) error {
	if c.World != "" {
		c.World = toWorldId(c.World)[0]
	}

	c.Object = strings.Title(c.Object)
	if c.Object != "" && !structs.IsValidAPIKind(c.Object) {
		return fmt.Errorf("invalid object type %s", c.Object)
	}

	client, err := client.New(c.Host, c.Port, c.IdleTimeout, c.ConnTimeout)
	if err != nil {
		return err
	}

	l := log.Sublogger("cli.Watch", map[string]interface{}{
		"world": c.World,
		"kind":  c.Object,
		"type":  c.Type,
		"id":    c.Id,
		"queue": c.Queue,
	})
	l.Info().Msg("Starting watch")
	defer l.Debug().Msg("Stopped watch")
	sub, err := client.OnChange(context.Background(), &structs.OnChangeRequest{
		Data: &structs.Change{
			World: c.World,
			Key:   c.Object,
			Type:  c.Type,
			Id:    c.Id,
		},
		Queue: c.Queue,
	})
	if err != nil {
		return err
	}

	var ackstream grpc.ClientStreamingClient[structs.AckRequest, structs.AckResponse]
	if c.Queue != "" {
		ackstream, err = client.AckStream(context.Background())
		if err != nil {
			return err
		}
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

		fmt.Printf("World: %s | Kind: %s | Type: %s | Id: %s\n", resp.Data.World, resp.Data.Key, resp.Data.Type, resp.Data.Id)
		if c.Queue != "" {
			ackstream.Send(&structs.AckRequest{Ack: []string{resp.Ack}})
		}
	}
}
