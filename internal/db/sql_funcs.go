package db

import (
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
	Select(dest interface{}, query string, args ...interface{}) error
}

// mstruct is a row of metadata
type mstruct struct {
	ID  string `db:"id"`
	Str string `db:"str"`
	Int int    `db:"int"`
}

func deleteModifiers(op sqlOperator, r Relation, expires_before_tick int) error {
	if !r.supportsModifiers() {
		return nil
	}

	qstr := fmt.Sprintf(`DELETE FROM %s WHERE tick_expires < :time;`, r.modTable())

	_, err := op.NamedExec(qstr, map[string]interface{}{
		"time": expires_before_tick,
	})
	return err
}

func modifiers(op sqlOperator, r Relation, token string, in []*ModifierFilter) ([]*structs.Modifier, string, error) {
	if !r.supportsModifiers() {
		return nil, token, fmt.Errorf("relation %s does not support modifiers", r)
	}
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromModifierFilters(r, tk, in)

	var out []*structs.Modifier
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setModifiers(op sqlOperator, r Relation, in []*structs.Modifier) error {
	if !r.supportsModifiers() {
		return fmt.Errorf("relation %s does not support modifiers", r)
	}
	if len(in) == 0 {
		return nil
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    subject, object, value, tick_expires, meta_key, meta_val, meta_reason
	) VALUES (
	    :subject, :object, :value, :tick_expires, :meta_key, :meta_val, :meta_reason
	);`, r.modTable())

	_, err := op.NamedExec(qstr, in)
	return err
}

func tuples(op sqlOperator, r Relation, token string, in []*TupleFilter) ([]*structs.Tuple, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromTupleFilters(r, tk, in)

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

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    subject, object, value
	) VALUES (
	    :subject, :object, :value
	) ON CONFLICT (subject, object) DO UPDATE SET value=EXCLUDED.value;`, r.tupleTable())

	_, err := op.NamedExec(qstr, in)
	return err
}

func plots(op sqlOperator, token string, in []*PlotFilter) ([]*structs.Plot, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromPlotFilters(tk, in)

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
		if f.OwnerFactionID != "" && !dbutils.IsValidID(f.OwnerFactionID) {
			return fmt.Errorf("plot owner faction id %s is invalid", f.OwnerFactionID)
		}
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    id, is_headquarters, area_id, owner_faction_id, size
	) VALUES (
	    :id, :is_headquarters, :area_id, :owner_faction_id, :size
	) ON CONFLICT (id) DO UPDATE SET 
	    is_headquarters=EXCLUDED.is_headquarters,
	    owner_faction_id=EXCLUDED.owner_faction_id,
	    size=EXCLUDED.size
	;`, tablePlots)

	_, err := op.NamedExec(qstr, in)
	return err
}

func routes(op sqlOperator, token string, in []*RouteFilter) ([]*structs.Route, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromRouteFilters(tk, in)

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

func people(op sqlOperator, token string, in []*PersonFilter) ([]*structs.Person, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromPersonFilters(tk, in)

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
		if !dbutils.IsValidID(f.BirthFamilyID) {
			return fmt.Errorf("person birth family id %s is invalid", f.BirthFamilyID)
		}
		if !dbutils.IsValidID(f.AreaID) {
			return fmt.Errorf("person area id %s is invalid", f.AreaID)
		}
		if f.JobID != "" && !dbutils.IsValidID(f.JobID) {
			return fmt.Errorf("person job id %s is invalid", f.JobID)
		}
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    id,
	    ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
	    first_name, last_name, birth_family_id, race,
	    area_id, job_id,
	    birth_tick, death_tick, is_male
	) VALUES (
	    :id,
	    :ethos_altruism, :ethos_ambition, :ethos_tradition, :ethos_pacifism, :ethos_piety, :ethos_caution,
	    :first_name, :last_name, :birth_family_id, :race,
	    :area_id, :job_id,
	    :birth_tick, :death_tick, :is_male
	) ON CONFLICT (id) DO UPDATE SET
	    ethos_altruism=EXCLUDED.ethos_altruism,
	    ethos_ambition=EXCLUDED.ethos_ambition,
	    ethos_tradition=EXCLUDED.ethos_tradition,
	    ethos_pacifism=EXCLUDED.ethos_pacifism,
	    ethos_piety=EXCLUDED.ethos_piety,
	    ethos_caution=EXCLUDED.ethos_caution,
	    area_id=EXCLUDED.area_id,
	    job_id=EXCLUDED.job_id,
	    death_tick=EXCLUDED.death_tick
	;`, tablePeople)

	_, err := op.NamedExec(qstr, in)
	return err
}

func landRights(op sqlOperator, token string, in []*LandRightFilter) ([]*structs.LandRight, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromLandRightFilters(tk, in)

	var out []*structs.LandRight
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setLandRights(op sqlOperator, in []*structs.LandRight) error {
	if len(in) == 0 {
		return nil
	}

	for _, f := range in {
		if !dbutils.IsValidID(f.ID) {
			return fmt.Errorf("land right id %s is invalid", f.ID)
		}
		if f.GoverningFactionID != "" && !dbutils.IsValidID(f.GoverningFactionID) {
			return fmt.Errorf("land right governing faction id %s is invalid", f.GoverningFactionID)
		}
		if f.ControllingFactionID != "" && !dbutils.IsValidID(f.ControllingFactionID) {
			return fmt.Errorf("land right controlling faction id %s is invalid", f.ControllingFactionID)
		}
		if !dbutils.IsValidID(f.AreaID) {
			return fmt.Errorf("land right area id %s is invalid", f.AreaID)
		}
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    id, governing_faction_id, controlling_faction_id, area_id, resource
	) VALUES (
	    :id, :governing_faction_id, :controlling_faction_id, :area_id, :resource
	) ON CONFLICT (id) DO UPDATE SET
	    governing_faction_id = EXCLUDED.governing_faction_id,
	    controlling_faction_id = EXCLUDED.controlling_faction_id
	;`, tableLandRights)

	_, err := op.NamedExec(qstr, in)
	return err
}

func jobs(op sqlOperator, token string, in []*JobFilter) ([]*structs.Job, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromJobFilters(tk, in)

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
	    id,
	    source_faction_id, source_area_id,
	    action,
	    target_area_id, target_meta_key, target_meta_val,
	    people_min, people_max,
	    tick_created, tick_starts, tick_ends,
	    secrecy,
	    is_illegal,
	    state
	) VALUES (
	    :id,
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

func governments(op sqlOperator, token string, in []*GovernmentFilter) ([]*structs.Government, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromGovernmentFilters(tk, in)

	var out []*structs.Government
	err = op.Select(&out, q, args...)
	if err != nil {
		return nil, token, err
	}

	return out, nextToken(tk, len(out)), nil
}

func setGovernments(op sqlOperator, in []*structs.Government) error {
	if len(in) == 0 {
		return nil
	}

	for _, f := range in {
		if !dbutils.IsValidID(f.ID) {
			return fmt.Errorf("government id %s is invalid", f.ID)
		}
	}

	qstr := fmt.Sprintf(`INSERT INTO %s (
	    id, tax_rate, tax_frequency
	) VALUES (
	    :id, :tax_rate, :tax_frequency
	) ON CONFLICT (id) DO UPDATE SET 
	    tax_rate=EXCLUDED.tax_rate,
	    tax_frequency=EXCLUDED.tax_frequency
	;`, tableGovernments)

	_, err := op.NamedExec(qstr, in)
	return err
}

func families(op sqlOperator, token string, in []*FamilyFilter) ([]*structs.Family, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromFamilyFilters(tk, in)

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
	    id, area_id, faction_id, is_child_bearing, male_id, female_id 
	) VALUES (
	    :id, :area_id, :faction_id, :is_child_bearing, :male_id, :female_id 
	) ON CONFLICT (id) DO UPDATE SET
	    area_id=EXCLUDED.area_id,
	    faction_id=EXCLUDED.faction_id,
	    is_child_bearing=EXCLUDED.is_child_bearing
	;`, tableFamilies)

	_, err := op.NamedExec(qstr, in)
	return err
}

func factions(op sqlOperator, token string, in []*FactionFilter) ([]*structs.Faction, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromFactionFilters(tk, in)

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
		if !dbutils.IsValidID(f.GovernmentID) {
			return fmt.Errorf("faction government id %s is invalid", f.GovernmentID)
		}
		if f.ParentFactionID != "" && !dbutils.IsValidID(f.ParentFactionID) {
			return fmt.Errorf("faction parent id %s is invalid", f.ParentFactionID)
		}
		if f.ReligionID != "" && !dbutils.IsValidID(f.ReligionID) {
			return fmt.Errorf("faction religion id %s is invalid", f.ReligionID)
		}
	}

	// We could make this shorter, but I like to be very specific in SQL :P
	qstr := fmt.Sprintf(`INSERT INTO %s (
            id, name,
	    ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
            action_frequency_ticks,
            leadership, wealth, cohesion, corruption,
	    is_covert, is_illegal,
	    government_id, is_government,
	    religion_id, is_religion,
	    is_member_by_birth,
	    espionage_offense, espionage_defense,
	    military_offense, military_defense,
            parent_faction_id,
	    parent_faction_relation
	) VALUES (
	    :id, :name,
	    :ethos_altruism, :ethos_ambition, :ethos_tradition, :ethos_pacifism, :ethos_piety, :ethos_caution,
	    :action_frequency_ticks,
	    :leadership, :wealth, :cohesion, :corruption,
	    :is_covert, :is_illegal,
	    :government_id, :is_government,
	    :religion_id, :is_religion,
	    :is_member_by_birth,
	    :espionage_offense, :espionage_defense,
	    :military_offense, :military_defense,
	    :parent_faction_id,
	    :parent_faction_relation
	) ON CONFLICT (id) DO UPDATE SET
	    name=EXCLUDED.name,
	    ethos_altruism=EXCLUDED.ethos_altruism,
	    ethos_ambition=EXCLUDED.ethos_ambition,
	    ethos_tradition=EXCLUDED.ethos_tradition,
	    ethos_pacifism=EXCLUDED.ethos_pacifism,
	    ethos_piety=EXCLUDED.ethos_piety,
	    ethos_caution=EXCLUDED.ethos_caution,
	    action_frequency_ticks=EXCLUDED.action_frequency_ticks,
	    leadership=EXCLUDED.leadership,
	    wealth=EXCLUDED.wealth,
	    cohesion=EXCLUDED.cohesion,
	    corruption=EXCLUDED.corruption,
	    is_covert=EXCLUDED.is_covert,
	    is_illegal=EXCLUDED.is_illegal,
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

func areas(op sqlOperator, token string, in []*AreaFilter) ([]*structs.Area, string, error) {
	tk, err := dbutils.ParseToken(token)
	if err != nil {
		return nil, token, err
	}

	q, args := sqlFromAreaFilters(tk, in)

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

	qstr := fmt.Sprintf(`INSERT INTO %s (id, governing_faction_id)
		        VALUES (:id, :governing_faction_id)
		        ON CONFLICT (id) DO UPDATE SET
		            governing_faction_id=EXCLUDED.governing_faction_id
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
