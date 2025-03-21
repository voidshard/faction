package client

import (
	"fmt"
	"net/url"

	"github.com/voidshard/faction/pkg/structs/api"
)

type watchBuilder struct {
	client *Client
	Req    *api.StreamEvents
}

func (b *watchBuilder) Do() (*EventStream, error) {
	v := url.Values{}
	v.Set("world", b.Req.World)
	v.Set("kind", b.Req.Kind)
	v.Set("controller", b.Req.Controller)
	v.Set("id", b.Req.Id)
	v.Set("queue", b.Req.Queue)
	return newEventStream(&url.URL{
		Scheme:   "ws",
		Host:     fmt.Sprintf("%s:%d", b.client.cfg.Host, b.client.cfg.Port),
		Path:     "/v1/event",
		RawQuery: v.Encode(),
	})
}

func (b *watchBuilder) World(world string) *watchBuilder {
	b.Req.World = world
	return b
}

func (b *watchBuilder) Kind(kind string) *watchBuilder {
	b.Req.Kind = kind
	return b
}

func (b *watchBuilder) Controller(controller string) *watchBuilder {
	b.Req.Controller = controller
	return b
}

func (b *watchBuilder) Id(id string) *watchBuilder {
	b.Req.Id = id
	return b
}

func (b *watchBuilder) Queue(queue string) *watchBuilder {
	b.Req.Queue = queue
	return b
}
