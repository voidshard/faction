package db

import (
	"math"
	"strings"

	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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
}

// mstruct is a row of metadata
type mstruct struct {
	ID  string `db:"id"`
	Str string `db:"str"`
	Int int    `db:"int"`
}

// lawStruct holds a law row. We don't provide these as a first class object,
// but internally a government(s) laws are written as individual rows.
type lawStruct struct {
	GovernmentID string          `db:"government_id"`
	MetaKey      structs.MetaKey `db:"meta_key"`
	MetaVal      string          `db:"meta_val"`
	Illegal      bool            `db:"illegal"`
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

func modifiersSum(op sqlOperator, r Relation, token string, in []*Filter) ([]*structs.Tuple, string, error) {
	if !r.SupportsModifiers() {
		return nil, token, fmt.Errorf("relation %s does not support modifiers", r)
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

func modifiers(op sqlOperator, r Relation, token string, in []*Filter) ([]*structs.Modifier, string, error) {
	if !r.SupportsModifiers() {
		return nil, token, fmt.Errorf("relation %s does not support modifiers", r)
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

func incrModifiers(op sqlOperator, r Relation, v int, in []*Filter) error {
	if len(in) == 0 {
		return nil
	}

	where, args, err := sqlWhereFromFilters(in, 1)
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

func tuples(op sqlOperator, r Relation, token string, in []*Filter) ([]*structs.Tuple, string, error) {
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

func incrTuples(op sqlOperator, r Relation, v int, in []*Filter) error {
	if len(in) == 0 {
		return nil
	}

	where, args, err := sqlWhereFromFilters(in, 0)
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

func plots(op sqlOperator, token string, in []*Filter) ([]*structs.Plot, string, error) {
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

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    id, area_id, faction_id, size, commodity, yield
	) VALUES (
	    :id, :area_id, :faction_id, :size, :commodity, :yield
	) ON CONFLICT (id) DO UPDATE SET 
	    faction_id=EXCLUDED.faction_id,
	    size=EXCLUDED.size,
	    commodity=EXCLUDED.commodity,
	    yield=EXCLUDED.yield
	;`, tablePlots)

	_, err := op.NamedExec(qstr, in)
	return err
}

func routes(op sqlOperator, token string, in []*Filter) ([]*structs.Route, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := sqlFromRouteFilters(tk, in)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Route
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setRoutes(op sqlOperator, in []*structs.Route) error {
	if len(in) == 0 {
		return nil
	}

	for _, f := range in {
		if !dbutils.IsValidID(f.SourceAreaID) {
			return fmt.Errorf("route source area id %s is invalid", f.SourceAreaID)
		}
		if !dbutils.IsValidID(f.TargetAreaID) {
			return fmt.Errorf("route target area id %s is invalid", f.TargetAreaID)
		}
		if f.SourceAreaID == f.TargetAreaID {
			return fmt.Errorf("source and target area ids are the same: %s", f.SourceAreaID)
		}
		if f.TravelTime < 0 { // we do not allow time travel, but instantaneous teleportation is fine
			return fmt.Errorf("invalid travel time, expected >= 0: %d", f.TravelTime)
		}
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    source_area_id, target_area_id, travel_time
	) VALUES (
	    :source_area_id, :target_area_id, :travel_time
	) ON CONFLICT (source_area_id, target_area_id) DO UPDATE SET
	    travel_time=EXCLUDED.travel_time;`, tableRoutes)

	_, err := op.NamedExec(qstr, in)
	return err
}

func people(op sqlOperator, token string, in []*Filter) ([]*structs.Person, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := sqlFromPersonFilters(tk, in)
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
		f.Clamp()
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    id,
	    ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
	    first_name, last_name, birth_family_id, race,
	    area_id, job_id,
	    birth_tick, death_tick, is_male, is_child,
	    preferred_profession, preferred_faction_id,
	    death_meta_reason, death_meta_key, death_meta_val
	) VALUES (
	    :id,
	    :ethos_altruism, :ethos_ambition, :ethos_tradition, :ethos_pacifism, :ethos_piety, :ethos_caution,
	    :first_name, :last_name, :birth_family_id, :race,
	    :area_id, :job_id,
	    :birth_tick, :death_tick,
	    :is_male, :is_child,
	    :preferred_profession, :preferred_faction_id,
	    :death_meta_reason, :death_meta_key, :death_meta_val
	) ON CONFLICT (id) DO UPDATE SET
	    ethos_altruism=EXCLUDED.ethos_altruism,
	    ethos_ambition=EXCLUDED.ethos_ambition,
	    ethos_tradition=EXCLUDED.ethos_tradition,
	    ethos_pacifism=EXCLUDED.ethos_pacifism,
	    ethos_piety=EXCLUDED.ethos_piety,
	    ethos_caution=EXCLUDED.ethos_caution,
	    area_id=EXCLUDED.area_id,
	    job_id=EXCLUDED.job_id,
	    is_child=EXCLUDED.is_child,
	    preferred_profession=EXCLUDED.preferred_profession,
	    preferred_faction_id=EXCLUDED.preferred_faction_id,
	    death_tick=EXCLUDED.death_tick,
	    death_meta_reason=EXCLUDED.death_meta_reason,
	    death_meta_key=EXCLUDED.death_meta_key,
	    death_meta_val=EXCLUDED.death_meta_val
	;`, tablePeople)

	_, err := op.NamedExec(qstr, in)
	return err
}

func jobs(op sqlOperator, token string, in []*Filter) ([]*structs.Job, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := sqlFromJobFilters(tk, in)
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
	}

	// we only need to update the job state.
	qstr := fmt.Sprintf(`INSERT INTO %s (
	    id, parent_job_id,
	    source_faction_id, source_area_id,
	    action,
	    target_area_id, target_meta_key, target_meta_val,
	    people_min, people_max,
	    tick_created, tick_starts, tick_ends,
	    secrecy,
	    is_illegal,
	    state
	) VALUES (
	    :id, :parent_job_id,
	    :source_faction_id, :source_area_id,
	    :action,
	    :target_area_id, :target_meta_key, :target_meta_val,
	    :people_min, :people_max,
	    :tick_created, :tick_starts, :tick_ends,
	    :secrecy,
	    :is_illegal,
	    :state
	) ON CONFLICT (id) DO UPDATE SET 
	   state=EXCLUDED.state
	;`, tableJobs)

	_, err := op.NamedExec(qstr, in)
	return err
}

func governments(op sqlOperator, token string, in []*Filter) ([]*structs.Government, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	// 1. read Government objects
	q, args, err := sqlFromGovernmentFilters(tk, in)
	if err != nil {
		return nil, token, err
	}

	var out []*structs.Government
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	// 2. read law(s) for relevant Government(s)
	gmap := map[string]*structs.Government{}
	ids := make([]string, len(out))
	for i, g := range out {
		ids[i] = g.ID
		gmap[g.ID] = g
		g.Outlawed = structs.NewLaws()
	}

	lawsQ, lawsArgs := sqlLawsFromGovernmentIDs(ids)

	var laws []*lawStruct
	err = op.Select(&laws, lawsQ, lawsArgs...)
	if err != nil {
		return nil, token, err
	}

	// 3. cobble together
	for _, l := range laws {
		g, ok := gmap[l.GovernmentID]
		if !ok { // should never happen, must have picked up row we didn't mean to
			continue
		}

		switch l.MetaKey {
		case structs.MetaKeyFaction:
			g.Outlawed.Factions[l.MetaVal] = l.Illegal
		case structs.MetaKeyCommodity:
			g.Outlawed.Commodities[l.MetaVal] = l.Illegal
		case structs.MetaKeyAction:
			g.Outlawed.Actions[structs.ActionType(l.MetaVal)] = l.Illegal
		case structs.MetaKeyResearch:
			g.Outlawed.Research[l.MetaVal] = l.Illegal
		case structs.MetaKeyReligion:
			g.Outlawed.Religions[l.MetaVal] = l.Illegal
		}
	}

	return out, nextToken(tk, len(out)), nil
}

// toLaws converts a Government object to a slice of lawStructs.
// Internally we save the laws in their own table, but we don't expose this because we expect the laws
// to be reasonbly small in number and reasonably static.
func toLaws(f *structs.Government) []*lawStruct {
	laws := []*lawStruct{}
	if f.Outlawed.Factions != nil {
		for k, v := range f.Outlawed.Factions {
			if !v {
				continue
			}
			laws = append(laws, &lawStruct{GovernmentID: f.ID, MetaKey: structs.MetaKeyFaction, MetaVal: k, Illegal: v})
		}
	}
	if f.Outlawed.Commodities != nil {
		for k, v := range f.Outlawed.Commodities {
			if !v {
				continue
			}
			laws = append(laws, &lawStruct{GovernmentID: f.ID, MetaKey: structs.MetaKeyCommodity, MetaVal: k, Illegal: v})
		}
	}
	if f.Outlawed.Actions != nil {
		for k, v := range f.Outlawed.Actions {
			if !v {
				continue
			}
			laws = append(laws, &lawStruct{GovernmentID: f.ID, MetaKey: structs.MetaKeyAction, MetaVal: string(k), Illegal: v})
		}
	}
	if f.Outlawed.Research != nil {
		for k, v := range f.Outlawed.Research {
			if !v {
				continue
			}
			laws = append(laws, &lawStruct{GovernmentID: f.ID, MetaKey: structs.MetaKeyResearch, MetaVal: k, Illegal: v})
		}
	}
	if f.Outlawed.Religions != nil {
		for k, v := range f.Outlawed.Religions {
			if !v {
				continue
			}
			laws = append(laws, &lawStruct{GovernmentID: f.ID, MetaKey: structs.MetaKeyReligion, MetaVal: k, Illegal: v})
		}
	}
	return laws
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
			laws = append(laws, toLaws(f)...)
		}
	}

	// 1. write Government objects
	qstr := fmt.Sprintf(`INSERT INTO %s (
	    id, tax_rate, tax_frequency
	) VALUES (
	    :id, :tax_rate, :tax_frequency
	) ON CONFLICT (id) DO UPDATE SET 
	    tax_rate=EXCLUDED.tax_rate,
	    tax_frequency=EXCLUDED.tax_frequency
	;`, tableGovernments)

	_, err := op.NamedExec(qstr, in)
	if err != nil {
		return err
	}

	// 2. delete all laws for the given government(s)
	qstr = fmt.Sprintf(`DELETE FROM %s WHERE government_id in (%s);`, tableLaws, strings.Join(idNames, ","))

	_, err = op.NamedExec(qstr, idArgs)
	if err != nil {
		return err
	}

	// 3. annnd now we can write the laws
	qstr = fmt.Sprintf(`INSERT INTO %s (
	    government_id, meta_key, meta_val, illegal
	) VALUES (
	    :government_id, :meta_key, :meta_val, :illegal
	);`, tableLaws)

	_, err = op.NamedExec(qstr, laws)
	if err != nil {
		return err
	}

	return err
}

func families(op sqlOperator, token string, in []*Filter) ([]*structs.Family, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := sqlFromFamilyFilters(tk, in)
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
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    id, area_id, faction_id, 
	    ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
	    is_child_bearing, max_child_bearing_tick,  pregnancy_end,
	    male_id, female_id,
	    ma_grandma_id, ma_grandpa_id, pa_grandma_id, pa_grandpa_id,
	    number_of_children
	) VALUES (
	    :id, :area_id, :faction_id, 
	    :ethos_altruism, :ethos_ambition, :ethos_tradition, :ethos_pacifism, :ethos_piety, :ethos_caution,
	    :is_child_bearing, :max_child_bearing_tick,  :pregnancy_end,
	    :male_id, :female_id,
	    :ma_grandma_id, :ma_grandpa_id, :pa_grandma_id, :pa_grandpa_id,
	    :number_of_children
	) ON CONFLICT (id) DO UPDATE SET
	    ethos_altruism=EXCLUDED.ethos_altruism,
	    ethos_ambition=EXCLUDED.ethos_ambition,
	    ethos_tradition=EXCLUDED.ethos_tradition,
	    ethos_pacifism=EXCLUDED.ethos_pacifism,
	    ethos_piety=EXCLUDED.ethos_piety,
	    ethos_caution=EXCLUDED.ethos_caution,
	    area_id=EXCLUDED.area_id,
	    faction_id=EXCLUDED.faction_id,
	    is_child_bearing=EXCLUDED.is_child_bearing,
	    pregnancy_end=EXCLUDED.pregnancy_end,
	    number_of_children=EXCLUDED.number_of_children
	;`, tableFamilies)

	_, err := op.NamedExec(qstr, in)
	return err
}

func factions(op sqlOperator, token string, in []*Filter) ([]*structs.Faction, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := sqlFromFactionFilters(tk, in)
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
	}

	// We could make this shorter, but I like to be very specific in SQL :P
	qstr := fmt.Sprintf(`INSERT INTO %s (
            id, name, home_area_id,
	    ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
            action_frequency_ticks,
            leadership, structure, wealth, cohesion, corruption,
	    is_covert,
	    government_id, is_government,
	    religion_id, is_religion,
	    is_member_by_birth,
	    espionage_offense, espionage_defense,
	    military_offense, military_defense,
            parent_faction_id,
	    parent_faction_relation
	) VALUES (
	    :id, :name,
	    :home_area_id,
	    :ethos_altruism, :ethos_ambition, :ethos_tradition, :ethos_pacifism, :ethos_piety, :ethos_caution,
	    :action_frequency_ticks,
	    :leadership, :structure, :wealth, :cohesion, :corruption,
	    :is_covert,
	    :government_id, :is_government,
	    :religion_id, :is_religion,
	    :is_member_by_birth,
	    :espionage_offense, :espionage_defense,
	    :military_offense, :military_defense,
	    :parent_faction_id,
	    :parent_faction_relation
	) ON CONFLICT (id) DO UPDATE SET
	    name=EXCLUDED.name,
	    home_area_id=EXCLUDED.home_area_id,
	    ethos_altruism=EXCLUDED.ethos_altruism,
	    ethos_ambition=EXCLUDED.ethos_ambition,
	    ethos_tradition=EXCLUDED.ethos_tradition,
	    ethos_pacifism=EXCLUDED.ethos_pacifism,
	    ethos_piety=EXCLUDED.ethos_piety,
	    ethos_caution=EXCLUDED.ethos_caution,
	    action_frequency_ticks=EXCLUDED.action_frequency_ticks,
	    leadership=EXCLUDED.leadership,
	    structure=EXCLUDED.structure,
	    wealth=EXCLUDED.wealth,
	    cohesion=EXCLUDED.cohesion,
	    corruption=EXCLUDED.corruption,
	    is_covert=EXCLUDED.is_covert,
	    government_id=EXCLUDED.government_id,
	    is_government=EXCLUDED.is_government,
	    religion_id=EXCLUDED.religion_id,
	    is_religion=EXCLUDED.is_religion,
	    is_member_by_birth=EXCLUDED.is_member_by_birth,
	    espionage_offense=EXCLUDED.espionage_offense,
	    espionage_defense=EXCLUDED.espionage_defense,
	    military_offense=EXCLUDED.military_offense,
	    military_defense=EXCLUDED.military_defense,
	    parent_faction_id=EXCLUDED.parent_faction_id,
	    parent_faction_relation=EXCLUDED.parent_faction_relation
	;`, tableFactions)

	_, err := op.NamedExec(qstr, in)
	return err
}

func areas(op sqlOperator, token string, in []*Filter) ([]*structs.Area, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args, err := sqlFromAreaFilters(tk, in)
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
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (id, government_id)
		        VALUES (:id, :government_id)
		        ON CONFLICT (id) DO UPDATE SET
		            government_id=EXCLUDED.government_id
		        ;`,
		tableAreas,
	)

	_, err := op.NamedExec(qstr, in)
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
