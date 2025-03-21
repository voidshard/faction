package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/voidshard/faction/pkg/client"
	"github.com/voidshard/faction/pkg/kind"
)

type cliEditCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Object struct {
		Kind string `positional-arg-name:"object" description:"Object to get"`
		Id   string `positional-arg-name:"id" description:"ID of object to get"`
	} `positional-args:"true" required:"true"`
}

func (c *cliEditCmd) Execute(args []string) error {
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

	// get object
	objs, err := conn.Get().Ids([]string{c.Object.Id}).World(c.World).Do(c.Object.Kind)
	if err != nil {
		return err
	}
	if len(objs) != 1 {
		return fmt.Errorf("object not found")
	}

	// write to temp file
	yamlData, err := dumpYaml(objs)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp("", "tmp-")
	if err != nil {
		return err
	}
	_, err = f.Write(yamlData)
	if err != nil {
		return err
	}
	f.Close()

	// calculate a hash of the file before editing
	preHash, err := calculateFileHash(f.Name())
	if err != nil {
		return err
	}

	// open the file in the editor and wait for it to close
	editcmd := os.Getenv("EDITOR")
	if editcmd == "" {
		editcmd = "vi" // what else would you use?
	}

	editor := exec.Command(editcmd, f.Name())
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	err = editor.Run()
	if err != nil {
		return err
	}

	// calculate a hash of the file after editing
	postHash, err := calculateFileHash(f.Name())
	if err != nil {
		return err
	}
	if preHash == postHash {
		fmt.Println("No changes detected")
		return nil
	}

	// update the object
	return applyYamlUpdate(conn, c.World, []string{f.Name()})
}
