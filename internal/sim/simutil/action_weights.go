package simutil

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	// restrictedActions must be explicitly enabled (see SetIsReligion() and SetIsGovernment())
	restrictedActions = map[structs.ActionType]bool{}
)

// ActionWeights is a helper for weighting actions, since we apply a lot of weights and calculations to them.
type ActionWeights struct {
	prob map[structs.ActionType]float64
	defn map[structs.ActionType]*config.Action
	goal map[structs.Goal][]structs.ActionType
}

// NewActionWeights creates a new ActionWeights helper.
func NewActionWeights(actions map[structs.ActionType]*config.Action) *ActionWeights {
	prob := map[structs.ActionType]float64{}
	goal := map[structs.Goal][]structs.ActionType{}
	for atype, act := range actions {
		// build a map of action types to probabilities
		p := act.Probability
		_, ok := restrictedActions[atype]
		if ok {
			p = 0.0 // restricted actions are 0 unless SetIsReligion() or SetIsGovernment() is called
		}
		prob[atype] = p

		// build a map of goals to actions
		for _, g := range act.Goals {
			cur, ok := goal[g]
			if !ok {
				cur = []structs.ActionType{}
			}
			goal[g] = append(cur, atype)
		}
	}
	return &ActionWeights{
		prob: prob,
		defn: actions,
		goal: goal,
	}
}

// SetIsReligion sets the probability of all religion-only actions to their starting values.
// If not set, the probability of these actions is 0.0
func (w *ActionWeights) SetIsReligion() {
	for _, a := range structs.ReligionOnlyActions {
		act, ok := w.defn[a]
		if ok {
			w.prob[a] = act.Probability
		}
	}
}

// SetIsGovernment sets the probability of all government-only actions to their starting values.
// If not set, the probability of these actions is 0.0
func (w *ActionWeights) SetIsGovernment() {
	for _, a := range structs.GovernmentOnlyActions {
		act, ok := w.defn[a]
		if ok {
			w.prob[a] = act.Probability
		}
	}
}

// WeightByGoal multiplies the probability of all actions that match the given goal(s) by the given multiplier.
func (w *ActionWeights) WeightByGoal(mult float64, goals ...structs.Goal) {
	for _, g := range goals {
		for _, a := range w.goal[g] {
			w.prob[a] *= mult
		}
	}
}

// WeightByEthos makes actions less likely the further they are away from the given ethos.
func (w *ActionWeights) WeightByEthos(e *structs.Ethos) {
	for atype, act := range w.defn {
		w.prob[atype] *= (1 - structs.EthosDistance(e, &act.Ethos))
	}
}

// WeightByTypes applies a given weight (MinTuple -> MaxTuple) to the probability of all actions of the given type.
// Nb. this is how we apply Faction Focus weights (see Focus struct, pkg/config/faction.go).
func (w *ActionWeights) WeightByTypes(weights map[structs.ActionType]int) {
	for atype, value := range weights {
		_, ok := w.prob[atype]
		if ok {
			w.prob[atype] *= 1.0 + float64(value/structs.MaxEthos)
		}
	}
}

// WeightByIllegal multiplies the probability of all actions that are illegal by the given multiplier.
// Nb. If the Action is illegal in multiple governments, the multiplier is applied only once.
func (w *ActionWeights) WeightByIllegal(mult float64, govs ...*structs.Government) {
	banned := map[structs.ActionType]bool{}
	for _, g := range govs {
		for a := range g.Outlawed.Actions {
			banned[a] = true
		}
	}
	for a := range banned {
		w.prob[a] *= mult
	}
}

func init() {
	for _, a := range structs.ReligionOnlyActions {
		restrictedActions[a] = true
	}
	for _, a := range structs.GovernmentOnlyActions {
		restrictedActions[a] = true
	}
}
