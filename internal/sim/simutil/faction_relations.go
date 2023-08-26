package simutil

import "github.com/voidshard/faction/pkg/structs"

type FactionRelations struct {
	// So aligned it's almost the same faction
	Federated []string // 9k -> 10k

	// Share many agreements, goals, usually enemies of the same factions
	Allied []string // 6k -> 9k

	// Share many goals & agreements
	Sympathetic []string // 3k -> 6k

	// Have been known to work together but lack common goals & formal agreements
	Friendly []string // 1k -> 3k

	// No strong feelings either way
	Neutral []string // 1k -> -1k

	// Have been known to work against one another
	Unfriendly []string // -1k -> -3k

	// Routinely work against one another, and have a history of conflict
	Rival []string // -3k -> -6k

	// Have a history of conflict, and are actively working against one another
	Hostile []string // -6k -> -9k

	// Sworn enemies, actively working to destroy one another
	Nemesis []string // -9k -> -10k
}

func NewFactionRelations() *FactionRelations {
	return &FactionRelations{
		Federated:   []string{},
		Allied:      []string{},
		Sympathetic: []string{},
		Friendly:    []string{},
		Neutral:     []string{},
		Unfriendly:  []string{},
		Rival:       []string{},
		Hostile:     []string{},
		Nemesis:     []string{},
	}
}

func (r *FactionRelations) Add(id string, w int) {
	if w < structs.MaxTuple/10 && w > structs.MinTuple/10 {
		r.Neutral = append(r.Neutral, id)
	} else if w > 0 { // positive
		if w > (structs.MaxTuple*9)/10 { // 90%+
			r.Federated = append(r.Federated, id)
		} else if w > (structs.MaxTuple*6)/10 { // 60%+
			r.Allied = append(r.Allied, id)
		} else if w > (structs.MaxTuple*3)/10 { // 30%+
			r.Sympathetic = append(r.Sympathetic, id)
		} else {
			r.Friendly = append(r.Friendly, id)
		}
	} else {
		if w < (structs.MinTuple*9)/10 { // 90%-
			r.Nemesis = append(r.Nemesis, id)
		} else if w < (structs.MinTuple*6)/10 { // 60%-
			r.Hostile = append(r.Hostile, id)
		} else if w < (structs.MinTuple*3)/10 { // 30%-
			r.Rival = append(r.Rival, id)
		} else {
			r.Unfriendly = append(r.Unfriendly, id)
		}
	}
}
