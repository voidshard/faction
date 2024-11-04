package controller

import (
	"github.com/voidshard/faction/pkg/structs"
)

type Controller interface {
	Init(m Manager) error
	Handle(ch *structs.Change) error
}
