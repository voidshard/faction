package simutil

import (
	"github.com/voidshard/faction/internal/random/rng"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// ActionWeights is a helper for weighting actions, since we apply a lot of weights and calculations to them.
type ActionWeights struct {
	prob map[structs.ActionType]float64
	defn map[structs.ActionType]*config.Action
	goal map[structs.Goal][]structs.ActionType

	normal  rng.Normalised
	choices []structs.ActionType
}

// NewActionWeights creates a new ActionWeights helper.
func NewActionWeights(actions map[structs.ActionType]*config.Action) *ActionWeights {
	prob := map[structs.ActionType]float64{}
	goal := map[structs.Goal][]structs.ActionType{}
	for atype, act := range actions {
		// build a map of action types to probabilities
		p := act.Probability
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

// Choose returns a random action type based on the current weights.
func (w *ActionWeights) Choose() structs.ActionType {
	return w.choices[w.normal.Int()]
}

// normalise builds a normalised list of probabilities and a list of action types.
// Returns the number of choices (> 0 probability), this should be checked
func (w *ActionWeights) Normalise() int {
	choices := []structs.ActionType{}
	probs := []float64{}
	for t, p := range w.prob {
		if p <= 0.0 {
			continue
		}
		choices = append(choices, t)
		probs = append(probs, p)
	}
	if len(choices) > 0 {
		w.normal = rng.NewNormalised(probs)
		w.choices = choices
	}
	return len(choices)
}

// WeightAction multiplies the probability of the given action by the given multiplier.
func (w *ActionWeights) WeightAction(mult float64, act structs.ActionType) {
	p, _ := w.prob[act]
	w.prob[act] = mult * p
	w.normal = nil
}

// WeightByMinPeople multiplies the probability of actions whose min people is more than `minPeople` by the given `mult`
func (w *ActionWeights) WeightByMinPeople(mult float64, i int) {
	w.normal = nil
	for atype, act := range w.defn {
		if act.MinPeople > i {
			w.prob[atype] *= mult
		}
	}
}

// WeightByCost multiplies the probability of actions whose min cost is more than `maxPrice` by the given `mult`
// Ie. we use this to prevent us from chosing actions out of our budget
func (w *ActionWeights) WeightByCost(mult, maxPrice float64) {
	w.normal = nil
	for atype, act := range w.defn {
		if act.Cost.Min > maxPrice {
			w.prob[atype] *= mult
		}
	}
}

// ApplyActionConditions applies conditions (see config/conditions.go) to the given faction.
//
// Ie. we rule out some actions based on the faction's settings.
//
// We use this to restrict say Crusades to religious factions, or only allow the government to
// do some actions.
func (w *ActionWeights) ApplyActionConditions(f *structs.FactionSummary) {
	for a, act := range w.defn {
		if act.Restricted == nil || len(act.Restricted) == 0 {
			continue // no restrictions
		}
		if !SourceMeetsActionConditions(f, act.Restricted) {
			w.prob[a] = 0.0
		}
	}
}

// WeightByGoal multiplies the probability of all actions that match the given goal(s) by the given multiplier.
func (w *ActionWeights) WeightByGoal(mult float64, goals ...structs.Goal) {
	w.normal = nil
	for _, g := range goals {
		for _, a := range w.goal[g] {
			w.prob[a] *= mult
		}
	}
}

// WeightByEthos makes actions less likely the further they are away from the given ethos.
func (w *ActionWeights) WeightByEthos(e *structs.Ethos) {
	w.normal = nil
	for atype, act := range w.defn {
		w.prob[atype] *= (1 - structs.EthosDistance(e, &act.Ethos))
	}
}

// WeightByTypes applies a given weight (MinTuple -> MaxTuple) to the probability of all actions of the given type.
// Nb. this is how we apply Faction Focus weights (see Focus struct, pkg/config/faction.go).
func (w *ActionWeights) WeightByTypes(weights map[structs.ActionType]int) {
	w.normal = nil
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
	w.normal = nil
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
