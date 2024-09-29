package main

import (
	"context"
	"fmt"

	"github.com/voidshard/faction/pkg/structs"
)

type cliListCmd struct {
	optCliConn

	Object struct {
		Name string `positional-arg-name:"object" description:"Object to list"`
	} `positional-args:"true" required:"true"`

	Limit  int `long:"limit" default:"100" description:"Limit number of results"`
	Offset int `short:"o" long:"offset" default:"0" description:"Offset results"`

	Labels map[string]string `short:"l" long:"labels" description:"Filter by labels"`
}

func (c *cliListCmd) Execute(args []string) error {
	obj := validObject(c.Object.Name)
	if obj == nil {
		return invalidObjectError()
	}

	conn, err := newClient(c.Host, c.Port, c.IdleTimeout, c.ConnTimeout)
	if err != nil {
		return err
	}

	lim := uint32(c.Limit)
	off := uint32(c.Offset)

	switch obj.(type) {
	case *structs.Actor:
		resp, err := conn.ListActors(context.TODO(), &structs.ListActorsRequest{Limit: &lim, Offset: &off, Labels: c.Labels})
		return dumpTable(resp.Data, resp.Error, err)
	case *structs.Faction:
		resp, err := conn.ListFactions(context.TODO(), &structs.ListFactionsRequest{Limit: &lim, Offset: &off, Labels: c.Labels})
		return dumpTable(resp.Data, resp.Error, err)
	case *structs.World:
		resp, err := conn.ListWorlds(context.TODO(), &structs.ListWorldsRequest{Limit: &lim, Offset: &off, Labels: c.Labels})
		return dumpTable(resp.Data, resp.Error, err)
	}

	return nil
}

func dumpTable[T structs.Object](data []T, err *structs.Error, connErr error) error {
	if connErr != nil {
		return connErr
	}
	if err != nil {
		return fmt.Errorf("error: %s", err.Message)
	}

	for _, d := range data {
		fmt.Println(d.String())
	}

	return nil
}
