package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/client"
	"github.com/voidshard/faction/pkg/kind"
)

type cliDeleteCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Object struct {
		Kind string   `positional-arg-name:"object" description:"Object to get"`
		Id   []string `positional-arg-name:"id" description:"ID of object to delete"`
	} `positional-args:"true" required:"true"`
}

func (c *cliDeleteCmd) Execute(args []string) error {
	if c.Object.Kind == "" {
		return fmt.Errorf("invalid object kind %s", c.Object.Kind)
	}

	if !kind.IsGlobal(c.Object.Kind) && c.World == "" {
		return fmt.Errorf("world must be set for non-global objects")
	}

	conn, err := client.New(client.NewConfig())
	if err != nil {
		return err
	}

	return conn.Delete(c.Object.Kind, c.World, c.Object.Id)
}
