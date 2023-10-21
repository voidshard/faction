package db

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/internal/log"
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
		return err
	}

	err = do(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// TuplesSumModsBySubject is a helper function that returns a map of object -> int
// for a given relation.
//
// If a relation doesn't support modifiers, this is simply tuples with the given subject & object(s).
// If it does support modifiers we add on to these tuples the sum of modifiers for the given subject & object(s) too.
func (f *FactionDB) TuplesSumModsBySubject(r Relation, subject string, objects ...string) (map[string]int, error) {
	q := Q(F(Object, In, objects), F(Subject, Equal, subject))
	if len(objects) == 0 {
		q = Q(F(Subject, Equal, subject))
	}

	var (
		tuples []*structs.Tuple
		token  string
		err    error
		result = map[string]int{}
	)

	for {
		tuples, token, err = f.Tuples(r, token, q)
		log.Debug().Err(err).Msg()("fetching tuples from database")
		if err != nil {
			return nil, err
		}

		for _, t := range tuples {
			result[t.Object] = t.Value
		}

		if token == "" {
			break
		}
	}

	if !r.SupportsModifiers() {
		return result, nil
	}

	for {
		tuples, token, err = f.ModifiersSum(r, token, q)
		log.Debug().Err(err).Msg()("fetching modifiers from database")
		if err != nil {
			return nil, err
		}

		for _, t := range tuples {
			v, _ := result[t.Object]
			result[t.Object] = v + t.Value
		}

		if token == "" {
			break
		}
	}

	return result, nil
}

func (f *FactionDB) LandSummary(areas, factions []string) (*structs.LandSummary, error) {
	filters := []*Filter{}
	if areas != nil && len(areas) > 0 {
		filters = append(filters, F(AreaID, In, areas))
	}
	if factions != nil && len(factions) > 0 {
		filters = append(filters, F(FactionID, In, factions))
	}
	if len(filters) == 0 {
		return nil, fmt.Errorf("no filters supplied")
	}

	q := Q(filters...)
	sum := structs.NewLandSummary()

	var (
		plots []*structs.Plot
		token string
		err   error
	)
	for {
		plots, token, err = f.Plots(token, q)
		log.Debug().Err(err).Msg()("fetching plots from database")
		if err != nil {
			return nil, err
		}

		for _, p := range plots {
			sum.Add(p)
		}

		if token == "" {
			break
		}
	}

	return sum, nil
}

// FactionSummary is a helper function that returns a summary of a faction,
// including related tuples summed with their corresponding modifiers.
//
// Basically this is an amalgamation of the Faction, Tuple(s) and Modifier(s) tables.
//
// Warning: This requires quite a few queries, so it's probably wise to limit exactly
// when & where this is called.
//
// It's useful for getting a high level view of a faction & all related info, but it's
// excessive for when you only need a small snippet of info.
//
// In general we call this when deciding on what a given faction wants to do when planning an action;
// because we need to check pretty much all of it's values & settings to decide wisely.
func (f *FactionDB) FactionSummary(rels []Relation, in ...string) ([]*structs.FactionSummary, error) {
	if rels == nil {
		rels = FactionSummaryRelations
	}

	tick, err := f.Database.Tick()
	if err != nil {
		return nil, err
	}

	ff := Q(F(ID, In, in))
	tfSub := Q(F(Subject, In, in))
	tfObj := Q(F(Object, In, in))
	mf := Q(F(TickExpires, Greater, tick)).Or(F(Subject, In, in))

	fdata := map[string]*structs.FactionSummary{}
	var (
		factions []*structs.Faction
		token    string
		tuples   []*structs.Tuple
	)

	for {
		factions, token, err = f.Factions(token, ff)
		log.Debug().Err(err).Msg()("fetching factions from database")
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
			tuples, token, err = f.Tuples(r, token, tf)
			log.Debug().Err(err).Msg()("fetching tuples from database")
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
			tuples, token, err = f.ModifiersSum(r, token, mf)
			log.Debug().Err(err).Msg()("fetching modifier summation from database")
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

	pf := Q()
	if in.Areas != nil {
		pf.Or(F(AreaID, In, in.Areas))
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
		people, ptoken, err = f.People(ptoken, pf)
		log.Debug().Err(err).Msg()("fetching people from database")
		if err != nil {
			return ret, err
		}

		pids := []string{}
		for _, p := range people {
			pids = append(pids, p.ID)
		}
		tf := Q(F(Subject, In, pids))

		// TODO: we probably should expand the filters so this is not needed
		for _, r := range demoRelations {
			for {
				tuples, ttoken, err = f.Tuples(r, ttoken, tf)
				log.Debug().Err(err).Msg()("fetching tuples from database")
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
	pfilters := Q(F(AreaID, In, areaIDs))

	var (
		plots []*structs.Plot
		token string
		err   error
	)

	result := map[string]map[string]bool{}

	for {
		plots, token, err = f.Plots(token, pfilters)
		log.Debug().Err(err).Msg()("fetching plots from database")
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

// FactionAreas returns a map of FactionID -> AreaID -> nil || Area
//
// That is, given a set of factions, this is where factions have influence.
// (Inverse of AreaFactions)
//
// if `includeAreas` we do the additional query to return the Area structs, otherwise
// the map will be populated with nil values (that indicate the faction has influence
// in an area, but the area data isn't fetched).
func (f *FactionDB) FactionAreas(includeAreas bool, factionIDs ...string) (map[string]map[string]*structs.Area, error) {
	pfilters := Q(F(FactionID, In, factionIDs))

	var (
		plots []*structs.Plot
		token string
		err   error
	)

	result := map[string]map[string]*structs.Area{}
	aset := mapset.NewSet[string]()

	for {
		plots, token, err = f.Plots(token, pfilters)
		log.Debug().Err(err).Msg()("fetching plots from database")
		if err != nil {
			return nil, err
		}
		for _, p := range plots {
			farea, ok := result[p.FactionID]
			if !ok {
				farea = map[string]*structs.Area{}
			}
			farea[p.AreaID] = nil
			result[p.FactionID] = farea
			aset.Add(p.AreaID)
		}
		if token == "" {
			break
		}
	}

	if !includeAreas {
		return result, nil
	}

	areas, _, err := f.Areas("", Q(F(ID, In, aset.ToSlice())))
	if err != nil {
		return result, err
	}
	amap := map[string]*structs.Area{}
	for _, a := range areas {
		amap[a.ID] = a
	}

	for _, fareas := range result {
		for areaID := range fareas {
			fareas[areaID] = amap[areaID]
		}
	}

	return result, nil
}

func (f *FactionDB) AreaGovernments(in ...string) (map[string]*structs.Government, error) {
	if len(in) == 0 {
		return nil, nil
	}

	// look up areas
	af := Q(F(ID, In, in))

	govIDs := map[string]bool{}
	areaToGovt := map[string]string{}
	var (
		areas []*structs.Area
		token string
		err   error
	)
	for {
		areas, token, err = f.Areas(token, af)
		log.Debug().Err(err).Msg()("fetching areas from database")
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

	if len(govIDs) == 0 {
		return map[string]*structs.Government{}, nil
	}

	// collect governments
	gids := []string{}
	for id := range govIDs {
		gids = append(gids, id)
	}
	gf := Q(F(ID, In, gids))

	var govs []*structs.Government
	govById := map[string]*structs.Government{}
	for {
		govs, token, err = f.Governments(token, gf)
		log.Debug().Err(err).Msg()("fetching governments from database")
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

	af := Q(F(ID, In, areaIDs))

	var (
		areas []*structs.Area
		token string
		err   error
	)

	for {
		areas, token, err = f.Areas(token, af)
		log.Debug().Err(err).Msg()("fetching areas from database")
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
