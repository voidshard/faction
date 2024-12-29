package main

const (
	docController = `Stand up a controller to manage object(s) changes`
)

type optLocker struct {
	Lk_Host string `long:"lk-host" env:"LK_HOST" description:"Locker host" default:"localhost"`
	Lk_Port int    `long:"lk-port" env:"LK_PORT" description:"Locker port" default:"6379"`
}

/**
	// connect to the locker
	locker, err := db.NewRedisLocker(&db.RedisConfig{
	    Host: c.Lk_Host,
	    Port: c.Lk_Port,
	})
	log.Info().Err(err).Str("host", c.Lk_Host).Int("port", c.Lk_Port).Msg("locker connection")
	if err != nil {
	    return err
	}
**/

type optsController struct {
}
