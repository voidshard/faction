package db

import (
	"fmt"
	"path/filepath"

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
)

var (
	// NB. UUIDs are 36 chars (eg. 123e4567-e89b-12d3-a456-426655440000)

	// Sqlite3 we use as a DB for users wanting to run in local-sim mode.
	createMeta = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id VARCHAR(255) PRIMARY KEY,
	str TEXT NOT NULL DEFAULT "",
	int INTEGER NOT NULL DEFAULT 0
    );`, tableMeta)

	createArea = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	    id VARCHAR(36) PRIMARY KEY,
	    governing_faction_id VARCHAR(36)
	);`, tableAreas)

	// table for structs.Faction object
	createFactions = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
            id VARCHAR(36) PRIMARY KEY,
            ethos_altruism INTEGER NOT NULL DEFAULT 0,
            ethos_ambition INTEGER NOT NULL DEFAULT 0,
            ethos_tradition INTEGER NOT NULL DEFAULT 0,
            ethos_pacifism INTEGER NOT NULL DEFAULT 0,
            ethos_piety INTEGER NOT NULL DEFAULT 0,
            ethos_caution INTEGER NOT NULL DEFAULT 0,
            action_frequency_ticks INTEGER NOT NULL DEFAULT 1,
            leadership INTEGER NOT NULL DEFAULT 0,
            wealth INTEGER NOT NULL DEFAULT 0,
            cohesion INTEGER NOT NULL DEFAULT 0,
            corruption INTEGER NOT NULL DEFAULT 0,
            is_covered BOOLEAN NOT NULL DEFAULT FALSE,
            is_illegal BOOLEAN NOT NULL DEFAULT FALSE,
            government_id VARCHAR(36) NOT NULL,
            is_government BOOLEAN NOT NULL DEFAULT FALSE,
            religion_id VARCHAR(36),
            is_religion BOOLEAN NOT NULL DEFAULT FALSE,
            is_member_by_birth BOOLEAN NOT NULL DEFAULT FALSE,
            espionage_offense INTEGER NOT NULL DEFAULT 0,
            espionage_defense INTEGER NOT NULL DEFAULT 0,
            military_offense INTEGER NOT NULL DEFAULT 0,
            military_defense INTEGER NOT NULL DEFAULT 0,
            parent_faction_id VARCHAR(36),
            parent_faction_relation INTEGER,
            research_science INTEGER NOT NULL DEFAULT 0,
            research_theology INTEGER NOT NULL DEFAULT 0,
            research_magic INTEGER NOT NULL DEFAULT 0,
            research_occult INTEGER NOT NULL DEFAULT 0
        );`, tableFactions)

	// indexes that we should create
	indexes = []string{}
)

// Sqlite represents a DB connection to sqlite
type Sqlite struct {
	*sqlDB
}

// NewSqlite3 opens a SQLite DB file.
// We also attempt to create / update tables to make it ready for use.
func NewSqlite3(cfg config.Database) (*Sqlite, error) {
	if cfg.Engine != config.DatabaseEngineSQLite3 {
		return nil, fmt.Errorf("invalid Engine, expected %s", config.DatabaseEngineSQLite3)
	}

	db, err := sqlx.Connect("sqlite3", filepath.Join(cfg.Location, cfg.Name))
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
	todo := []string{createMeta}
	todo = append(todo, indexes...)
	for _, ddl := range todo {
		_, err := s.conn.Exec(ddl)
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
