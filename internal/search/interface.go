package search

import (
	"context"

	"github.com/voidshard/faction/pkg/structs"
)

type Search interface {
	Index(ctx context.Context, world string, in []structs.Object, flush bool) error
	Delete(ctx context.Context, world, kind, id string) error
	// Find(ctx context.Context, world, kind string, query string, limit, offset int) ([]structs.Object, error)
}
