package client

import (
	"github.com/voidshard/faction/pkg/kind"
	"github.com/voidshard/faction/pkg/structs/api"
	v1 "github.com/voidshard/faction/pkg/structs/v1"
)

type getBuilder struct {
	client *Client
	Req    *api.GetRequest
}

func (r *getBuilder) Do(k string) ([]v1.Object, error) {
	if r.Req.Ids != nil && len(r.Req.Ids) > 0 {
		r.Req.Labels = nil
	}
	resp, err := r.client.doGet(k, r.Req)
	if err != nil {
		return nil, err
	}
	objects := []v1.Object{}
	for _, d := range resp.Data {
		obj, err := kind.New(k, d)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, err
}

func (r *getBuilder) Worlds() ([]*v1.World, error) {
	resp, err := r.client.doGet(kindWorld, r.Req)
	if err != nil {
		return nil, err
	}
	return toObject(resp.Data, kindWorld, &v1.World{})
}

func (r *getBuilder) Actors() ([]*v1.Actor, error) {
	resp, err := r.client.doGet(kindActor, r.Req)
	if err != nil {
		return nil, err
	}
	return toObject(resp.Data, kindActor, &v1.Actor{})
}

func (r *getBuilder) Ids(ids []string) *getBuilder {
	r.Req.Ids = append(r.Req.Ids, ids...)
	return r
}

func (r *getBuilder) Limit(limit int64) *getBuilder {
	if limit > 0 {
		r.Req.Limit = limit
	}
	return r
}

func (r *getBuilder) Offset(offset int64) *getBuilder {
	if offset > 0 {
		r.Req.Offset = offset
	}
	return r
}

func (r *getBuilder) Labels(labels map[string]string) *getBuilder {
	for k, v := range labels {
		r.Req.Labels[k] = v
	}
	return r
}

func (r *getBuilder) World(world string) *getBuilder {
	r.Req.World = world
	return r
}
