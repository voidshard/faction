package simutil

import (
	mapset "github.com/deckarep/golang-set/v2"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"
)

// ServicesForHire returns a map of services that can be hired, and the faction(s) that offer the
// service within the given areas.
func ServicesForHire(dbconn *db.FactionDB, areas []string) (map[structs.ActionType][]string, error) {
	// all the actions that can be paid for
	actionsForHire := []string{}
	for actionType, _ := range structs.ActionsForMercenaries {
		actionsForHire = append(actionsForHire, string(actionType))
	}
	for actionType, _ := range structs.ActionsForSpies {
		actionsForHire = append(actionsForHire, string(actionType))
	}

	// all factions within the given areas
	areaToFactionMap, err := dbconn.AreaFactions(areas...)
	if err != nil {
		return nil, err
	}

	factionIDs := mapset.NewSet[string]()
	for _, factionToBool := range areaToFactionMap {
		for factionID, _ := range factionToBool {
			factionIDs.Add(factionID)
		}
	}

	// look up who prefers what action(s)
	q := db.Q(db.F(db.Subject, db.In, factionIDs), db.F(db.Object, db.In, actionsForHire))
	var (
		tuples []*structs.Tuple
		token  string
	)

	inprogress := map[structs.ActionType]mapset.Set[string]{}
	for {
		tuples, token, err = dbconn.Tuples(db.RelationFactionActionTypeWeight, token, q)
		if err != nil {
			return nil, err
		}

		for _, tuple := range tuples {
			if tuple.Value <= 0 {
				continue
			}
			act := structs.ActionType(tuple.Object)
			who, ok := inprogress[act]
			if !ok {
				who = mapset.NewSet[string]()
			}
			who.Add(tuple.Subject)
			inprogress[act] = who
		}

		if token == "" {
			break
		}
	}

	// annnd finally shove our answer into a nice map
	done := map[structs.ActionType][]string{}
	for act, who := range inprogress {
		done[act] = who.ToSlice()
	}

	return done, nil
}
