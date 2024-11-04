package controller

import (
	"github.com/voidshard/faction/pkg/structs"
)

type subscriber struct {
	stream structs.API_OnChangeClient
	kill   chan bool
}

func newSubscriber(stream structs.API_OnChangeClient) *subscriber {
	return &subscriber{stream: stream, kill: make(chan bool)}
}
