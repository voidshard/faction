package main

import (
	"context"
	"fmt"

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
		return invalidObjectError(c.Object.Name)
	}

	_, isWorld := obj.(*structs.World)
	if isWorld {
		c.Object.Id = toWorldId(c.Object.Id...)
	} else if !isWorld && c.World == "" {
		return fmt.Errorf("world must be set for %s", c.Object.Name)
	}
	c.World = toWorldId(c.World)[0]

	conn, err := newClient(c.Host, c.Port, c.IdleTimeout, c.ConnTimeout)
	if err != nil {
		return err
	}

	objs, err := getObjects(conn, c.World, obj, c.Object.Id)
	if err != nil {
		return err
	}

	yamlData, err := dumpYaml(objs)
	if yamlData != nil {
		fmt.Println(string(yamlData))
	}
	return err
}

func getObjects(client structs.APIClient, world string, obj marshalable, ids []string) ([]marshalable, error) {
	data := []marshalable{}
	switch obj.(type) {
	case *structs.Faction:
		resp, err := client.Factions(context.TODO(), &structs.GetFactionsRequest{World: world, Ids: ids})
		if err != nil {
			return nil, err
		} else if resp == nil {
			return nil, fmt.Errorf("nil response")
		} else if resp.Error != nil {
			return nil, fmt.Errorf("error: %s", resp.Error.Message)
		}
		for _, d := range resp.Data {
			data = append(data, d)
		}
	case *structs.Actor:
		resp, err := client.Actors(context.TODO(), &structs.GetActorsRequest{World: world, Ids: ids})
		if err != nil {
			return nil, err
		} else if resp == nil {
			return nil, fmt.Errorf("nil response")
		} else if resp.Error != nil {
			return nil, fmt.Errorf("error: %s", resp.Error.Message)
		}
		for _, d := range resp.Data {
			data = append(data, d)
		}
	case *structs.World:
		resp, err := client.Worlds(context.TODO(), &structs.GetWorldsRequest{Ids: ids})
		if err != nil {
			return nil, err
		} else if resp == nil {
			return nil, fmt.Errorf("nil response")
		} else if resp.Error != nil {
			return nil, fmt.Errorf("error: %s", resp.Error.Message)
		}
		for _, d := range resp.Data {
			data = append(data, d)
		}
	}
	return data, nil
}
