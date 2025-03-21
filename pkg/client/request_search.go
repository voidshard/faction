package client

import (
	"github.com/voidshard/faction/pkg/kind"
	"github.com/voidshard/faction/pkg/structs/api"
	v1 "github.com/voidshard/faction/pkg/structs/v1"
)

type Operation string

const (
	Equal       Operation = "eq"
	LessThan    Operation = "lt"
	GreaterThan Operation = "gt"
)

type searchBuilder struct {
	client *Client
	world  string
	Req    *api.SearchRequest
}

func (s *searchBuilder) Do() ([]v1.Object, error) {
	resp, err := s.client.search(s.world, s.Req)
	if err != nil {
		return nil, err
	}
	objects := []v1.Object{}
	for _, d := range resp.Data {
		obj, err := kind.New(s.Req.Kind, d)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

func (s *searchBuilder) Worlds() ([]*v1.World, error) {
	resp, err := s.client.search(s.world, s.Req)
	if err != nil {
		return nil, err
	}
	return toObject(resp.Data, kindWorld, &v1.World{})
}

func (s *searchBuilder) Actors() ([]*v1.Actor, error) {
	resp, err := s.client.search(s.world, s.Req)
	if err != nil {
		return nil, err
	}
	return toObject(resp.Data, kindActor, &v1.Actor{})
}

func toObject[T v1.Object](data []interface{}, k string, o T) ([]T, error) {
	objects := []T{}
	for _, d := range data {
		obj, err := kind.New(k, d)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj.(T))
	}
	return objects, nil
}

func (s *searchBuilder) RandomWeight(i float64) *searchBuilder {
	s.Req.RandomWeight = i
	return s
}

func (s *searchBuilder) All(field string, value interface{}, op ...Operation) *searchBuilder {
	s.Req.All = append(s.Req.All, toMatch(field, value, op))
	return s
}

func (s *searchBuilder) Any(field string, value interface{}, op ...Operation) *searchBuilder {
	s.Req.Any = append(s.Req.Any, toMatch(field, value, op))
	return s
}

func (s *searchBuilder) Not(field string, value interface{}, op ...Operation) *searchBuilder {
	s.Req.Not = append(s.Req.Not, toMatch(field, value, op))
	return s
}

func (s *searchBuilder) Score(field string, value interface{}, weight float64, op ...Operation) *searchBuilder {
	match := toMatch(field, value, op)
	s.Req.Score = append(s.Req.Score, v1.Score{
		Match:  match,
		Weight: weight,
	})
	return s
}

func toMatch(field string, value interface{}, op []Operation) v1.Match {
	comparison := Equal
	if len(op) > 0 {
		comparison = op[0]
	}
	return v1.Match{
		Field: field,
		Value: value,
		Op:    string(comparison),
	}
}
