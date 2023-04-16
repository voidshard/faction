package db

import (
	"github.com/voidshard/faction/pkg/structs"
)

type Iterator[T interface{}] interface {
	Next() ([]T, error)
	HasNext() bool
}

type iterator[T interface{}] struct {
	token string
	get   func(string) ([]T, string, error)
	done  bool
}

func (i *iterator[T]) HasNext() bool {
	return !i.done
}

func (i *iterator[T]) Next() ([]T, error) {
	if i.done {
		return []T{}, nil
	}
	data, token, err := i.get(i.token)
	if err != nil {
		return nil, err
	}
	i.token = token
	i.done = token == ""
	return data, nil
}

func (f *FactionDB) IterPeople(filters ...*PersonFilter) Iterator[*structs.Person] {
	return &iterator[*structs.Person]{
		get: func(token string) ([]*structs.Person, string, error) {
			return f.People(token, filters...)
		},
	}
}
