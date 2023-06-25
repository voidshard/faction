package db

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// Read only functions. We can read these outside of a transaction.
type Reader interface {
	// Generally;
	// - All Filter fields are considered "AND" and all Filter
	// objects (in a list of filters) are considered "OR"
	// Ie. Filter{id=a AND value = 3} OR Filter{id=b}
	//
	// - All 'Set' operations are Upserts ("insert or update").
	//
	// - Add get / search / list functions take a string 'token' and return
	// such a token.
	// > the first 'get' can pass an empty string as a token
	// > An empty string in a reply implies that there are no more results.
	// Although this token is fairly obvious (currently) this isn't meant to
	// be parsed by the end user(s), and it's design & encoding may change.

	Tick() (int, error)
	Areas(token string, q *Query) ([]*structs.Area, string, error)
	Factions(token string, q *Query) ([]*structs.Faction, string, error)
	Families(token string, q *Query) ([]*structs.Family, string, error)
	Governments(token string, q *Query) ([]*structs.Government, string, error)
	Jobs(token string, q *Query) ([]*structs.Job, string, error)
	People(token string, q *Query) ([]*structs.Person, string, error)
	Plots(token string, q *Query) ([]*structs.Plot, string, error)
	Tuples(table Relation, token string, q *Query) ([]*structs.Tuple, string, error)
	Routes(token string, q *Query) ([]*structs.Route, string, error)
	Meta(id string) (string, int, error)
	Modifiers(table Relation, token string, q *Query) ([]*structs.Modifier, string, error)
	ModifiersSum(table Relation, token string, q *Query) ([]*structs.Tuple, string, error)
}

type Writer interface {
	SetTick(i int) error
	SetAreas(in ...*structs.Area) error
	SetFactions(in ...*structs.Faction) error
	SetFamilies(in ...*structs.Family) error
	SetGovernments(in ...*structs.Government) error
	SetJobs(in ...*structs.Job) error
	SetPeople(in ...*structs.Person) error
	SetPlots(in ...*structs.Plot) error
	SetTuples(table Relation, in ...*structs.Tuple) error
	SetModifiers(table Relation, in ...*structs.Modifier) error
	DeleteModifiers(table Relation, expires_before_tick int) error
	SetRoutes(in ...*structs.Route) error
	SetMeta(id, str_val string, int_val int) error
	IncrTuples(table Relation, v int, q *Query) error
}

type ReaderWriter interface {
	Reader
	Writer
}

// Database is something the simulation uses to record data.
// We don't expect users to supply their own implementations (if so, they can add one to internal/),
type Database interface {
	Reader

	Transaction() (Transaction, error)
}

// Transaction grants all read functions & includes write functions.
type Transaction interface {
	ReaderWriter

	Commit() error
	Rollback() error
}

func New(cfg *config.Database) (*FactionDB, error) {
	var dbconn Database
	var err error

	switch cfg.Driver {
	case config.DatabaseSQLite3:
		dbconn, err = NewSqlite3(cfg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown database driver: %s", cfg.Driver)
	}

	return &FactionDB{dbconn}, nil
}
