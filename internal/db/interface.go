package db

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

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
	Areas(token string, f ...*AreaFilter) ([]*structs.Area, string, error)
	Factions(token string, f ...*FactionFilter) ([]*structs.Faction, string, error)
	Families(token string, f ...*FamilyFilter) ([]*structs.Family, string, error)
	Governments(token string, f ...*GovernmentFilter) ([]*structs.Government, string, error)
	Jobs(token string, f ...*JobFilter) ([]*structs.Job, string, error)
	LandRights(token string, f ...*LandRightFilter) ([]*structs.LandRight, string, error)
	People(token string, f ...*PersonFilter) ([]*structs.Person, string, error)
	Plots(token string, f ...*PlotFilter) ([]*structs.Plot, string, error)
	Tuples(table Relation, token string, f ...*TupleFilter) ([]*structs.Tuple, string, error)
	Routes(token string, f ...*RouteFilter) ([]*structs.Route, string, error)
	Meta(id string) (string, int, error)
	Modifiers(table Relation, token string, f ...*ModifierFilter) ([]*structs.Modifier, string, error)
	ModifiersSum(table Relation, token string, f ...*ModifierFilter) ([]*structs.Tuple, string, error)
}

// Database is something the simulation uses to record data.
// We don't expect users to supply their own implementations (if so, they can add one to internal/),
// we keep the interface here to outline what the Simulation needs.
type Database interface {
	Reader
	Transaction() (Transaction, error)
}

type Transaction interface {
	Reader

	Commit() error
	Rollback() error

	SetTick(i int) error

	SetAreas(in ...*structs.Area) error

	SetFactions(in ...*structs.Faction) error

	SetFamilies(in ...*structs.Family) error

	SetGovernments(in ...*structs.Government) error

	SetJobs(in ...*structs.Job) error

	SetLandRights(in ...*structs.LandRight) error

	SetPeople(in ...*structs.Person) error

	SetPlots(in ...*structs.Plot) error

	SetTuples(table Relation, in ...*structs.Tuple) error

	SetModifiers(table Relation, in ...*structs.Modifier) error
	DeleteModifiers(table Relation, expires_before_tick int) error

	SetRoutes(in ...*structs.Route) error

	SetMeta(id, str_val string, int_val int) error
}

func New(cfg *config.Database) (Database, error) {
	switch cfg.Driver {
	case config.DatabaseSQLite3:
		return NewSqlite3(cfg)
	default:
		return nil, fmt.Errorf("unknown database driver: %s", cfg.Driver)
	}
}
