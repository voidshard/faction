package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/internal/service/api"
)

const (
	docApi = `Run the API server`
)

type optsDatabase struct {
	Db_Host     string `long:"db-host" env:"DB_HOST" description:"Database host" default:"localhost"`
	Db_Port     int    `long:"db-port" env:"DB_PORT" description:"Database port" default:"27017"`
	Db_Username string `long:"db-username" env:"DB_USERNAME" description:"Database user" default:"admin"`
	Db_Password string `long:"db-password" env:"DB_PASSWORD" description:"Database password" default:"admin"`
	Db_Database string `long:"db-database" env:"DB_DATABASE" description:"Database name" default:"faction"`
}

type optsQueue struct {
	Q_Host     string `long:"q-host" env:"Q_HOST" description:"Queue host" default:"localhost"`
	Q_Port     int    `long:"q-port" env:"Q_PORT" description:"Queue port" default:"5672"`
	Q_Username string `long:"q-username" env:"Q_USERNAME" description:"Queue user" default:"admin"`
	Q_Password string `long:"q-password" env:"Q_PASSWORD" description:"Queue password" default:"admin"`
}

type optsAPI struct {
	optsGeneral
	optsDatabase
	optsQueue

	// MaxMessageAge is the maximum age of a message before it is considered stale
	MaxMessageAge time.Duration `long:"max-message-age" description:"Maximum age of a message before it is considered stale" default:"10m"`
	Routines      int           `long:"routines" description:"Number of routines to use for processing" default:"5"`
	Port          int           `long:"port" description:"Port to listen on" default:"8080"`
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
	log.Info().Err(err).Str("host", c.Q_Host).Int("port", c.Q_Port).Str("username", c.Q_Username).Msg("queue connection")

	// setup the API server
	srv := api.NewService(&api.Config{Routines: c.Routines, MaxMessageAge: c.MaxMessageAge}, database, qu)
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
