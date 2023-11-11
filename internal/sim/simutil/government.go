package simutil

import (
	"github.com/voidshard/faction/pkg/structs"
)

func IsIllegalAction(action string, govs ...*structs.Government) bool {
	for _, g := range govs {
		_, ok := g.Outlawed.Actions[action]
		if ok {
			return true
		}
	}
	return false
}
