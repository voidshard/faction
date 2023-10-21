package base

import (
	"fmt"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"

	mapset "github.com/deckarep/golang-set/v2"
)

func (s *Base) applyFactionCreated(tick int, events []*structs.Event) error {
	// all the factions that were created
	subjects := eventSubjects(events)

	// find all the areas we live in
	subjectAreas, err := s.dbconn.FactionAreas(false, subjects...)
	if err != nil {
		return err
	}
	areasNeeded := mapset.NewSet[string]()
	for _, fareas := range subjectAreas { // map[factionID]map[areaID]nil
		for areaID := range fareas {
			areasNeeded.Add(areaID)
		}
	}

	// find all of the factions that are in those areas
	factionsByArea, err := s.dbconn.AreaFactions(areasNeeded.ToSlice()...) // map[areaID]map[factionID]bool

	// small optimization so we don't redo rows
	doneRelation := map[string]bool{}
	relationKey := func(a, b string) string {
		if a < b {
			return fmt.Sprintf("%s:%s", a, b)
		}
		return fmt.Sprintf("%s:%s", b, a)
	}

	// now, run through all factions in areas in which we are, we'll need to fetch their summaries ..
	for factionID, myareas := range subjectAreas {
		factionsInMyAreas := mapset.NewSet[string]()
		factionsInMyAreas.Add(factionID) // add myself

		for areaID := range myareas {
			otherIDs, ok := factionsByArea[areaID]
			if !ok {
				continue
			}
			for otherID := range otherIDs {
				_, ok := doneRelation[relationKey(factionID, otherID)]
				if ok {
					continue // we've done this relation bidirectionally already
				}
				factionsInMyAreas.Add(otherID)
			}
		}
		if factionsInMyAreas.Cardinality() == 0 {
			continue // we've done already
		}

		summaries, err := s.dbconn.FactionSummary(
			[]db.Relation{
				db.RelationFactionActionTypeWeight,
				db.RelationFactionTopicResearchWeight,
			},
			factionsInMyAreas.ToSlice()...,
		)
		if err != nil {
			return err
		}

		// find myself (so zen, much wow)
		var me *structs.FactionSummary
		for _, s := range summaries {
			if s.ID == factionID {
				me = s
				break
			}
		}
		if me == nil {
			continue
		}

		// now, for each faction in my areas, we'll need to create a relationship
		trust := []*structs.Tuple{}
		for _, other := range summaries {
			if other.ID == me.ID {
				continue
			}
			trust = append(trust,
				&structs.Tuple{
					Subject: me.ID,
					Object:  other.ID,
					Value:   determineBaseTrust(me, other),
				},
				&structs.Tuple{
					Subject: other.ID,
					Object:  me.ID,
					Value:   determineBaseTrust(other, me),
				},
			)
			doneRelation[relationKey(me.ID, other.ID)] = true
		}

		// write out our trust relationships

		err = s.dbconn.SetTuples(db.RelationFactionFactionTrust, trust...)
		if err != nil {
			return err
		}
	}

	return nil
}

func determineBaseTrust(f1, f2 *structs.FactionSummary) int {
	// In general we give larger negatives than positives. Trust is much harder to build than suspicion.

	// base value starts from Ethos distance
	// y = -0.4x + 0.2 where x is the distance between the two factions' ethos (0-1)
	// Yields
	// ->  1/5 where x = 0
	// -> -1/5 where x = 1
	// ->  0   where x = 0.5
	dist := structs.EthosDistance(&f1.Faction.Ethos, &f2.Faction.Ethos)
	base := int(float64(structs.MaxEthos)*(-0.4*dist+0.2)) - structs.MaxEthos/20

	// apply some basic "are we alike" modifiers
	if f1.Faction.IsGovernment && f2.Faction.IsCovert {
		base -= structs.MaxEthos / 5
	}
	if f1.Faction.IsCovert && f2.Faction.IsCovert {
		base += structs.MaxEthos / 20
	}
	if f1.Faction.ReligionID != "" {
		if f1.Faction.ReligionID == f2.Faction.ReligionID {
			base += structs.MaxEthos / 10
		} else {
			if f1.Faction.IsReligion {
				base -= structs.MaxEthos / 5
			} else {
				base -= structs.MaxEthos / 10
			}
		}
	}

	// more specific modifiers for favoured actions / resources
	research := map[string]bool{}
	for topic := range f1.Research { // my research topics
		research[topic] = true
	}
	sharedResearch := false
	for topic := range f2.Research { // their research topics
		_, ok := research[topic]
		if ok {
			base -= structs.MaxEthos / 10 // research rivalry
			sharedResearch = true
		}
	}
	if !sharedResearch {
		base += structs.MaxEthos / 10 // research cooperation
	}

	// favoured actions
	actions := map[structs.ActionType]bool{}
	for act := range f2.Actions { // their actions
		actions[act] = true
		switch act {
		case structs.ActionTypeFestival, structs.ActionTypeCharity:
			base += structs.MaxEthos / 20 // they're nice / charitable
		case structs.ActionTypeWar, structs.ActionTypeCrusade, structs.ActionTypeShadowWar:
			base -= structs.MaxEthos / 5 // they're warlike
		case structs.ActionTypeRaid, structs.ActionTypePillage, structs.ActionTypeEnslave:
			base -= structs.MaxEthos / 10 // they're raiders
		}
	}
	for act := range f1.Actions { // my actions
		switch act {
		case structs.ActionTypeTrade, structs.ActionTypeCraft, structs.ActionTypeHarvest:
			_, ok := actions[structs.ActionTypeSteal]
			if ok {
				base -= structs.MaxEthos / 10 // dirty theif!
			}
		case structs.ActionTypeSteal, structs.ActionTypeRaid, structs.ActionTypePillage, structs.ActionTypeEnslave:
			_, trade := actions[structs.ActionTypeTrade]
			_, craft := actions[structs.ActionTypeCraft]
			_, harvest := actions[structs.ActionTypeHarvest]
			if trade || craft || harvest {
				base -= structs.MaxEthos / 5 // looks like a tempting target ;)
			}
		case structs.ActionTypeConcealSecrets:
			_, investigate := actions[structs.ActionTypeGatherSecrets]
			_, spy := actions[structs.ActionTypeHireSpies]
			if investigate || spy {
				base -= structs.MaxEthos / 5 // we don't like people who snoop on us
			}
		case structs.ActionTypeConsolidate, structs.ActionTypePropoganda:
			_, frame := actions[structs.ActionTypeFrame]
			_, blackmail := actions[structs.ActionTypeBlackmail]
			_, propoganda := actions[structs.ActionTypePropoganda]
			if frame || blackmail || propoganda {
				base -= structs.MaxEthos / 5 // we don't like people who talk about us
			}
		}
	}

	return base
}
