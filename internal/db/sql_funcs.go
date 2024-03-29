package db

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// sqlOperator is something that can perform an sql operation read/write
// We do this so we can have some lower level funcs that perform the query logic regardless
// of whether we are in a transaction or not.
//
// Basically both sqlx.Tx and sqlx.DB implement this interface so we can use them interchangeably.
// As in, we can run the same code without having to worry about whether we are in a transaction or not.
type sqlOperator interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
}

// mstruct is a row of metadata
type mstruct struct {
	ID  string `db:"id"`
	Str string `db:"str"`
	Int int    `db:"int"`
}

// lawStruct holds a law row. We don't provide these as a first class object,
// but internally laws are written as individual rows.
type lawStruct struct {
	SourceID string          `db:"source_id"`
	MetaKey  structs.MetaKey `db:"meta_key"`
	MetaVal  string          `db:"meta_val"`
	Illegal  bool            `db:"illegal"`
}

// rankedPlot holds a Plot along side a rank setting (for windowed queries)
type rankedPlot struct {
	structs.Plot
	Rank int `db:"rnk"`
}

// rankedTuple holds a Tuple along side a rank setting (for windowed queries)
type rankedTuple struct {
	structs.Tuple
	Rank int `db:"rnk"`
}

func factionLeadership(op sqlOperator, limit int, ids ...string) (map[string]*FactionLeadership, error) {
	result := map[string]*FactionLeadership{}
	if len(ids) == 0 {
		return result, nil
	}

	if limit < 1 {
		limit = 1
	}

	inClause, vals := sqlIn(ids)

	// Essentially for tuples in our list of faction IDs (the "object" here) we partition
	// sort by value and then take the top N (limit) for each faction.
	// Thus we return the highest ranked people of each faction given to us, down to some given limit.
	//
	// stackoverflow.com/questions/28119176/select-top-n-record-from-each-group-sqlite
	qstr := fmt.Sprintf(`SELECT * FROM (
	    SELECT ROW_NUMBER() OVER (
		PARTITION BY object
		ORDER BY value DESC
	    ) as rnk, *
	    FROM %s
	    WHERE object IN (%s)
	)
	WHERE rnk <= %d;`, RelationPersonFactionRank.tupleTable(), inClause, limit)

	var out []*rankedTuple
	err := op.Select(&out, qstr, vals...)
	if err != nil {
		return nil, err
	}

	// TODO: consider a join?
	peopleIDs := []string{}
	for _, t := range out {
		peopleIDs = append(peopleIDs, t.Subject)
	}

	peeps, _, err := people(op, "", Q(F(ID, In, peopleIDs), F(DeathTick, Equal, 0)))
	if err != nil {
		return nil, err
	}

	pmap := map[string]*structs.Person{}
	for _, p := range peeps {
		pmap[p.ID] = p
	}

	for _, t := range out {
		leaders, ok := result[t.Object]
		if !ok {
			leaders = NewFactionLeadership()
		}
		person, ok := pmap[t.Subject]
		if ok {
			leaders.Add(structs.FactionRank(t.Value), person)
			result[t.Object] = leaders
		}
	}

	return result, nil
}

func factionPlots(op sqlOperator, limit int, ids ...string) (map[string][]*structs.Plot, error) {
	result := map[string][]*structs.Plot{}
	if len(ids) == 0 {
		return result, nil
	}

	if limit < 1 {
		limit = 1
	}

	inClause, vals := sqlIn(ids)

	qstr := fmt.Sprintf(`SELECT * FROM (
	    SELECT ROW_NUMBER() OVER (
		PARTITION BY faction_id
		ORDER BY value DESC
	    ) as rnk, *
	    FROM %s
	    WHERE faction_id IN (%s)
	)
	WHERE rnk <= %d;`, tablePlots, inClause, limit)

	var out []*rankedPlot
	err := op.Select(&out, qstr, vals...)
	if err != nil {
		return nil, err
	}

	for _, t := range out {
		fplots, ok := result[t.FactionID]
		if !ok {
			fplots = []*structs.Plot{}
		}
		fplots = append(fplots, &(t).Plot)
		result[t.FactionID] = fplots
	}

	return result, nil
}

func deleteModifiers(op sqlOperator, r Relation, expires_before_tick int) error {
	if !r.SupportsModifiers() {
		return nil
	}

	qstr := fmt.Sprintf(`DELETE FROM %s WHERE tick_expires < :time;`, r.modTable())

	_, err := op.NamedExec(qstr, map[string]interface{}{
		"time": expires_before_tick,
	})
	return err
}

func modifiersSum(op sqlOperator, r Relation, token string, in *Query) ([]*structs.Tuple, string, error) {
	if !r.SupportsModifiers() {
		return nil, token, fmt.Errorf("relation %s does not support modifiers", r)
	}

	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := sqlSummationTuplesFromModifiers(r, tk, in)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Tuple
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func modifiers(op sqlOperator, r Relation, token string, in *Query) ([]*structs.Modifier, string, error) {
	if !r.SupportsModifiers() {
		return nil, token, fmt.Errorf("relation %s does not support modifiers", r)
	}

	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := sqlFromModifierFilters(r, tk, in)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Modifier
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setModifiers(op sqlOperator, r Relation, in []*structs.Modifier) error {
	if !r.SupportsModifiers() {
		return fmt.Errorf("relation %s does not support modifiers", r)
	}
	if len(in) == 0 {
		return nil
	}
	for _, i := range in {
		i.Value = clampInt(i.Value, structs.MinTuple, structs.MaxTuple)
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    subject, object, value, tick_expires, meta_key, meta_val, meta_reason
	) VALUES (
	    :subject, :object, :value, :tick_expires, :meta_key, :meta_val, :meta_reason
	);`, r.modTable())

	_, err := op.NamedExec(qstr, in)
	return err
}

func incrModifiers(op sqlOperator, r Relation, v int, in *Query) error {
	if in == nil || len(in.filters) == 0 {
		return nil
	}

	where, args, err := in.sqlQuery(1) // 1 is taken
	if err != nil {
		return err
	}
	args = append([]interface{}{v}, args...) // add our value to the front

	qstr := fmt.Sprintf(
		`UPDATE %s SET value = MAX(MIN(value + $1, %d), %d) %s;`,
		r.modTable(),
		structs.MaxTuple, structs.MinTuple,
		where,
	)

	_, err = op.Exec(qstr, args...)
	return err
}

func countTuples(op sqlOperator, r Relation, in *Query) (int, error) {
	in.DisableSort() // not that we use sort anyway but :shrug:
	q, args, err := genericCountSQLFromFilters(in, r.tupleTable())
	if err != nil {
		return 0, err
	}
	var count int
	return count, op.Get(&count, q, args...)
}

func tuples(op sqlOperator, r Relation, token string, in *Query) ([]*structs.Tuple, string, error) {
	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := sqlFromTupleFilters(r, tk, in)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Tuple
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setTuples(op sqlOperator, r Relation, in []*structs.Tuple) error {
	if len(in) == 0 {
		return nil
	}
	for _, i := range in {
		i.Value = clampInt(i.Value, structs.MinTuple, structs.MaxTuple)
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    subject, object, value
	) VALUES (
	    :subject, :object, :value
	) ON CONFLICT (subject, object) DO UPDATE SET value=EXCLUDED.value;`, r.tupleTable())

	_, err := op.NamedExec(qstr, in)
	return err
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func incrTuples(op sqlOperator, r Relation, v int, in *Query) error {
	if in == nil || len(in.filters) == 0 {
		return nil
	}

	where, args, err := in.sqlQuery(1) // 1 is taken
	if err != nil {
		return err
	}

	args = append([]interface{}{v}, args...) // add our value to the front
	qstr := fmt.Sprintf(
		`UPDATE %s SET value = MAX(MIN(value + $1, %d), %d) %s;`,
		r.tupleTable(),
		structs.MaxTuple, structs.MinTuple,
		where,
	)

	_, err = op.Exec(qstr, args...)
	return err
}

func plots(op sqlOperator, token string, in *Query) ([]*structs.Plot, string, error) {
	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := sqlFromPlotFilters(tk, in)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Plot
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setPlots(op sqlOperator, in []*structs.Plot) error {
	if len(in) == 0 {
		return nil
	}

	for _, f := range in {
		if !dbutils.IsValidID(f.ID) {
			return fmt.Errorf("plot id %s is invalid", f.ID)
		}
		if !dbutils.IsValidID(f.AreaID) {
			return fmt.Errorf("plot area id %s is invalid", f.AreaID)
		}
		if f.FactionID != "" && !dbutils.IsValidID(f.FactionID) {
			return fmt.Errorf("plot faction id %s is invalid", f.FactionID)
		}
		if f.Commodity == "" {
			f.Yield = 0
		}
		if f.Yield < 0 {
			f.Yield = 0
		}
		if f.Size < 1 {
			f.Size = 1
		}
	}

	_, err := op.NamedExec(sqlInsertPlots, in)
	return err
}

func people(op sqlOperator, token string, in *Query) ([]*structs.Person, string, error) {
	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := genericSQLFromFilters(tk, in, tablePeople)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Person
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setPeople(op sqlOperator, in []*structs.Person) error {
	if len(in) == 0 {
		return nil
	}

	for _, f := range in {
		if !dbutils.IsValidID(f.ID) {
			return fmt.Errorf("person id %s is invalid", f.ID)
		}
		if f.BirthFamilyID != "" && !dbutils.IsValidID(f.BirthFamilyID) {
			return fmt.Errorf("person birth family id %s is invalid", f.BirthFamilyID)
		}
		if !dbutils.IsValidID(f.AreaID) {
			return fmt.Errorf("person area id %s is invalid", f.AreaID)
		}
		if f.JobID != "" && !dbutils.IsValidID(f.JobID) {
			return fmt.Errorf("person job id %s is invalid", f.JobID)
		}
		if f.Race == "" {
			return fmt.Errorf("person id %s race required", f.ID)
		}
		if f.Culture == "" {
			return fmt.Errorf("person id %s culture required", f.ID)
		}
		f.Clamp()
		if f.Random == 0 {
			f.Random = rng.Intn(structs.RandomMax)
		}
	}

	_, err := op.NamedExec(sqlInsertPeople, in)
	return err
}

func events(op sqlOperator, token string, in *Query) ([]*structs.Event, string, error) {
	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := genericSQLFromFilters(tk, in, tableEvents)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Event
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setEvents(op sqlOperator, in []*structs.Event) error {
	if len(in) == 0 {
		return nil
	}

	for _, f := range in {
		if !dbutils.IsValidID(f.ID) {
			return fmt.Errorf("event id %s is invalid", f.ID)
		}
	}

	_, err := op.NamedExec(sqlInsertEvents, in)
	return err
}

func jobs(op sqlOperator, token string, in *Query) ([]*structs.Job, string, error) {
	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := genericSQLFromFilters(tk, in, tableJobs)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Job
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setJobs(op sqlOperator, in []*structs.Job) error {
	if len(in) == 0 {
		return nil
	}

	for _, f := range in {
		if !dbutils.IsValidID(f.ID) {
			return fmt.Errorf("job id %s is invalid", f.ID)
		}
		if !dbutils.IsValidID(f.SourceFactionID) {
			return fmt.Errorf("job source faction id %s is invalid", f.SourceFactionID)
		}
		if !dbutils.IsValidID(f.SourceAreaID) {
			return fmt.Errorf("job source area id %s is invalid", f.SourceAreaID)
		}
		if !dbutils.IsValidID(f.TargetFactionID) {
			return fmt.Errorf("job target faction id %s is invalid", f.TargetFactionID)
		}
	}

	_, err := op.NamedExec(sqlInsertJobs, in)
	return err
}

func assignJob(op sqlOperator, jobID string, peopleIDs []string) (int, error) {
	if len(peopleIDs) == 0 || jobID == "" {
		return 0, nil
	}
	if !dbutils.IsValidID(jobID) {
		return 0, fmt.Errorf("job id %s is invalid", jobID)
	}
	for _, id := range peopleIDs {
		if !dbutils.IsValidID(id) {
			return 0, fmt.Errorf("person id %s is invalid", id)
		}
	}

	bindvar, args := sqlIn(peopleIDs)

	// note we specifiy here job_id = '' to ensure we can't clash with other inprogress job assignments
	peopleUpdate := fmt.Sprintf(`UPDATE %s SET job_id = $1 WHERE job_id = '' AND id IN (%s);`, tablePeople, bindvar)
	result, err := op.Exec(peopleUpdate, args...)
	if err != nil {
		return 0, err
	}

	// only update the job by how many rows we actually altered
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	jobUpdate := fmt.Sprintf(`UPDATE %s SET people_now = people_now + %d WHERE id = $1;`, tableJobs, affected)
	_, err = op.Exec(jobUpdate, jobID)
	return int(affected), err
}

func governments(op sqlOperator, token string, in *Query) ([]*structs.Government, string, error) {
	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	// 1. read Government objects
	q, args, err := genericSQLFromFilters(tk, in, tableGovernments)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Government
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	// 2. read law(s) for relevant Government(s)
	ids := make([]string, len(out))
	for i, g := range out {
		ids[i] = g.ID
	}

	laws, err := readLaws(op, ids)
	for _, g := range out {
		laws, ok := laws[g.ID]
		if !ok {
			// govt. has no laws
			laws = structs.NewLaws()
		}
		g.Outlawed = laws
	}

	return out, nextToken(tk, len(out)), nil
}

func readLaws(op sqlOperator, in []string) (map[string]*structs.Laws, error) {
	lawsQ, lawsArgs := sqlLawsFromSourceIDs(in)
	var lawRows []*lawStruct
	err := op.Select(&lawRows, lawsQ, lawsArgs...)
	if err != nil {
		return nil, err
	}

	result := map[string]*structs.Laws{}

	for _, l := range lawRows {
		law, ok := result[l.SourceID]
		if !ok {
			law = structs.NewLaws()
			result[l.SourceID] = law
		}

		switch l.MetaKey {
		case structs.MetaKeyFaction:
			law.Factions[l.MetaVal] = l.Illegal
		case structs.MetaKeyCommodity:
			law.Commodities[l.MetaVal] = l.Illegal
		case structs.MetaKeyAction:
			law.Actions[l.MetaVal] = l.Illegal
		case structs.MetaKeyResearch:
			law.Research[l.MetaVal] = l.Illegal
		case structs.MetaKeyReligion:
			law.Religions[l.MetaVal] = l.Illegal
		}
	}

	return result, nil
}

// toLaws converts a Government object to a slice of lawStructs.
// Internally we save the laws in their own table, but we don't expose this because we expect the laws
// to be reasonbly small in number and reasonably static.
func toLawRows(id string, laws *structs.Laws) []*lawStruct {
	rows := []*lawStruct{}
	if laws.Factions != nil {
		for k, v := range laws.Factions {
			if !v {
				continue
			}
			rows = append(rows, &lawStruct{SourceID: id, MetaKey: structs.MetaKeyFaction, MetaVal: k, Illegal: v})
		}
	}
	if laws.Commodities != nil {
		for k, v := range laws.Commodities {
			if !v {
				continue
			}
			rows = append(rows, &lawStruct{SourceID: id, MetaKey: structs.MetaKeyCommodity, MetaVal: k, Illegal: v})
		}
	}
	if laws.Actions != nil {
		for k, v := range laws.Actions {
			if !v {
				continue
			}
			rows = append(rows, &lawStruct{SourceID: id, MetaKey: structs.MetaKeyAction, MetaVal: string(k), Illegal: v})
		}
	}
	if laws.Research != nil {
		for k, v := range laws.Research {
			if !v {
				continue
			}
			rows = append(rows, &lawStruct{SourceID: id, MetaKey: structs.MetaKeyResearch, MetaVal: k, Illegal: v})
		}
	}
	if laws.Religions != nil {
		for k, v := range laws.Religions {
			if !v {
				continue
			}
			rows = append(rows, &lawStruct{SourceID: id, MetaKey: structs.MetaKeyReligion, MetaVal: k, Illegal: v})
		}
	}
	return rows
}

func setGovernments(op sqlOperator, in []*structs.Government) error {
	if len(in) == 0 {
		return nil
	}

	laws := []*lawStruct{}
	ids := make([]string, len(in))
	idNames := make([]string, len(in))
	idArgs := map[string]interface{}{}
	for i, f := range in {
		if !dbutils.IsValidID(f.ID) {
			return fmt.Errorf("government id %s is invalid", f.ID)
		}

		ids[i] = f.ID

		// seriously wtf why doesn't "sqlx.In" work ..
		idNames[i] = fmt.Sprintf(":%d", i)
		idArgs[fmt.Sprintf("%d", i)] = f.ID

		if f.Outlawed != nil {
			laws = append(laws, toLawRows(f.ID, f.Outlawed)...)
		}
	}

	// 1. write Government objects
	_, err := op.NamedExec(sqlInsertGovernments, in)
	if err != nil {
		return err
	}

	// 2. delete all laws for the given government(s)
	qstr := fmt.Sprintf(`DELETE FROM %s WHERE source_id in (%s);`, tableLaws, strings.Join(idNames, ","))

	_, err = op.NamedExec(qstr, idArgs)
	if err != nil {
		return err
	}

	// 3. annnd now we can write the laws
	_, err = op.NamedExec(sqlInsertLaws, laws)
	if err != nil {
		return err
	}

	return err
}

func families(op sqlOperator, token string, in *Query) ([]*structs.Family, string, error) {
	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := genericSQLFromFilters(tk, in, tableFamilies)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Family
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setFamilies(op sqlOperator, in []*structs.Family) error {
	if len(in) == 0 {
		return nil
	}

	for _, f := range in {
		if !dbutils.IsValidID(f.ID) {
			return fmt.Errorf("family id %s is invalid", f.ID)
		}
		if f.FactionID != "" && !dbutils.IsValidID(f.FactionID) {
			return fmt.Errorf("family faction id %s is invalid", f.FactionID)
		}
		if !dbutils.IsValidID(f.AreaID) {
			return fmt.Errorf("family area id %s is invalid", f.AreaID)
		}
		if !dbutils.IsValidID(f.MaleID) {
			return fmt.Errorf("family male id %s is invalid", f.MaleID)
		}
		if !dbutils.IsValidID(f.FemaleID) {
			return fmt.Errorf("family female id %s is invalid", f.FemaleID)
		}
		if f.Race == "" {
			return fmt.Errorf("family %s race is requred", f.ID)
		}
		if f.Culture == "" {
			return fmt.Errorf("family %s culture is requred", f.ID)
		}
		if f.Random == 0 {
			f.Random = rng.Intn(structs.RandomMax)
		}
	}

	_, err := op.NamedExec(sqlInsertFamilies, in)
	return err
}

func factionChildren(op sqlOperator, id string, rs []structs.FactionRelation) ([]*structs.Faction, error) {
	if !dbutils.IsValidID(id) {
		return nil, fmt.Errorf("faction id %s is invalid", id)
	}
	if rs == nil || len(rs) == 0 {
		rs = structs.AllFactionRelations
	}

	// TODO we could make this smarter with GreaterThan / LessThan but we have 4 relation types right now,
	// so it's not a huge deal
	in := []int{}
	for i, v := range rs {
		in[i] = int(v)
	}
	where, vals := sqlInInt(in)

	qstr := fmt.Sprintf(sqlFactionChildren, where, where)

	var out []*structs.Faction
	return out, op.Select(&out, qstr, append([]interface{}{id}, vals...)...)
}

func factionParents(op sqlOperator, id string, rs []structs.FactionRelation) ([]*structs.Faction, error) {
	if !dbutils.IsValidID(id) {
		return nil, fmt.Errorf("faction id %s is invalid", id)
	}
	if rs == nil || len(rs) == 0 {
		rs = structs.AllFactionRelations
	}

	// TODO we could make this smarter with GreaterThan / LessThan but we have 4 relation types right now,
	// so it's not a huge deal
	in := []int{}
	for i, v := range rs {
		in[i] = int(v)
	}
	where, vals := sqlInInt(in)

	qstr := fmt.Sprintf(sqlFactionParents, where, where)

	var out []*structs.Faction
	return out, op.Select(&out, qstr, append([]interface{}{id}, vals...)...)
}

func factions(op sqlOperator, token string, in *Query) ([]*structs.Faction, string, error) {
	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := genericSQLFromFilters(tk, in, tableFactions)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Faction
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setFactions(op sqlOperator, in []*structs.Faction) error {
	if len(in) == 0 {
		return nil
	}
	for _, f := range in {
		if !dbutils.IsValidID(f.ID) {
			return fmt.Errorf("faction id %s is invalid", f.ID)
		}
		if f.GovernmentID != "" && !dbutils.IsValidID(f.GovernmentID) {
			return fmt.Errorf("faction government id %s is invalid", f.GovernmentID)
		}
		if f.ParentFactionID != "" && !dbutils.IsValidID(f.ParentFactionID) {
			return fmt.Errorf("faction parent id %s is invalid", f.ParentFactionID)
		}
		if (f.IsReligion || f.ReligionID != "") && !dbutils.IsValidID(f.ReligionID) {
			return fmt.Errorf("faction religion id %s is invalid", f.ReligionID)
		}
		f.Clamp() // for ethos
		f.Cohesion = clampInt(f.Cohesion, 0, structs.MaxTuple)
		f.Corruption = clampInt(f.Corruption, 0, structs.MaxTuple)
		f.Wealth = clampInt(f.Wealth, 0, math.MaxInt64)
		f.EspionageOffense = clampInt(f.EspionageOffense, 0, math.MaxInt64)
		f.EspionageDefense = clampInt(f.EspionageDefense, 0, math.MaxInt64)
		f.MilitaryOffense = clampInt(f.MilitaryOffense, 0, math.MaxInt64)
		f.MilitaryDefense = clampInt(f.MilitaryDefense, 0, math.MaxInt64)
		f.Members = clampInt(f.Members, 0, math.MaxInt64)
		f.Plots = clampInt(f.Plots, 0, math.MaxInt64)
		f.Areas = clampInt(f.Areas, 0, math.MaxInt64)
	}

	_, err := op.NamedExec(sqlInsertFactions, in)
	return err
}

func areas(op sqlOperator, token string, in *Query) ([]*structs.Area, string, error) {
	if in == nil {
		in = Q()
	}

	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := genericSQLFromFilters(tk, in, tableAreas)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Area
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

// setAreas saves area information to the database
func setAreas(op sqlOperator, in []*structs.Area) error {
	if len(in) == 0 {
		return nil
	}
	for _, a := range in {
		if !dbutils.IsValidID(a.ID) {
			return fmt.Errorf("area id %s is invalid", a.ID)
		}
		if a.Random == 0 {
			a.Random = rng.Intn(structs.RandomMax)
		}
	}

	_, err := op.NamedExec(sqlInsertAreas, in)
	return err
}

// meta returns some metadata, if set
func meta(op sqlOperator, id string) (string, int, error) {
	if !dbutils.IsValidName(id) {
		return "", 0, fmt.Errorf("metadata key %s is invalid", id)
	}

	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE id=$1 LIMIT 1;",
		tableMeta,
	)

	result := []*mstruct{}
	err := op.Select(&result, query, id)
	if err != nil || len(result) == 0 {
		return "", 0, err
	}

	return result[0].Str, result[0].Int, nil
}

// setMeta sets some data in our meta table
func setMeta(op sqlOperator, id, strv string, intv int) error {
	if !dbutils.IsValidName(id) {
		return fmt.Errorf("metadata key %s is invalid", id)
	}

	// update schema version to current
	qstr := fmt.Sprintf(`INSERT INTO %s (id, str, int)
		VALUES (:id, :str, :int) 
		ON CONFLICT (id) DO UPDATE SET
		    int=EXCLUDED.int,
		    str=EXCLUDED.str
		;`,
		tableMeta,
	)
	_, err := op.NamedExec(qstr, map[string]interface{}{
		"id":  id,
		"str": strv,
		"int": intv,
	})
	return err
}

func nextToken(tk *dbutils.IterToken, rowsFetched int) string {
	if rowsFetched < tk.Limit {
		return ""
	}
	tk.Offset += tk.Limit
	return tk.String()
}
