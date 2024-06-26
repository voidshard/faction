package config

import (
	"os"
	"path/filepath"
	"strings"
)

// DatabaseDriver denotes what kind of data store we'll use
type DatabaseDriver string

const (
	// DatabaseSQLite3 uses sqlite3 to store data.
	// Easy to manage, good for testing, not recommended for large datasets / large simulations.
	DatabaseSQLite3 DatabaseDriver = "sqlite3"

	// DatabasePostgres uses postgres to store data.
	// Probably you want to use this for anything serious -- this is also required for anything involving the
	// queue.
	DatabasePostgres DatabaseDriver = "postgres"
)

// Database configuration struct
type Database struct {
	// Driver denotes the underlying db implementation (see above)
	Driver DatabaseDriver

	// Location of database; where to find
	// [sqlite3]: file path
	// [postgres]: connection string
	Location string
}

func DefaultDatabase() *Database {
	if isLocalMode() {
		// This is useful for debugging without needing to edit code.
		return &Database{
			Driver:   DatabaseSQLite3,
			Location: filepath.Join(os.TempDir(), "faction.sqlite"),
		}
	}
	return &Database{
		Driver:   DatabasePostgres,
		Location: "postgres://postgres:test@localhost:5432/faction?sslmode=disable",
	}
}

func isLocalMode() bool {
	localmode := strings.ToLower(os.Getenv("ENABLE_LOCAL_MODE"))
	return localmode == "true" || localmode == "1" || localmode == "yes" || localmode == "on"
}
