package client

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/voidshard/faction/pkg/structs"
	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

func NewFromEnv() (structs.APIClient, error) {
	// read from env
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	idle := os.Getenv("IDLE_TIMEOUT")
	if idle == "" {
		idle = "10s"
	}
	conntimeout := os.Getenv("CONN_TIMEOUT")
	if conntimeout == "" {
		conntimeout = "5s"
	}

	// parse values
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	idleDuration, err := time.ParseDuration(idle)
	if err != nil {
		return nil, err
	}
	conntimeoutDuration, err := time.ParseDuration(conntimeout)
	if err != nil {
		return nil, err
	}

	return New(host, portInt, idleDuration, conntimeoutDuration)
}

func New(host string, port int, idle, conntimeout time.Duration) (structs.APIClient, error) {
	conn, err := rpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		rpc.WithTransportCredentials(insecure.NewCredentials()),
		rpc.WithIdleTimeout(idle),
		rpc.WithConnectParams(rpc.ConnectParams{
			MinConnectTimeout: conntimeout,
			Backoff:           backoff.DefaultConfig,
		}),
		//grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}
	return structs.NewAPIClient(conn), nil
}
