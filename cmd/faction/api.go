package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/internal/search"
	"github.com/voidshard/faction/internal/service/api"
)

const (
	docApi = `Run the API server`
)

type optDatabase struct {
	Db_Host     string `long:"db-host" env:"DB_HOST" description:"Database host" default:"localhost"`
	Db_Port     int    `long:"db-port" env:"DB_PORT" description:"Database port" default:"27017"`
	Db_Username string `long:"db-username" env:"DB_USERNAME" description:"Database user" default:"admin"`
	Db_Password string `long:"db-password" env:"DB_PASSWORD" description:"Database password" default:"admin"`
	Db_Database string `long:"db-database" env:"DB_DATABASE" description:"Database name" default:"faction"`
}

type optSearchbase struct {
	Sb_Address       string        `long:"sb-address" env:"SB_ADDRESS" description:"Search host, accepts comma delimited" default:"http://localhost:9200"`
	Sb_Username      string        `long:"sb-username" env:"SB_USERNAME" description:"Database user" default:"admin"`
	Sb_Password      string        `long:"sb-password" env:"SB_PASSWORD" description:"Database password" default:"Opensearch123!"`
	Sb_FlushInterval time.Duration `long:"sb-flush-interval" env:"SB_FLUSH_INTERVAL" description:"Flush interval for searchbase" default:"5s"`
	Sb_WriteRoutines int           `long:"sb-write-workers" env:"SB_WRITE_WORKERS" description:"Number of routines to use for searchbase" default:"2"`
}

type optQueue struct {
	Q_Host     string `long:"q-host" env:"Q_HOST" description:"Queue host" default:"localhost"`
	Q_Port     int    `long:"q-port" env:"Q_PORT" description:"Queue port" default:"5672"`
	Q_Username string `long:"q-username" env:"Q_USERNAME" description:"Queue user" default:"admin"`
	Q_Password string `long:"q-password" env:"Q_PASSWORD" description:"Queue password" default:"admin"`

	Q_ReplyMaxTime  time.Duration `long:"q-reply-max-time" env:"Q_REPLY_MAX_TIME" description:"Maximum time to wait for a reply" default:"60s"`
	Q_ReplyRoutines int           `long:"q-reply-routines" env:"Q_REPLY_ROUTINES" description:"Number of routines to use for replies" default:"5"`
}

type optsAPI struct {
	optGeneral
	optDatabase
	optSearchbase
	optQueue

	// MaxMessageAge is the maximum age of a message before it is considered stale
	MaxMessageAge time.Duration `env:"MAX_MESSAGE_AGE" long:"max-message-age" description:"Maximum age of a message before it is considered stale" default:"10m"`
	Routines      int           `env:"ROUTINES" long:"routines" description:"Number of routines to use for processing" default:"5"`
	Port          int           `env:"PORT" long:"port" description:"Port to listen on" default:"8080"`
	FlushSearch   bool          `env:"FLUSH_SEARCH" long:"flush-search" description:"Wait for writes to searchbase before returning API writes (slow)"`
}

func (c *optsAPI) Execute(args []string) error {
	if c.Debug {
		os.Setenv("DEBUG", "true")
		log.SetGlobalLevel()
	}

	// connect to the database
	database, err := db.NewMongo(&db.MongoConfig{
		Host:     c.Db_Host,
		Port:     c.Db_Port,
		Username: c.Db_Username,
		Password: c.Db_Password,
		Database: c.Db_Database,
	})
	log.Info().Err(err).Str("host", c.Db_Host).Int("port", c.Db_Port).Str("username", c.Db_Username).Str("database", c.Db_Database).Msg("database connection")
	if err != nil {
		return err
	}

	// connect to the queue
	qu := queue.NewRabbitQueue(&queue.RabbitConfig{
		Host:     c.Q_Host,
		Port:     c.Q_Port,
		Username: c.Q_Username,
		Password: c.Q_Password,
	})
	log.Info().Str("host", c.Q_Host).Int("port", c.Q_Port).Str("username", c.Q_Username).Msg("queue connection")

	// connect to search
	sb, err := search.NewOpensearch(&search.OpensearchConfig{
		Address:       c.Sb_Address,
		Username:      c.Sb_Username,
		Password:      c.Sb_Password,
		FlushInterval: c.Sb_FlushInterval,
		WriteRoutines: c.Sb_WriteRoutines,
	})
	log.Info().Err(err).Str("address", c.Sb_Address).Str("username", c.Sb_Username).Msg("search connection")
	if err != nil {
		return err
	}

	// setup the API server
	srv := api.NewService(&api.Config{
		Routines:      c.Routines,
		MaxMessageAge: c.MaxMessageAge,
		FlushSearch:   c.FlushSearch,
	}, database, qu, sb)
	server := api.NewServer(srv)

	// basic signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Info().Str("signal", sig.String()).Msg("signal received")
		server.Stop()
	}()

	return server.Serve(c.Port)
}
