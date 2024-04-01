package simutil

import (
	mapset "github.com/deckarep/golang-set/v2"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

func SourceMeetsActionConditions(f *FactionContext, conditions [][]config.Condition) bool {
	if conditions == nil {
		return false
	}
	for _, check := range conditions {
		allow := true
		for _, condition := range check {
			switch condition {
			case config.ConditionAlways:
				// pass
			case config.ConditionSrcFactionIsCovert:
				allow = allow && f.Summary.Faction.IsCovert
			case config.ConditionSrcFactionIsNotCovert:
				allow = allow && !f.Summary.Faction.IsCovert
			case config.ConditionSrcFactionIsGovernment:
				allow = allow && f.Summary.Faction.IsGovernment
			case config.ConditionSrcFactionIsNotGovernment:
				allow = allow && !f.Summary.Faction.IsGovernment
			case config.ConditionSrcFactionIsReligion:
				allow = allow && f.Summary.Faction.IsReligion
			case config.ConditionSrcFactionHasReligion:
				allow = allow && f.Summary.Faction.ReligionID != ""
			case config.ConditionSrcFactionStructurePyramid:
				allow = allow && f.Summary.Faction.Structure == structs.FactionStructure_Pyramid
			case config.ConditionSrcFactionStructureLoose:
				allow = allow && f.Summary.Faction.Structure == structs.FactionStructure_Loose
			case config.ConditionSrcFactionStructureCell:
				allow = allow && f.Summary.Faction.Structure == structs.FactionStructure_Cell
			case config.ConditionSrcFactionHasHarvestablePlot:
				allow = allow && len(f.Land.Commodities) > 0
			}
			if !allow {
				// break out of inner loop early
				break
			}
		}
		if allow {
			return true
		}
	}
	return false
}

// ServicesForHire returns a map of services that can be hired, and the faction(s) that offer the
// service within the given areas.
func ServicesForHire(actions map[string]*config.Action, dbconn *db.FactionDB, areas []string) (map[string][]string, error) {
	// all the actions that can be paid for
	actionsForHire := []string{}
	for _, cfg := range actions {
		if cfg.MercenaryActions != nil {
			actionsForHire = append(actionsForHire, cfg.MercenaryActions...)
		}
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

	inprogress := map[string]mapset.Set[string]{}
	for {
		tuples, token, err = dbconn.Tuples(db.RelationFactionActionTypeWeight, token, q)
		if err != nil {
			return nil, err
		}

		for _, tuple := range tuples {
			if tuple.Value <= 0 {
				continue
			}
			who, ok := inprogress[tuple.Object]
			if !ok {
				who = mapset.NewSet[string]()
			}
			who.Add(tuple.Subject)
			inprogress[tuple.Object] = who
		}

		if token == "" {
			break
		}
	}

	// annnd finally shove our answer into a nice map
	done := map[string][]string{}
	for act, who := range inprogress {
		done[act] = who.ToSlice()
	}

	return done, nil
}
