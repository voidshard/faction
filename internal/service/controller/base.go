package controller

import (
	"context"

	"github.com/voidshard/faction/pkg/structs"
)

type Search interface {
	Find(ctx context.Context, world, kind string, query string, limit, offset int) ([]structs.Object, error)
}

type Controller interface {
	Init(ctx context.Context, api structs.APIClient, sh Search) error
	Reconcile(ctx context.Context, ch *structs.Change) error
}
