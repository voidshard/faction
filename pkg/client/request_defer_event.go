package client

import "github.com/voidshard/faction/pkg/structs/api"

type deferEventBuilder struct {
	client *Client
	Req    *api.DeferEventRequest
}

func (b *deferEventBuilder) ByTick(tick uint64) *deferEventBuilder {
	b.Req.ByTick = tick
	return b
}

func (b *deferEventBuilder) ToTick(tick uint64) *deferEventBuilder {
	b.Req.ToTick = tick
	return b
}

func (b *deferEventBuilder) Controller(controller string) *deferEventBuilder {
	b.Req.Controller = controller
	return b
}

func (b *deferEventBuilder) Do() (uint64, error) {
	resp, err := b.client.doDefer(b.Req)
	if err != nil {
		return 0, err
	}
	return resp.ToTick, nil
}
