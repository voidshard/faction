package db

import (
	"fmt"

	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"
)

// FactionDB is a helper struct that extends a Database implementation
// with helpful additional functions.
//
// These functions are useful to callers, but can be supplied regardless of
// the implementation assuming that the Database interface is met.
//
// Ie. all implementations of Database would supply these helper functions
// by doing the same thing, so we don't have to ask implementors to do it,
// rather we simply embed the Database into ourselves and build atop it.
type FactionDB struct {
	Database
}

type DemographicQuery struct {
	// Only count people who are based in one of these areas
	Areas []string

	// Restrict count to these objects
	Objects []string
}

// InTransaction is a helper function that adds automatic commit / rollback
// depending on if an error is returned.
func (f *FactionDB) InTransaction(do func(tx ReaderWriter) error) error {
	tx, err := f.Database.Transaction()
	if err != nil {
		return nil
	}

	err = do(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// FactionSummary is a helper function that returns a summary of a faction,
// including related tuples summed with their corresponding modifiers.
func (f *FactionDB) FactionSummary(in ...string) ([]*structs.FactionSummary, error) {
	tick, err := f.Database.Tick()
	if err != nil {
		return nil, err
	}

	ff := []*FactionFilter{}
	tf := []*TupleFilter{}
	mf := []*ModifierFilter{}

	for _, id := range in {
		ff = append(ff, &FactionFilter{ID: id})
		tf = append(tf, &TupleFilter{Subject: id})
		mf = append(mf, &ModifierFilter{MinTickExpires: tick, TupleFilter: TupleFilter{Subject: id}})
	}

	fdata := map[string]*structs.FactionSummary{}
	var (
		factions []*structs.Faction
		token    string
		tuples   []*structs.Tuple
	)

	for {
		factions, token, err = f.Factions(token, ff...)
		if err != nil {
			return nil, err
		}

		for _, f := range factions {
			fdata[f.ID] = &structs.FactionSummary{
				Faction:        *f,
				Professions:    map[string]int{},
				Actions:        map[structs.ActionType]int{},
				Research:       map[string]int{},
				ResearchWeight: map[string]int{},
				Trust:          map[string]int{},
			}
		}

		if token == "" {
			break
		}
	}

	for _, r := range []Relation{
		RelationFactionProfessionWeight,
		RelationFactionActionTypeWeight,
		RelationFactionTopicResearch,
		RelationFactionTopicResearchWeight,
		RelationFactionFactionTrust,
	} {
		for {
			tuples, token, err = f.Tuples(r, token, tf...)
			if err != nil {
				return nil, err
			}

			for _, t := range tuples {
				f, ok := fdata[t.Subject]
				if !ok {
					continue
				}

				switch r {
				case RelationFactionProfessionWeight:
					f.Professions[t.Object] += t.Value
				case RelationFactionActionTypeWeight:
					f.Actions[structs.ActionType(t.Object)] += t.Value
				case RelationFactionTopicResearch:
					f.Research[t.Object] += t.Value
				case RelationFactionTopicResearchWeight:
					f.ResearchWeight[t.Object] += t.Value
				case RelationFactionFactionTrust:
					f.Trust[t.Object] += t.Value
				}

			}

			if token == "" {
				break
			}
		}

		if !r.SupportsModifiers() {
			continue
		}

		for {
			tuples, token, err = f.ModifiersSum(r, token, mf...)
			if err != nil {
				return nil, err
			}

			for _, t := range tuples {
				f, ok := fdata[t.Subject]
				if !ok {
					continue
				}
				switch r {
				case RelationFactionProfessionWeight:
					f.Professions[t.Object] += t.Value
				case RelationFactionActionTypeWeight:
					f.Actions[structs.ActionType(t.Object)] += t.Value
				case RelationFactionTopicResearch:
					f.Research[t.Object] += t.Value
				case RelationFactionFactionTrust:
					f.Trust[t.Object] += t.Value
				}
			}

			if token == "" {
				break
			}
		}
	}

	final := make([]*structs.FactionSummary, len(fdata))
	i := 0
	for _, f := range fdata {
		final[i] = f
		i++
	}

	return final, nil
}

// Demographics returns a map of Relation -> DemographicStats for the given filter(s) (areas & objects).
//
// Since it's totally impractical / not too helpful to return this for all relations
// (eg. inter personal trust), we only return a few relations, namely:
// - RelationPersonReligionFaith
// - RelationPersonProfessionSkill
// - RelationPersonFactionAffiliation
//
// For each of these we return a count of the number of people with scores within some bounds
// (see DemographicStats) for a given Object.
//
// We return a map of Relation -> Object -> DemographicStats
func (f *FactionDB) Demographics(tables []Relation, in *DemographicQuery) (map[Relation]map[string]*structs.DemographicStats, error) {
	pf := []*PersonFilter{}
	if in != nil {
		if in.Areas != nil {
			pf = make([]*PersonFilter, len(in.Areas))
			for i, area := range in.Areas {
				pf[i] = &PersonFilter{AreaID: area}
			}
		}
	}

	demoRelations := []Relation{}
	for _, r := range tables {
		// whitelist allowed relations
		if r == RelationPersonReligionFaith || r == RelationPersonProfessionSkill || r == RelationPersonFactionAffiliation {
			demoRelations = append(demoRelations, r)
		}
	}

	ret := map[Relation]map[string]*structs.DemographicStats{}
	for _, r := range demoRelations {
		ret[r] = map[string]*structs.DemographicStats{}
	}

	initialPToken := dbutils.NewToken()
	initialPToken.Limit = 500
	ptoken := initialPToken.String()

	var (
		people []*structs.Person
		tuples []*structs.Tuple
		ttoken string
		err    error
	)

	for {
		people, ptoken, err = f.People(ptoken, pf...)
		if err != nil {
			return ret, err
		}

		tf := []*TupleFilter{}
		for _, p := range people {
			if in != nil && in.Objects != nil {
				for _, obj := range in.Objects {
					tf = append(tf, &TupleFilter{Subject: p.ID, Object: obj})
				}
			} else {
				tf = append(tf, &TupleFilter{Subject: p.ID})
			}
		}

		for _, r := range demoRelations {
			stats := ret[r] // we made sure to set these

			for {
				tuples, ttoken, err = f.Tuples(r, ttoken, tf...)
				if err != nil {
					return ret, err
				}

				for _, t := range tuples {
					objStats, ok := stats[t.Object]
					if !ok {
						objStats = &structs.DemographicStats{}
					}

					objStats.Add(t.Value)

					stats[t.Object] = objStats
				}

				if ttoken == "" {
					break
				}
			}

			ret[r] = stats
		}

		if ptoken == "" {
			break
		}
	}

	return ret, nil
}

// FactionAreas returns a map of FactionID -> AreaID -> true
// for all areas that the given faction either has a Plot (a building of some kind) or a
// LandRight (a claim to work the land).
//
// That is, the faction A with Plot in Area B and LandRight in Area C would return:
// map[A]map[B:true, C:true]
func (f *FactionDB) FactionAreas(factionIDs ...string) (map[string]map[string]bool, error) {
	pfilters := make([]*PlotFilter, len(factionIDs))
	lfilters := make([]*LandRightFilter, len(factionIDs))

	for i, id := range factionIDs {
		pfilters[i] = &PlotFilter{FactionID: id}
		lfilters[i] = &LandRightFilter{FactionID: id}
	}

	var (
		land  []*structs.LandRight
		plots []*structs.Plot
		token string
		err   error
	)

	result := map[string]map[string]bool{}
	for {
		land, token, err = f.LandRights(token, lfilters...)
		if err != nil {
			return nil, err
		}
		for _, l := range land {
			farea, ok := result[l.FactionID]
			if !ok {
				farea = map[string]bool{}
			}
			farea[l.AreaID] = true
			result[l.FactionID] = farea
		}
		if token == "" {
			break
		}
	}

	for {
		plots, token, err = f.Plots(token, pfilters...)
		if err != nil {
			return nil, err
		}
		for _, p := range plots {
			farea, ok := result[p.FactionID]
			if !ok {
				farea = map[string]bool{}
			}
			farea[p.AreaID] = true
			result[p.FactionID] = farea
		}
		if token == "" {
			break
		}
	}

	return result, nil
}

// SetAreaGovernment is a helper function that changes the government id of some area(s)
func (f *FactionDB) SetAreaGovernment(govID string, areaIDs []string) error {
	if !structs.IsValidID(govID) {
		return fmt.Errorf("invalid goverment id: %s", govID)
	}

	af := []*AreaFilter{}
	for _, a := range areaIDs {
		if !structs.IsValidID(a) {
			return fmt.Errorf("invalid area id: %s", a)
		}
		af = append(af, &AreaFilter{ID: a})
	}

	var (
		areas []*structs.Area
		token string
		err   error
	)

	for {
		areas, token, err = f.Areas(token, af...)
		if err != nil {
			return err
		}

		for _, a := range areas {
			a.GovernmentID = govID
		}

		err = f.InTransaction(func(tx ReaderWriter) error {
			return tx.SetAreas(areas...)
		})
		if err != nil {
			return err
		}

		if token == "" {
			break
		}
	}

	return nil
}
