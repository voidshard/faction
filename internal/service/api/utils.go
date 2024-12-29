package api

import (
	"bytes"
	"fmt"
)

const (
	factionMaxRelations = 20
	factionMaxMemories  = 10 // nb. this is per relation

	actorMaxRelations = 10
	actorMaxMemories  = 5 // nb. this is per relation
)

/*
func tidyRelations(in structs.Interactive, maxRelations int) {
	// dedupe
	relations := []*structs.Relationship{}
	seen := map[string]bool{}
	for _, r := range append(in.GetAllies(), in.GetEnemies()...) {
		_, ok := seen[r.ID.Value]
		if ok {
			continue
		}
		seen[r.ID.Value] = true
		relations = append(relations, r)
	}

	// divide into allies and enemies
	sort.Sort(structs.Relationships(relations))
	i := 0
	for ; i < len(relations); i++ {
		if relations[i].Trust > 0 { // first positive trust
			// Easy case since we're sorted by low -> high trust, set enemies from the front
			in.SetEnemies(relations[:math.Min(i, maxRelations)])

			// Harder case, set allies from the back, then reverse (high -> low trust)
			// Nb. since we've found Trust > 0 we know there's at least one ally
			numAllies := math.Min(len(relations)-i-1, maxRelations)
			allies := relations[len(relations)-numAllies:]
			sort.Sort(sort.Reverse(structs.Relationships(allies)))
			in.SetAllies(allies)
			break
		}
	}
	if i >= len(relations) { // Covers the case where we have no allies
		in.SetEnemies(relations[:math.Min(i, maxRelations)])
		in.SetAllies(nil)
	}
}

func tidyRelationMemories(in []*structs.Memory, maxMemories int) []*structs.Memory {
	dedupe := map[string][]*structs.Memory{}
	for _, m := range in {
		cat, ok := dedupe[m.Category]
		if !ok {
			cat = []*structs.Memory{}
		}
		dedupe[m.Category] = append(cat, m)
	}

	all := []*structs.Memory{}
	for cat, mems := range dedupe {
		if cat == "" {
			// no uniqueness applied to empty category
			all = append(all, mems...)
		} else {
			// sort and append only the greatest
			sort.Sort(structs.Memories(mems))
			all = append(all, mems[0])
		}
	}

	sort.Sort(structs.Memories(all))
	return all[:math.Min(len(all), maxMemories)]
}

func tidyMemories(in structs.Interactive, maxMemories int) {
	for _, relation := range in.GetAllies() {
		memories := relation.GetMemories()
		relation.SetMemories(tidyRelationMemories(memories, maxMemories))
	}
	for _, relation := range in.GetEnemies() {
		memories := relation.GetMemories()
		relation.SetMemories(tidyRelationMemories(memories, maxMemories))
	}
}
*/

func defaultValue(value *uint64, otherwise uint64) uint64 {
	if value == nil {
		return otherwise
	}
	return *value
}

func clamp(value *uint64, min, max, ifNull uint64) uint64 {
	v := defaultValue(value, ifNull)
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}

func encodeRequest(in marshalable) ([]byte, error) {
	method := in.Kind()
	data, err := in.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return bytes.Join([][]byte{[]byte(method), data}, []byte("|")), nil
}

func decodeRequest(data []byte) (string, []byte, error) {
	sections := bytes.SplitN(data, []byte("|"), 2)
	if len(sections) != 2 {
		return "", nil, fmt.Errorf("invalid request")
	}
	return string(sections[0]), sections[1], nil
}
