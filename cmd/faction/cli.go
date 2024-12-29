package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	docCli = `CLI for operations on the API`
)

type optCliConn struct {
	Host        string        `long:"host" env:"HOST" description:"API host" default:"localhost"`
	Port        int           `long:"port" env:"PORT" description:"API port" default:"5000"`
	IdleTimeout time.Duration `long:"idle-timeout" env:"IDLE_TIMEOUT" description:"Idle timeout" default:"30s"`
	ConnTimeout time.Duration `long:"conn-timeout" env:"CONN_TIMEOUT" description:"Connection timeout" default:"5s"`
}

type optCliGlobal struct {
	World string `long:"world" short:"w" env:"WORLD" description:"World to use"`
}

type optsCli struct {
	Get    cliGetCmd    `command:"get" description:"Get objects"`
	List   cliListCmd   `command:"list" description:"List objects"`
	Create cliCreateCmd `command:"create" description:"Create objects"`
	Delete cliDeleteCmd `command:"delete" description:"Delete objects"`
	Watch  cliWatchCmd  `command:"watch" description:"Watch objects"`
	Edit   cliEditCmd   `command:"edit" description:"Edit objects"`

	Help cliHelpCmd `command:"help" description:"Help about available objects"`
}

type cliHelpCmd struct {
}

func (c *cliHelpCmd) Execute(args []string) error {
	objects := map[string][]string{}
	for k, v := range shortNames {
		abbreviations, ok := objects[v]
		if !ok {
			abbreviations = []string{}
		}
		objects[v] = append(abbreviations, k)
	}

	fmt.Println("Objects:")
	for k := range objects {
		abbreviations, ok := objects[k]
		abbr := ""
		if ok {
			abbr = fmt.Sprintf(" (abbreviations: %s)", strings.Join(abbreviations, ", "))
		}
		fmt.Printf("\t%s%s\n", k, abbr)
		extra, ok := help[k]
		if ok {
			fmt.Printf("\t\t%s\n", extra)
		}
	}

	return nil
}
