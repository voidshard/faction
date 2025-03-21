package controller

import (
	"context"

	"github.com/voidshard/faction/pkg/structs/v1"
)

type Controller interface {
	Init(m Manager) error
	Handle(ctx context.Context, ch *v1.Change) (bool, error)
}
