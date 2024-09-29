package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/voidshard/faction/pkg/structs"
	"gopkg.in/yaml.v3"
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
		return invalidObjectError()
	}

	_, isWorld := obj.(*structs.World)
	if !isWorld && c.World == "" {
		return fmt.Errorf("world must be set for %s", c.Object.Name)
	}

	toWrite, err := readObjectsFromFile(obj, c.Files)
	if err != nil {
		return err
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
	case *structs.Actor:
		in := make([]*structs.Actor, len(toWrite))
		for i, v := range toWrite {
			in[i] = v.(*structs.Actor)
		}
		resp, err := conn.SetActors(context.TODO(), &structs.SetActorsRequest{World: c.World, Data: in})
		return checkErr(err, resp.Error)
	case *structs.Faction:
		in := make([]*structs.Faction, len(toWrite))
		for i, v := range toWrite {
			in[i] = v.(*structs.Faction)
		}
		resp, err := conn.SetFactions(context.TODO(), &structs.SetFactionsRequest{World: c.World, Data: in})
		return checkErr(err, resp.Error)
	case *structs.World:
		for i, v := range toWrite {
			_, err = conn.SetWorld(context.TODO(), &structs.SetWorldRequest{Data: v.(*structs.World)})
			err = checkErr(err, nil)
			if err != nil {
				return fmt.Errorf("error on object %d: %w", i, err)
			}
		}
	}

	return nil
}

func readObjectsFromFile[T structs.Object](objToRead T, files []string) ([]T, error) {
	found := []T{}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		for _, chunk := range bytes.Split(data, []byte("\n---\n")) {
			obj := new(T)
			err = yaml.Unmarshal(chunk, obj)
			if err != nil {
				return nil, err
			}
			found = append(found, *obj)
		}
	}
	return found, nil
}
