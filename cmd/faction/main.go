package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type optsGeneral struct {
	Debug bool `long:"debug" env:"DEBUG" description:"Enable debug mode"`
}

type optsDatabase struct {
	DatabaseDriver string `long:"database-driver" env:"DATABASE_DRIVER" description:"Database driver" choice:"postgres" choice:"sqlite3" default:"postgres"`
	DatabaseURL    string `long:"database-url" env:"DATABASE_URL" description:"Database connection string" default:"postgres://faction:faction@localhost:5432/faction?sslmode=disable&search_path=faction"`
}

type optsQueue struct {
	QueueDriver      string `long:"queue-driver" env:"QUEUE_DRIVER" description:"Queue driver" choice:"igor" choice:"inmemory" default:"igor"`
	QueueURL         string `long:"queue-url" env:"QUEUE_URL" description:"Queue connection string" default:"igor-redis:6379"`
	QueueDatabaseURL string `long:"queue-database-url" env:"QUEUE_DATABASE_URL" description:"Queue database connection string" default:"postgres://igor:igor@localhost:5432/igor?sslmode=disable"`
}

var parser = flags.NewParser(nil, flags.Default)

func init() {
	parser.AddCommand("worker", "Run faction background worker", docWorker, &optsWorker{})
	parser.AddCommand("migrate", "Postgres Migration Operations", docMigrate, &optsMigrate{})
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
