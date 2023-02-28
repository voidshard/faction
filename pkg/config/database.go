package config

// DatabaseEngine denotes what kind of data store we'll use
type DatabaseEngine string

const (
	// DatabaseEngineSQLite3 uses sqlite3 to store data.
	// Easy to manage, good for testing, not recommended for large datasets / large simulations.
	DatabaseEngineSQLite3 DatabaseEngine = "sqlite3"
)

// Database configuration struct
type Database struct {
	// Engine denotes the underlying db implementation
	Engine DatabaseEngine

	// Name of database
	// [sqlite3]: file name
	// [postgres]: database name
	Name string

	// Location of database; where to find
	// [sqlite3]: folder path
	// [postgres]: connection string
	Location string
}
