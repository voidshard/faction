package main

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/voidshard/faction/pkg/structs"
)

type cliGetCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Object struct {
		Name string   `positional-arg-name:"object" description:"Object to get"`
		Id   []string `positional-arg-name:"id" description:"ID of object to get"`
	} `positional-args:"true" required:"true"`
}

func (c *cliGetCmd) Execute(args []string) error {
	obj := validObject(c.Object.Name)
	if obj == nil {
		return invalidObjectError()
	}

	_, isWorld := obj.(*structs.World)
	if !isWorld && c.World == "" {
		return fmt.Errorf("world must be set for %s", c.Object.Name)
	}

	conn, err := newClient(c.Host, c.Port, c.IdleTimeout, c.ConnTimeout)
	if err != nil {
		return err
	}

	switch obj.(type) {
	case *structs.Faction:
		resp, err := conn.Factions(context.TODO(), &structs.GetFactionsRequest{World: c.World, Ids: c.Object.Id})
		return dumpYaml(resp.Data, resp.Error, err)
	case *structs.Actor:
		resp, err := conn.Actors(context.TODO(), &structs.GetActorsRequest{World: c.World, Ids: c.Object.Id})
		return dumpYaml(resp.Data, resp.Error, err)
	case *structs.World:
		resp, err := conn.Worlds(context.TODO(), &structs.GetWorldsRequest{Ids: c.Object.Id})
		return dumpYaml(resp.Data, resp.Error, err)
	}

	return nil
}

func dumpYaml[T structs.Object](data []T, err *structs.Error, connErr error) error {
	if connErr != nil {
		return connErr
	}
	if err != nil {
		return fmt.Errorf("error: %s", err.Message)
	}

	for i, d := range data {
		b, err := yaml.Marshal(d)
		if err != nil {
			return err
		}
		if i > 0 && i < len(data) {
			fmt.Println("---")
		}
		fmt.Println(strings.TrimSpace(string(b)))
	}

	return nil
}
