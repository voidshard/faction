/*
	Sqlite3 we use as a DB for users wanting to run in local mode.

It is not recommended to use this for anything other than testing, or very small simulations.

We keep the SQL here compatible with Postgres.
*/
package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/voidshard/faction/pkg/config"
)

const (
	// metaSchemaVersion is the current db schema version.
	// This applies only for sqlite
	metaSchemaVersion = "db-schema-version"

	// currentSchemaVersion of the db schema. Should be updated
	// when we update the tables so we can handle migrations
	currentSchemaVersion = 1

	tupleTableCreateTemplate = `CREATE TABLE IF NOT EXISTS %s (
	    subject VARCHAR(255) NOT NULL, 
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    UNIQUE(subject, object)
	);`

	modifierTableCreateTemplate = `CREATE TABLE IF NOT EXISTS %s (
	    subject VARCHAR(255) NOT NULL,
	    object VARCHAR(255) NOT NULL,
	    value INTEGER NOT NULL DEFAULT 0,
	    tick_expires INTEGER NOT NULL DEFAULT 0,
	    meta_key VARCHAR(255) DEFAULT '',
	    meta_val VARCHAR(255) DEFAULT '',
	    meta_reason TEXT DEFAULT ''
	);`

	tupleIndexCreateTemplate = `CREATE INDEX IF NOT EXISTS subject_%s ON %s (subject);`
)

var (
	// NB. UUIDs are 36 chars (eg. 123e4567-e89b-12d3-a456-426655440000)
	createMeta = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	    id VARCHAR(255) PRIMARY KEY,
	    str VARCHAR(255) DEFAULT '',
	    int INTEGER NOT NULL DEFAULT 0
	);`, tableMeta)

	createArea = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	    id VARCHAR(36) PRIMARY KEY,
	    government_id VARCHAR(36),
	    biome VARCHAR(255) DEFAULT '',
	    random INTEGER NOT NULL DEFAULT 0
	);`, tableAreas)

	createGovernments = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	    id VARCHAR(36) PRIMARY KEY,
	    tax_rate REAL NOT NULL DEFAULT 0.10,
	    tax_frequency INTEGER NOT NULL DEFAULT 1
	);`, tableGovernments)

	createLaws = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	    source_id VARCHAR(36) NOT NULL,
	    meta_key VARCHAR(255) NOT NULL,
	    meta_val VARCHAR(255) NOT NULL,
	    illegal BOOLEAN NOT NULL DEFAULT FALSE,
	    UNIQUE(source_id, meta_key, meta_val)
	);`, tableLaws)

	createPlots = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	    id VARCHAR(36) PRIMARY KEY,
	    area_id VARCHAR(36) NOT NULL,
	    faction_id VARCHAR(36) DEFAULT '',
	    hidden INTEGER NOT NULL DEFAULT 0,
	    value REAL NOT NULL DEFAULT 0,
	    size INTEGER NOT NULL DEFAULT 1,
	    commodity VARCHAR(255) DEFAULT '',
	    yield INTEGER NOT NULL DEFAULT 0
	);`, tablePlots)

	createJobs = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	    id VARCHAR(36) PRIMARY KEY,
	    parent_job_id VARCHAR(36) DEFAULT '',
	    source_faction_id VARCHAR(36) NOT NULL,
	    source_area_id VARCHAR(36) NOT NULL,
	    action VARCHAR(255) NOT NULL,
	    priority INTEGER NOT NULL DEFAULT 0,
	    conscription BOOLEAN NOT NULL DEFAULT FALSE,
	    target_area_id VARCHAR(36) NOT NULL,
	    target_faction_id VARCHAR(36) DEFAULT '',
	    target_meta_key VARCHAR(20) DEFAULT '',
	    target_meta_val VARCHAR(255) DEFAULT '',
	    people_min INTEGER NOT NULL DEFAULT 1,
	    people_max INTEGER NOT NULL DEFAULT 0,
	    people_now INTEGER NOT NULL DEFAULT 0,
	    tick_created INTEGER NOT NULL DEFAULT 0,
	    tick_starts INTEGER NOT NULL DEFAULT 0,
	    tick_ends INTEGER NOT NULL DEFAULT 1,
	    secrecy INTEGER NOT NULL DEFAULT 0,
	    is_illegal BOOLEAN NOT NULL DEFAULT FALSE,
	    state VARCHAR(20) NOT NULL DEFAULT 'pending'
	);`, tableJobs)

	createFamilies = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	    id VARCHAR(36) PRIMARY KEY,
	    race VARCHAR(255) DEFAULT '',
	    culture VARCHAR(255) DEFAULT '',
            ethos_altruism INTEGER NOT NULL DEFAULT 0,
            ethos_ambition INTEGER NOT NULL DEFAULT 0,
            ethos_tradition INTEGER NOT NULL DEFAULT 0,
            ethos_pacifism INTEGER NOT NULL DEFAULT 0,
            ethos_piety INTEGER NOT NULL DEFAULT 0,
            ethos_caution INTEGER NOT NULL DEFAULT 0,
	    area_id VARCHAR(36) NOT NULL,
	    social_class VARCHAR(255) DEFAULT '',
	    faction_id VARCHAR(36) DEFAULT '',
	    is_child_bearing BOOLEAN NOT NULL DEFAULT FALSE,
	    male_id VARCHAR(36) NOT NULL,
	    female_id VARCHAR(36) NOT NULL,
	    max_child_bearing_tick INTEGER NOT NULL DEFAULT 0,
	    pregnancy_end INTEGER NOT NULL DEFAULT 0,
	    ma_grandma_id VARCHAR(36) DEFAULT '',
	    ma_grandpa_id VARCHAR(36) DEFAULT '',
	    pa_grandma_id VARCHAR(36) DEFAULT '',
	    pa_grandpa_id VARCHAR(36) DEFAULT '',
	    number_of_children INTEGER NOT NULL DEFAULT 0,
	    random INTEGER NOT NULL DEFAULT 0,
	    UNIQUE (male_id, female_id)
	);`, tableFamilies)

	createPeople = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
            id VARCHAR(36) PRIMARY KEY,
	    first_name VARCHAR(255) DEFAULT '',
	    last_name VARCHAR(255) DEFAULT '',
	    birth_family_id VARCHAR(36) DEFAULT '',
	    race VARCHAR(255) DEFAULT '',
	    culture VARCHAR(255) DEFAULT '',
            ethos_altruism INTEGER NOT NULL DEFAULT 0,
            ethos_ambition INTEGER NOT NULL DEFAULT 0,
            ethos_tradition INTEGER NOT NULL DEFAULT 0,
            ethos_pacifism INTEGER NOT NULL DEFAULT 0,
            ethos_piety INTEGER NOT NULL DEFAULT 0,
            ethos_caution INTEGER NOT NULL DEFAULT 0,
	    area_id VARCHAR(36) NOT NULL,
	    job_id VARCHAR(36) DEFAULT '',
	    birth_tick INTEGER NOT NULL DEFAULT 1,
	    death_tick INTEGER NOT NULL DEFAULT 0,
	    is_male BOOLEAN NOT NULL DEFAULT FALSE,
	    adulthood_tick INTEGER NOT NULL DEFAULT 0,
	    preferred_profession VARCHAR(255) DEFAULT '',
	    preferred_faction_id VARCHAR(36) DEFAULT '',
	    death_meta_reason TEXT DEFAULT '',
	    death_meta_key VARCHAR(255) DEFAULT '',
	    death_meta_val VARCHAR(255) DEFAULT '',
	    natural_death_tick INTEGER NOT NULL DEFAULT 0,
	    random INTEGER NOT NULL DEFAULT 0
	);`, tablePeople)

	createFactions = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
            id VARCHAR(36) PRIMARY KEY,
	    name VARCHAR(255) DEFAULT '',
	    home_area_id VARCHAR(36) NOT NULL,
	    hq_plot_id VARCHAR(36) NOT NULL,
            ethos_altruism INTEGER NOT NULL DEFAULT 0,
            ethos_ambition INTEGER NOT NULL DEFAULT 0,
            ethos_tradition INTEGER NOT NULL DEFAULT 0,
            ethos_pacifism INTEGER NOT NULL DEFAULT 0,
            ethos_piety INTEGER NOT NULL DEFAULT 0,
            ethos_caution INTEGER NOT NULL DEFAULT 0,
            action_frequency_ticks INTEGER NOT NULL DEFAULT 1,
            leadership INTEGER NOT NULL DEFAULT 0,
	    structure INTEGER NOT NULL DEFAULT 0,
            wealth INTEGER NOT NULL DEFAULT 0,
            cohesion INTEGER NOT NULL DEFAULT 0,
            corruption INTEGER NOT NULL DEFAULT 0,
            is_covert BOOLEAN NOT NULL DEFAULT FALSE,
            government_id VARCHAR(36) NOT NULL,
            is_government BOOLEAN NOT NULL DEFAULT FALSE,
            religion_id VARCHAR(36) DEFAULT '',
            is_religion BOOLEAN NOT NULL DEFAULT FALSE,
            is_member_by_birth BOOLEAN NOT NULL DEFAULT FALSE,
	    espionage_offense INTEGER NOT NULL DEFAULT 0,
	    espionage_defense INTEGER NOT NULL DEFAULT 0,
	    military_offense INTEGER NOT NULL DEFAULT 0,
	    military_defense INTEGER NOT NULL DEFAULT 0,
            parent_faction_id VARCHAR(36) DEFAULT '',
            parent_faction_relation INTEGER NOT NULL DEFAULT 0,
	    members INTEGER NOT NULL DEFAULT 0,
	    vassals INTEGER NOT NULL DEFAULT 0,
	    plots INTEGER NOT NULL DEFAULT 0,
	    areas INTEGER NOT NULL DEFAULT 0
        );`, tableFactions)

	createEvents = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	    id VARCHAR(36) PRIMARY KEY,
	    type VARCHAR(255) DEFAULT '',
	    tick INTEGER NOT NULL DEFAULT 0,
	    message TEXT DEFAULT '',
	    subject_meta_key VARCHAR(255) DEFAULT '',
	    subject_meta_val VARCHAR(255) DEFAULT '',
	    cause_meta_key VARCHAR(255) DEFAULT '',
	    cause_meta_val VARCHAR(255) DEFAULT ''
	);`, tableEvents)

	// indexes that we should create
	// TODO
	indexes = []string{
		// we look up laws by the government id
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS laws_source_id ON %s (source_id);`, tableLaws),

		// we look up plots by their owner(s) in order to determine where faction(s) can perform actions
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS plot_owner ON %s (faction_id);`, tablePlots),

		// jobs we want to search by area & state (to see if we need to add people)
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS jobs_target_area ON %s (target_area_id, state);`, tableJobs),

		// events we want to search by tick (ie. what happened this tick?)
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS event_tick_type ON %s (tick, type);`, tableEvents),

		// families we mostly run over to see if we need to add children on a tick, so area + child_bearing
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS fam_child_bearing ON %s (area_id, is_child_bearing);`, tableFamilies),

		// hunt for people in given areas, by job, family, or their preferred_profession & area
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS peo_area ON %s (area_id);`, tablePeople),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS peo_prof ON %s (area_id, preferred_profession);`, tablePeople),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS peo_job ON %s (job_id);`, tablePeople),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS peo_fam ON %s (birth_family_id);`, tablePeople),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS peo_rand ON %s (random);`, tablePeople),

		// factions we really only look up either by ID, by action_frequency or by government
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS fact_government ON %s (government_id);`, tableFactions),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS fact_action_fq ON %s (action_frequency_ticks);`, tableFactions),
	}
)

// Sqlite represents a DB connection to sqlite
type Sqlite struct {
	*sqlDB
}

// NewSqlite3 opens a SQLite DB file.
// We also attempt to create / update tables to make it ready for use.
func NewSqlite3(cfg *config.Database) (*Sqlite, error) {
	if cfg.Driver != config.DatabaseSQLite3 {
		return nil, fmt.Errorf("invalid Engine, expected %s", config.DatabaseSQLite3)
	}

	db, err := sqlx.Connect("sqlite3", cfg.Location)
	if err != nil {
		return nil, err
	}
	me := &Sqlite{&sqlDB{conn: db}}
	return me, me.setupDatabase()
}

// setupDatabase tries to bring an SQLite db into line with our schema, upgrading as
// needed.
// We consider this acceptable in the case of local SQlite - but postgres is an entire
// 'nother kettle of fish.
func (s *Sqlite) setupDatabase() error {
	err := s.createTables()
	if err != nil {
		return err
	}

	txn, err := s.Transaction()
	if err != nil {
		return err
	}

	err = txn.SetMeta(metaSchemaVersion, "", currentSchemaVersion)
	if err != nil {
		txn.Rollback()
		return err
	}

	return txn.Commit()
}

// createTables & wrap up an errors we get.
// We'll try to press on despite errors.
func (s *Sqlite) createTables() error {
	var final error
	todo := []string{
		createMeta,
		createArea,
		createPlots,
		createGovernments,
		createLaws,
		createJobs,
		createFamilies,
		createPeople,
		createFactions,
		createEvents,
	}
	todoIx := []string{}

	for _, r := range allRelations {
		todo = append(todo, fmt.Sprintf(tupleTableCreateTemplate, r.tupleTable()))
		todoIx = append(todoIx, fmt.Sprintf(tupleIndexCreateTemplate, r.tupleTable(), r.tupleTable()))
		if r.SupportsModifiers() {
			todo = append(todo, fmt.Sprintf(modifierTableCreateTemplate, r.modTable()))
			todoIx = append(todoIx, fmt.Sprintf(tupleIndexCreateTemplate, r.modTable(), r.modTable()))
		}
	}

	todo = append(todo, indexes...)
	todo = append(todo, todoIx...)

	for _, ddl := range todo {
		_, err := s.conn.Exec(ddl)
		fmt.Println(ddl)
		if err != nil {
			if final == nil {
				final = err
			} else {
				final = fmt.Errorf("%w %v", final, err)
			}
		}
	}
	return final
}
