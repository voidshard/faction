package config

import (
	"github.com/voidshard/faction/pkg/structs"
)

type Faction struct {
	// Ethos settings
	EthosMean      structs.Ethos
	EthosDeviation structs.Ethos

	// Actions that can be selected as a "focus"
	FocusActions []structs.ActionType
	// Probability of a faction having some number of focuses (by index).
	//
	// Eg. given [0.5, 0.3, 0.2], a faction has a 50% chance of having 0 focuses,
	// a 30% chance of having 1 focus, and a 20% chance of having 2 focuses.
	//
	// A 'focus' is a action(s) that the given faction prefers to perform
	// ("business as usual"). A faction with an illegal focus is covert by definition.
	//
	// Nb. you almost certainly want a faction to have at least one focus
	FocusProbability []float64
}
