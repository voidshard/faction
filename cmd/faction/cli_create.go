package main

import (
	"context"
	"fmt"

	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/structs"
)

type cliCreateCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Object struct {
		Name string `positional-arg-name:"object" description:"Object to create"`
	} `positional-args:"true" required:"true"`

	Files []string `short:"f" long:"file" description:"File(s) to read object(s) from"`
}

func (c *cliCreateCmd) Execute(args []string) error {
	obj := validObject(c.Object.Name)
	if obj == nil {
		return invalidObjectError(c.Object.Name)
	}

	toWrite, err := readObjectsFromFile(obj, c.Files)
	if err != nil {
		return err
	}

	cliGivenWorld := toWorldId(c.World)[0]
	byWorld := map[string][]structs.Object{}
	for _, v := range toWrite {
		if v.GetWorld() == "" { // object does not have world set
			if c.World == "" {
				return fmt.Errorf("world not set on either object or command line")
			}
			v.SetWorld(cliGivenWorld) // set world to the one given on the command line
		} else { // object has world set
			if c.World == "" {
				// this is fine, object has a world set and command line doesn't specify one
			} else if v.GetWorld() != cliGivenWorld {
				return fmt.Errorf("world mismatch: %s != %s for object %s", v.GetWorld(), cliGivenWorld, v.GetId())
			}
		}
		byWorld[v.GetWorld()] = append(byWorld[v.GetWorld()], v)
	}

	conn, err := newClient(c.Host, c.Port, c.IdleTimeout, c.ConnTimeout)
	if err != nil {
		return err
	}

	for world, subobjects := range byWorld {
		err = updateObjects(conn, world, obj, subobjects)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateObjects(conn structs.APIClient, world string, obj structs.Object, objs []structs.Object) error {
	var (
		rerr *structs.Error
	)
	switch obj.(type) {
	case *structs.Actor:
		in := make([]*structs.Actor, len(objs))
		for i, v := range objs {
			in[i] = v.(*structs.Actor)
		}
		resp, err := conn.SetActors(context.TODO(), &structs.SetActorsRequest{World: world, Data: in})
		if err != nil {
			return err
		} else if resp == nil {
			return fmt.Errorf("nil response")
		}
		rerr = resp.Error
	case *structs.Faction:
		in := make([]*structs.Faction, len(objs))
		for i, v := range objs {
			in[i] = v.(*structs.Faction)
		}
		resp, err := conn.SetFactions(context.TODO(), &structs.SetFactionsRequest{World: world, Data: in})
		if err != nil {
			return err
		} else if resp == nil {
			return fmt.Errorf("nil response")
		}
		rerr = resp.Error
	case *structs.World:
		for _, v := range objs {
			resp, err := conn.SetWorld(context.TODO(), &structs.SetWorldRequest{Data: v.(*structs.World)})
			if err != nil {
				return err
			} else if resp == nil {
				return fmt.Errorf("nil response")
			} else if resp.Error != nil {
				return fmt.Errorf("error: %s", resp.Error.Message)
			}
		}
	}

	var err error
	if rerr != nil {
		err = fmt.Errorf("error: %s", rerr.Message)
	}
	log.Debug().Err(err).Msg("updateObjects response")
	return err
}
