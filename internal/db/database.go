package db

import (
	"fmt"

	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	// FactionSummaryRelations are the default relations we consider when
	// building a faction summary.
	FactionSummaryRelations = []Relation{
		RelationFactionProfessionWeight,
		RelationFactionActionTypeWeight,
		RelationFactionTopicResearch,
		RelationFactionTopicResearchWeight,
		RelationFactionFactionTrust,
		RelationPersonFactionRank,
	}
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

	// Restrict Rank / Affiliation to these factions
	Factions []string

	// Restict faith to these religions
	Religions []string

	// Restrict skills to these professions
	Professions []string
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
//
// Basically this is an amalgamation of the Faction, Tuple(s) and Modifier(s) tables.
//
// Warning: This requires quite a few queries, so it's probably wise to limit exactly
// when & where this is called.
// It's useful for getting a high level view of a faction & all related info, but it's
// excessive for when you only need a small snippet of info.
func (f *FactionDB) FactionSummary(rels []Relation, in ...string) ([]*structs.FactionSummary, error) {
	if rels == nil {
		rels = FactionSummaryRelations
	}

	tick, err := f.Database.Tick()
	if err != nil {
		return nil, err
	}

	ff := []*Filter{F(ID, In, in)}
	tfSub := []*Filter{F(Subject, In, in)}
	tfObj := []*Filter{F(Object, In, in)}
	mf := []*Filter{F(TickExpires, Greater, tick), F(Subject, In, in)}

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
			fdata[f.ID] = structs.NewFactionSummary(f)
		}

		if token == "" {
			break
		}
	}

	for _, r := range rels {
		tf := tfSub
		if r == RelationPersonFactionRank {
			tf = tfObj
		}

		for {
			tuples, token, err = f.Tuples(r, token, tf...)
			if err != nil {
				return nil, err
			}

			for _, t := range tuples {
				var (
					f  *structs.FactionSummary
					ok bool
				)

				if r == RelationPersonFactionRank {
					f, ok = fdata[t.Object]
				} else {
					f, ok = fdata[t.Subject]
				}
				if !ok {
					continue
				}

				switch r {
				case RelationFactionProfessionWeight:
					f.Professions[t.Object] += t.Value
				case RelationFactionActionTypeWeight:
					f.Actions[structs.ActionType(t.Object)] += t.Value
				case RelationFactionTopicResearch:
					f.ResearchProgress[t.Object] += t.Value
				case RelationFactionTopicResearchWeight:
					f.Research[t.Object] += t.Value
				case RelationFactionFactionTrust:
					f.Trust[t.Object] += t.Value
				case RelationPersonFactionRank:
					f.Ranks.Add(structs.FactionRank(t.Value), 1)
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

// Demographics return information for the given filter(s) (areas & objects).
//
// Since it's totally impractical / not too helpful to return this for all relations
// (eg. inter personal trust), we only return a few relations, namely:
// - RelationPersonReligionFaith
// - RelationPersonProfessionSkill
// - RelationPersonFactionAffiliation
// - RelationPersonFactionRank
//
// For each of these we return a count of the number of people with scores within some bounds
// (see DemographicStats) for a given Object.
func (f *FactionDB) Demographics(in *DemographicQuery) (*structs.Demographics, error) {
	var (
		onlyFactions    map[string]bool
		onlyReligions   map[string]bool
		onlyProfessions map[string]bool
	)
	if in != nil {
		if in.Factions != nil {
			for _, id := range in.Factions {
				onlyFactions[id] = true
			}
		}
		if in.Religions != nil {
			for _, id := range in.Religions {
				onlyReligions[id] = true
			}
		}
		if in.Professions != nil {
			for _, id := range in.Professions {
				onlyProfessions[id] = true
			}
		}
	}
	pf := []*Filter{}
	if in.Areas != nil {
		pf = append(pf, F(AreaID, In, in.Areas))
	}

	permit := func(item string, set map[string]bool) bool {
		if set == nil {
			return true
		}
		v, _ := set[item]
		return v
	}

	demoRelations := []Relation{
		RelationPersonReligionFaith,
		RelationPersonProfessionSkill,
		RelationPersonFactionAffiliation,
		RelationPersonFactionRank,
	}

	ret := structs.NewDemographics()

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

		pids := []string{}
		for _, p := range people {
			pids = append(pids, p.ID)
		}
		tf := []*Filter{F(Subject, In, pids)}

		// TODO: we probably should expand the filters so this is not needed
		for _, r := range demoRelations {
			for {
				tuples, ttoken, err = f.Tuples(r, ttoken, tf...)
				if err != nil {
					return ret, err
				}

				for _, t := range tuples {
					switch r {
					case RelationPersonReligionFaith:
						if permit(t.Object, onlyReligions) {
							ret.AddFaith(t.Object, t.Value)
						}
					case RelationPersonProfessionSkill:
						if permit(t.Object, onlyProfessions) {
							ret.AddProfession(t.Object, t.Value)
						}
					case RelationPersonFactionAffiliation:
						if permit(t.Object, onlyFactions) {
							ret.AddAffiliation(t.Object, t.Value)
						}
					case RelationPersonFactionRank:
						if permit(t.Object, onlyFactions) {
							ret.AddRank(t.Object, structs.FactionRank(t.Value))
						}
					}
				}

				if ttoken == "" {
					break
				}
			}
		}

		if ptoken == "" {
			break
		}
	}

	return ret, nil
}

// AreaFactions returns a map of AreaID -> []*Faction
//
// That is, given a set of areas, which factions have influence there.
// (Inverse of FactionAreas)
func (f *FactionDB) AreaFactions(areaIDs ...string) (map[string]map[string]bool, error) {
	pfilters := []*Filter{F(AreaID, In, areaIDs)}

	var (
		plots []*structs.Plot
		token string
		err   error
	)

	result := map[string]map[string]bool{}

	for {
		plots, token, err = f.Plots(token, pfilters...)
		if err != nil {
			return nil, err
		}
		for _, p := range plots {
			areaf, ok := result[p.AreaID]
			if !ok {
				areaf = map[string]bool{}
			}
			areaf[p.FactionID] = true
			result[p.AreaID] = areaf
		}
		if token == "" {
			break
		}
	}

	return result, nil
}

// FactionAreas returns a map of FactionID -> AreaID -> true
//
// That is, given a set of factions, this is where factions have influence.
// (Inverse of AreaFactions)
func (f *FactionDB) FactionAreas(factionIDs ...string) (map[string]map[string]bool, error) {
	pfilters := []*Filter{F(FactionID, In, factionIDs)}

	var (
		plots []*structs.Plot
		token string
		err   error
	)

	result := map[string]map[string]bool{}

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

func (f *FactionDB) AreaGovernments(in ...string) (map[string]*structs.Government, error) {
	if len(in) == 0 {
		return nil, nil
	}

	// look up areas
	af := []*Filter{F(ID, In, in)}

	var (
		govIDs     map[string]bool
		areaToGovt map[string]string
		areas      []*structs.Area
		token      string
		err        error
	)
	for {
		areas, token, err = f.Areas(token, af...)
		if err != nil {
			return nil, err
		}

		for _, a := range areas {
			// nb. areas don't have to have a government
			if dbutils.IsValidID(a.GovernmentID) {
				govIDs[a.GovernmentID] = true
				areaToGovt[a.ID] = a.GovernmentID
			}
		}

		if token == "" {
			break
		}
	}

	// collect governments
	gids := []string{}
	for id := range govIDs {
		gids = append(gids, id)
	}
	gf := []*Filter{F(ID, In, gids)}

	var (
		govs    []*structs.Government
		govById map[string]*structs.Government
	)
	for {
		govs, token, err = f.Governments(token, gf...)
		if err != nil {
			return nil, err
		}

		for _, g := range govs {
			govById[g.ID] = g
		}

		if token == "" {
			break
		}
	}

	// zip up area -> gov
	ret := map[string]*structs.Government{}
	for areaID, govtID := range areaToGovt {
		gov, ok := govById[govtID]
		if !ok {
			continue
		}
		ret[areaID] = gov
	}

	return ret, nil
}

// SetAreaGovernment is a helper function that changes the government id of some area(s)
func (f *FactionDB) SetAreaGovernment(govID string, areaIDs []string) error {
	if !structs.IsValidID(govID) {
		return fmt.Errorf("invalid goverment id: %s", govID)
	}

	af := []*Filter{F(ID, In, areaIDs)}

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
