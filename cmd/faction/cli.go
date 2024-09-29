package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/voidshard/faction/pkg/structs"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	docCli = `CLI for operations on the API`
)

var (
	objects = map[string]structs.Object{ // map of objects we can interact with
		"world":   &structs.World{},
		"actor":   &structs.Actor{},
		"faction": &structs.Faction{},
	}
	help = map[string]string{ // help text for each object
		"world":   "Worlds are the top level object in the system, every other object exists within a world.",
		"actor":   "Actors are individual entities that can interact & form factions.",
		"faction": "Factions are groups of actors that can work together, form common ground(s) or simply work against other factions.",
	}
	shortNames = map[string]string{ // allows short hand names because we're lazy
		"wo": "world",
		"ac": "actor",
		"fa": "faction",
	}
)

type optCliConn struct {
	Host        string        `long:"host" env:"HOST" description:"API host" default:"localhost"`
	Port        int           `long:"port" env:"PORT" description:"API port" default:"5000"`
	IdleTimeout time.Duration `long:"idle-timeout" description:"Idle timeout" default:"30s"`
	ConnTimeout time.Duration `long:"conn-timeout" description:"Connection timeout" default:"5s"`
}

type optCliGlobal struct {
	World string `long:"world" short:"w" env:"WORLD" description:"World to use"`
}

type optsCli struct {
	Get    cliGetCmd    `command:"get" description:"Get objects"`
	List   cliListCmd   `command:"list" description:"List objects"`
	Create cliCreateCmd `command:"create" description:"Create objects"`
	Delete cliDeleteCmd `command:"delete" description:"Delete objects"`

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

func newClient(host string, port int, idle, conntimeout time.Duration) (structs.APIClient, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithIdleTimeout(idle),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: conntimeout,
			Backoff:           backoff.DefaultConfig,
		}),
	)
	if err != nil {
		return nil, err
	}
	return structs.NewAPIClient(conn), nil
}

func validObject(name string) structs.Object {
	name = strings.ToLower(name)

	longname, ok := shortNames[name]
	if ok {
		name = longname
	}

	obj, _ := objects[name]
	return obj
}

func invalidObjectError() error {
	valid := []string{}
	for k := range objects {
		valid = append(valid, k)
	}
	return fmt.Errorf("Invalid object. Valid names: %s", strings.Join(valid, ", "))
}
