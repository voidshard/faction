package main

import (
	"context"
	"fmt"

	"github.com/voidshard/faction/pkg/structs"
)

type cliDeleteCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Object struct {
		Name string   `positional-arg-name:"object" description:"Object to delete"`
		Id   []string `positional-arg-name:"id" description:"ID of object to delete"`
	}
}

func (c *cliDeleteCmd) Execute(args []string) error {
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

	checkErr := func(err error, errResp *structs.Error) error {
		if err != nil {
			return err
		}
		if errResp != nil {
			return fmt.Errorf("error: %s", errResp.Message)
		}
		return nil
	}

	switch obj.(type) {
	case *structs.Faction:
		for _, id := range c.Object.Id {
			resp, err := conn.DeleteFaction(context.TODO(), &structs.DeleteFactionRequest{World: c.World, Id: id})
			err = checkErr(err, resp.Error)
			if err != nil {
				return err
			}
		}
	case *structs.World:
		for _, id := range c.Object.Id {
			resp, err := conn.DeleteWorld(context.TODO(), &structs.DeleteWorldRequest{Id: id})
			err = checkErr(err, resp.Error)
			if err != nil {
				return err
			}
		}
	case *structs.Actor:
		for _, id := range c.Object.Id {
			resp, err := conn.DeleteActor(context.TODO(), &structs.DeleteActorRequest{Id: id})
			err = checkErr(err, resp.Error)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
