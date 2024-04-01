package main

import (
	"errors"
	"fmt"
	"time"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/voidshard/faction/pkg/config"
)

const (
	docMigrate = `Run the migration`
)

type optsMigrateGeneral struct {
	Source string `long:"source" env:"MIGRATIONS_SOURCE" description:"Source of the migrations" default:"file://migrations/prod"`
}

type cmdMigrateUp struct {
	optsMigrateGeneral
	optsGeneral
	optsDatabase

	Steps int `long:"steps" env:"MIGRATIONS_UP_STEPS" description:"Number of steps to migrate up (defaults to all)"`
}

type cmdMigrateDown struct {
	optsMigrateGeneral
	optsGeneral
	optsDatabase

	Steps int `long:"steps" env:"MIGRATIONS_DOWN_STEPS" description:"Number of steps to migrate down (defaults to all)"`
}

type cmdMigrateVersion struct {
	optsMigrateGeneral
	optsGeneral
	optsDatabase
}

type cmdMigrateForce struct {
	optsMigrateGeneral
	optsGeneral
	optsDatabase

	Version int `long:"version" env:"MIGRATIONS_FORCE_VERSION" description:"Force-set the current version of the database"`
}

type cmdMigrateWait struct {
	optsMigrateGeneral
	optsGeneral
	optsDatabase

	Version int           `long:"version" env:"MIGRATIONS_WAIT_VERSION" description:"Min version to wait for"`
	Timeout time.Duration `long:"timeout" env:"MIGRATIONS_WAIT_TIMEOUT" description:"Time to wait before erroring (optional)"`
}

type optsMigrate struct {
	Up cmdMigrateUp `command:"up" description:"Up version the database"`

	Down cmdMigrateDown `command:"down" description:"Down version the database"`

	Version cmdMigrateVersion `command:"version" description:"Get or force-set the current version of the database"`

	Wait cmdMigrateWait `command:"wait" description:"Wait (block) for the database to be at least the given version"`
}

func (c *cmdMigrateForce) Execute(args []string) error {
	m, db, err := buildMigrate(c.Source, &config.Database{Location: c.DatabaseURL})
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println("Forcing version to", c.Version)
	return m.Force(c.Version)
}

func (c *cmdMigrateUp) Execute(args []string) error {
	m, db, err := buildMigrate(c.Source, &config.Database{Location: c.DatabaseURL})
	if err != nil {
		return err
	}
	defer db.Close()

	if c.Steps < 0 { // force steps to be positive
		c.Steps = -1 * c.Steps
	}

	if c.Steps == 0 {
		fmt.Println("Migrating up to highest available version")
		return m.Up()
	} else {
		fmt.Println("Migrating up", c.Steps, "steps")
		return m.Steps(c.Steps)
	}
}

func (c *cmdMigrateDown) Execute(args []string) error {
	m, db, err := buildMigrate(c.Source, &config.Database{Location: c.DatabaseURL})
	if err != nil {
		return err
	}
	defer db.Close()

	if c.Steps > 0 { // force steps to be negative
		c.Steps = -1 * c.Steps
	}

	if c.Steps == 0 {
		fmt.Println("Migrating down to lowest available version")
		return m.Down()
	} else {
		fmt.Println("Migrating down", c.Steps, "steps")
		return m.Steps(c.Steps)
	}
}

func (c *cmdMigrateVersion) Execute(args []string) error {
	m, db, err := buildMigrate(c.Source, &config.Database{Location: c.DatabaseURL})
	if err != nil {
		return err
	}
	defer db.Close()

	v, err := getVersion(m)
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}

func (c *cmdMigrateWait) Execute(args []string) error {
	m, db, err := buildMigrate(c.Source, &config.Database{Location: c.DatabaseURL})
	if err != nil {
		return err
	}
	defer db.Close()

	lastV := -1
	end := time.Now().Add(c.Timeout)
	for {
		nowV, err := getVersion(m)
		if err != nil {
			return err
		}
		if nowV != lastV {
			fmt.Println("Database is at version", nowV, "desired:", c.Version)
		}
		if nowV >= c.Version {
			return nil
		}
		if c.Timeout > 0 && time.Now().After(end) {
			return errors.New(fmt.Sprintf("timeout waiting for version %d", c.Version))
		}

		lastV = nowV
		time.Sleep(1 * time.Second)
	}
}

func getVersion(m *migrate.Migrate) (int, error) {
	vuint, isDirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return -1, nil
	}
	if err != nil {
		return -1, err
	}
	if isDirty {
		fmt.Println("Warning: Database is dirty, manual intervention may be needed to force a version")
	}
	return int(vuint), nil
}

func buildMigrate(source string, opts *config.Database) (*migrate.Migrate, *sql.DB, error) {
	db, err := sql.Open("postgres", opts.Location)
	if err != nil {
		return nil, nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		defer db.Close()
		return nil, nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(source, "postgres", driver)
	if err != nil {
		defer db.Close()
		return nil, nil, err
	}

	return m, db, nil
}
