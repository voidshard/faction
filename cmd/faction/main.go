package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type optGeneral struct {
	Debug bool `long:"debug" env:"DEBUG" description:"Enable debug logging"`
}

var cmdAPI optsAPI
var cmdCli optsCli
var parser = flags.NewParser(nil, flags.Default)

func init() {
	parser.AddCommand("api", "Run API Server", docApi, &cmdAPI)
	parser.AddCommand("cli", "Run the CLI", docCli, &cmdCli)
}

func main() {
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}
