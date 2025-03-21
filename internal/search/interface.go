package search

import (
	"context"

	v1 "github.com/voidshard/faction/pkg/structs/v1"
)

type Search interface {
	// Index adds an object to the search index.
	Index(ctx context.Context, world string, in []v1.Object, flush bool) error

	// Delete removes an object from the search index.
	Delete(ctx context.Context, world, kind, id string) error

	// Find returns a list of IDs that match the given query.
	// IDs are returned in order of relevance based on given scoring, most relevant first.
	Find(ctx context.Context, world string, q *v1.Query) ([]string, error)
}
