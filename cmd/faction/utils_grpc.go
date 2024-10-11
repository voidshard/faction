package main

import (
	"fmt"
	"time"

	"github.com/voidshard/faction/pkg/structs"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

func dial(host string, port int, idle, conntimeout time.Duration) (*grpc.ClientConn, error) {
	return grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithIdleTimeout(idle),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: conntimeout,
			Backoff:           backoff.DefaultConfig,
		}),
		//grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
}

func newClient(host string, port int, idle, conntimeout time.Duration) (structs.APIClient, error) {
	conn, err := dial(host, port, idle, conntimeout)
	if err != nil {
		return nil, err
	}
	return structs.NewAPIClient(conn), nil
}
