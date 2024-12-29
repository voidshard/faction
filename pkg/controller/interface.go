package controller

import (
	"context"

	"github.com/voidshard/faction/pkg/structs"
)

type Controller interface {
	Init(m Manager) error
	Handle(ctx context.Context, ch *structs.Change) (bool, error)
}
