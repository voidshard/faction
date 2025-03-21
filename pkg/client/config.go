package client

import (
	"os"
	"strconv"

	"github.com/voidshard/faction/pkg/util/log"
)

type Config struct {
	Host string
	Port int
}

func NewConfig() *Config {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		portInt = 5000
		log.Warn().Err(err).Msg("Failed to parse port, using default")
	}

	return &Config{Host: host, Port: portInt}
}
