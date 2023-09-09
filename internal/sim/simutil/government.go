package simutil

import (
	"github.com/voidshard/faction/pkg/structs"
)

func IsIllegalAction(act structs.ActionType, govs ...*structs.Government) bool {
	for _, g := range govs {
		_, ok := g.Outlawed.Actions[act]
		if ok {
			return true
		}
	}
	return false
}
