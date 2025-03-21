package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/client"
	"github.com/voidshard/faction/pkg/kind"
)

type cliGetCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Object struct {
		Kind string `positional-arg-name:"object" description:"Object to get"`
	} `positional-args:"true" required:"true"`

	Limit  int64 `long:"limit" default:"100" description:"Limit number of results"`
	Offset int64 `short:"o" long:"offset" default:"0" description:"Offset results"`

	Labels map[string]string `short:"l" long:"labels" description:"Filter by labels"`
}

func (c *cliGetCmd) Execute(args []string) error {
	c.Object.Kind = validKind(c.Object.Kind)
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

	objs, err := conn.Get().Ids(args).Limit(c.Limit).Offset(c.Offset).Labels(c.Labels).World(c.World).Do(c.Object.Kind)
	if err != nil {
		return err
	}

	yamlData, err := dumpYaml(objs)
	if yamlData != nil {
		fmt.Println(string(yamlData))
	}
	return err
}
