package main

import (
	"context"
	"fmt"

	"github.com/voidshard/faction/pkg/client"
	"github.com/voidshard/faction/pkg/structs"
)

type cliListCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Object struct {
		Name string `positional-arg-name:"object" description:"Object to list"`
	} `positional-args:"true" required:"true"`

	Limit  uint64 `long:"limit" default:"100" description:"Limit number of results"`
	Offset uint64 `short:"o" long:"offset" default:"0" description:"Offset results"`

	Labels map[string]string `short:"l" long:"labels" description:"Filter by labels"`
}

func (c *cliListCmd) Execute(args []string) error {
	obj := validObject(c.Object.Name)
	if obj == nil {
		return invalidObjectError(c.Object.Name)
	}

	_, isWorld := obj.(*structs.World)
	if !isWorld && c.World == "" {
		return fmt.Errorf("world must be set for %s", c.Object.Name)
	}
	c.World = toWorldId(c.World)[0]

	conn, err := client.New(c.Host, c.Port, c.IdleTimeout, c.ConnTimeout)
	if err != nil {
		return err
	}

	lim := uint64(c.Limit)
	off := uint64(c.Offset)

	switch obj.(type) {
	case *structs.Actor:
		resp, err := conn.ListActors(context.TODO(), &structs.ListActorsRequest{World: c.World, Limit: &lim, Offset: &off, Labels: c.Labels})
		if resp == nil {
			return err
		}
		return dumpTable(resp.Data, resp.Error, err)
	case *structs.Faction:
		resp, err := conn.ListFactions(context.TODO(), &structs.ListFactionsRequest{World: c.World, Limit: &lim, Offset: &off, Labels: c.Labels})
		if resp == nil {
			return err
		}
		return dumpTable(resp.Data, resp.Error, err)
	case *structs.World:
		resp, err := conn.ListWorlds(context.TODO(), &structs.ListWorldsRequest{Limit: &lim, Offset: &off, Labels: c.Labels})
		if resp == nil {
			return err
		}
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
