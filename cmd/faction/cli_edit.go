package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/voidshard/faction/pkg/client"
	"github.com/voidshard/faction/pkg/structs"
)

type cliEditCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Object struct {
		Name string `positional-arg-name:"object" description:"Object to get"`
		Id   string `positional-arg-name:"id" description:"ID of object to get"`
	} `positional-args:"true" required:"true"`
}

func (c *cliEditCmd) Execute(args []string) error {
	obj := validObject(c.Object.Name)
	if obj == nil {
		return invalidObjectError(c.Object.Name)
	}

	_, isWorld := obj.(*structs.World)
	if isWorld {
		c.Object.Id = toWorldId(c.Object.Id)[0]
	} else if !isWorld && c.World == "" {
		return fmt.Errorf("world must be set for %s", c.Object.Name)
	}
	c.World = toWorldId(c.World)[0]

	conn, err := client.New(c.Host, c.Port, c.IdleTimeout, c.ConnTimeout)
	if err != nil {
		return err
	}

	// read object from API
	objs, err := getObjects(conn, c.World, obj, []string{c.Object.Id})
	if err != nil {
		return err
	}
	if len(objs) != 1 {
		return fmt.Errorf("expected 1 object, got %d", len(objs))
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

	// read the file back in
	edited, err := readObjectsFromFile(obj, []string{f.Name()})
	if err != nil {
		return err
	}

	// update the object
	return updateObjects(conn, c.World, obj, edited)
}
