package db

import (
	"context"

	v1 "github.com/voidshard/faction/pkg/structs/v1"
)

type Database interface {
	Get(c context.Context, world, kind string, id []string, out interface{}) error
	List(c context.Context, world, kind string, labels map[string]string, limit, offset int64, out interface{}) error
	Set(c context.Context, world, etag string, in []v1.Object) (*Result, error)
	Delete(c context.Context, world, kind string, id string) error

	Close()
}

// Result returns information about a batch write operation.
//
// On batch writes we use optimistic locking with our Etags. So rows are only written if the Etag
// supplied matches the current Etag in the database. The database then will assign a new Etag
// and we return the row ID -> new Etag mapping here.
type Result struct {
	// Written maps ID -> new Etag of written rows
	Written map[string]string
}

func NewResult() *Result {
	return &Result{
		Written: map[string]string{},
	}
}

func (r *Result) Merge(other *Result) {
	for k, v := range other.Written {
		r.Written[k] = v
	}
}
