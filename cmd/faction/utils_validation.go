package main

import (
	"fmt"
	"strings"

	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/util/uuid"
)

var (
	objects = map[string]structs.Object{ // map of objects we can interact with
		"world":   &structs.World{},
		"actor":   &structs.Actor{},
		"faction": &structs.Faction{},
		"culture": &structs.Culture{},
		"race":    &structs.Race{},
		"job":     &structs.Job{},
	}
	help = map[string]string{ // help text for each object
		"world":   "Worlds are the top level object in the system, every other object exists within a world.",
		"actor":   "Actors are individual entities that can interact & form factions.",
		"faction": "Factions are groups of actors that can work together, form common ground(s) or simply work against other factions.",
		"culture": "Cultures define classes, hierarchies, accepted belief, practices & skills.",
		"race":    "Race defines the physical characteristics of an actor.",
		"job":     "Jobs are some call to action put in place by a faction(s) for actor(s) to complete.",
	}
	shortNames = map[string]string{ // allows short hand names because we're lazy
		"wo": "world",
		"ac": "actor",
		"fa": "faction",
		"jo": "job",
	}
)

func validObject(name string) structs.Object {
	name = strings.ToLower(name)

	longname, ok := shortNames[name]
	if ok {
		name = longname
	}

	obj, _ := objects[name]
	return obj
}

func invalidObjectError(in string) error {
	valid := []string{}
	for k := range objects {
		valid = append(valid, k)
	}
	return fmt.Errorf("Invalid object '%s'. Valid names: %s", in, strings.Join(valid, ", "))
}

// toWorldId takes a list of strings and converts them to UUIDs if they're not already.
func toWorldId(in ...string) []string {
	out := []string{}
	for _, w := range in {
		if uuid.IsValidUUID(w) {
			log.Debug().Str("id", w).Msg("using world UUID")
			out = append(out, w)
		} else {
			id := uuid.NewID(w).String()
			log.Debug().Str("world", w).Str("id", id).Msg("using world name to determine UUID")
			out = append(out, id)
		}
	}
	return out
}
