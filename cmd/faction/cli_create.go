package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/client"
	"github.com/voidshard/faction/pkg/kind"
)

type cliCreateCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Files []string `short:"f" long:"file" description:"File(s) to read object(s) from"`
}

func (c *cliCreateCmd) Execute(args []string) error {
	conn, err := client.New(client.NewConfig())
	if err != nil {
		return err
	}
	return applyYamlUpdate(conn, c.World, c.Files)
}

func applyYamlUpdate(conn *client.Client, world string, files []string) error {
	toWrite, err := readObjectsFromFile(files)
	if err != nil {
		return err
	}
	// set world from cli if needed
	for _, v := range toWrite {
		if kind.IsGlobal(v.GetKind()) {
			// if the object doesn't need a world then it's fine as is
			continue
		} else if v.GetWorld() == "" && world != "" {
			// If object doesn't have a world but we have one from the cli, set it
			v.SetWorld(world)
		} else if v.GetWorld() != "" && world == "" {
			// If the object has a world set and we weren't given one from the cli
			continue
		} else if v.GetWorld() == world {
			// If the object has a world set and we were given one from the cli
			continue
		} else {
			// Either; we have two different conflicting world ids (object vs cli global)
			// or we have an object that requires a world and we've no default world var set
			return fmt.Errorf("world required but is not set or ambiguous", v.GetWorld(), world)
		}
	}

	return conn.Set(toWrite)

}
