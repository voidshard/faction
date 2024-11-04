package client

import (
	"fmt"
	"time"

	"github.com/voidshard/faction/pkg/structs"
	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

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
