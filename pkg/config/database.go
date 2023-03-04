package config

// DatabaseDriver denotes what kind of data store we'll use
type DatabaseDriver string

const (
	// DatabaseSQLite3 uses sqlite3 to store data.
	// Easy to manage, good for testing, not recommended for large datasets / large simulations.
	DatabaseSQLite3 DatabaseDriver = "sqlite3"
)

// Database configuration struct
type Database struct {
	// Driver denotes the underlying db implementation
	Driver DatabaseDriver

	// Name of database
	// [sqlite3]: file name
	// [postgres]: database name
	Name string

	// Location of database; where to find
	// [sqlite3]: folder path
	// [postgres]: connection string
	Location string
}

func DefaultDatabase() *Database {
	return &Database{
		Driver: DatabaseSQLite3,
		Name:   "faction-sim.sqlite",
	}
}
