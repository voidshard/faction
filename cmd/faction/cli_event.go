package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/client"
)

type cliEventCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Object struct {
		Kind string `positional-arg-name:"object" description:"Object to get"`
		Id   string `positional-arg-name:"id" description:"ID of object to get"`
	} `positional-args:"true" required:"true"`

	Controller string `long:"controller" short:"c" description:"Controller to use"`

	ToTick uint64 `long:"to-tick" short:"t" description:"Tick to defer to"`
	ByTick uint64 `long:"by-tick" short:"b" description:"Tick to defer by (relative to current)"`
}

func (c *cliEventCmd) Execute(args []string) error {
	if c.ToTick == 0 && c.ByTick == 0 {
		return fmt.Errorf("must set either --to-tick or --by-tick")
	}
	if c.ToTick != 0 && c.ByTick != 0 {
		return fmt.Errorf("only one of --to-tick or --by-tick can be set")
	}

	c.Object.Kind = validKind(c.Object.Kind)
	if c.Object.Kind == "" {
		return fmt.Errorf("invalid object kind %s", c.Object.Kind)
	}

	if c.World == "" {
		return fmt.Errorf("world must be set for non-global objects")
	}

	conn, err := client.New(client.NewConfig())
	if err != nil {
		return err
	}

	// setup the request
	def := conn.Defer(c.Object.Kind, c.World, c.Object.Id)
	if c.Controller != "" {
		def = def.Controller(c.Controller)
	}
	if c.ToTick != 0 {
		def = def.ToTick(c.ToTick)
	}
	if c.ByTick != 0 {
		def = def.ByTick(c.ByTick)
	}

	toTick, err := def.Do()
	if err != nil {
		return err
	}
	fmt.Printf(
		"{Kind: %s, World: %s, Id: %s, Controller: %s} => tick: %d\n",
		c.Object.Kind, c.World, c.Object.Id, c.Controller, toTick,
	)
	return nil
}
