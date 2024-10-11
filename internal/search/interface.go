package search

import (
	"context"

	"github.com/voidshard/faction/pkg/structs"
)

type Search interface {
	IndexActors(context.Context, string, []*structs.Actor, bool) error
	IndexFactions(context.Context, string, []*structs.Faction, bool) error

	DeleteActor(context.Context, string, string) error
	DeleteFaction(context.Context, string, string) error
}
