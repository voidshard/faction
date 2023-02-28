package db

import (
	"github.com/voidshard/faction/pkg/structs"
)

// Database is something the simulation uses to record data.
// We don't expect users to supply their own implementations (if so, they can add one to internal/),
// we keep the interface here to outline what the Simulation needs.
type Database interface {
	Transaction() (Editor, error)
}

type Editor interface {
	// Generally;
	// - All Filter fields are considered "AND" and all Filter
	// objects (in a list of filters) are considered "OR"
	// Ie. Filter{id=a AND value = 3} OR Filter{id=b}
	//
	// - All 'Set' operations are Upserts ("insert or update").
	//
	// - Add get / search / list functions take a string 'token' and return
	// such a token.
	// > An empty string in a reply implies that there are no more results.
	// Although this token is fairly obvious (currently) this isn't meant to
	// be parsed by the end user(s), and it's design & encoding may change.

	Commit() error
	Rollback() error

	// Clock get, set
	Tick() (int, error)
	SetTick(i int) error

	Areas(token string, f ...*AreaFilter) ([]*structs.Area, string, error)
	SetAreas(in []*structs.Area) error

	Factions(token string, f ...*FactionFilter) ([]*structs.Faction, string, error)
	SetFactions(in []*structs.Faction) error

	Families(token string, f ...*FamilyFilter) ([]*structs.Family, string, error)
	SetFamilies(in []*structs.Family) error

	Governments(token string, ids ...string) ([]*structs.Government, string, error)
	SetGovernments(in []*structs.Government) error

	Jobs(token string, f ...*JobFilter) ([]*structs.Job, string, error)
	SetJobs(in []*Job) error

	LandRights(token string, f ...*LandRightFilter) ([]*structs.LandRight, string, error)
	SetLandRighs(in []*structs.LandRight) error

	People(token string, f ...*PersonFilter) ([]*structs.Person, string, error)
	SetPeople(in []*structs.Person) error

	Tuples(table Relation, token string, f ...*TupleFilter) ([]*structs.Tuple, string, error)
	SetTuples(table Relation, in ...*structs.Tuple) error
	ComputeTuples(table Relation, token string, tick int, f ...*TupleFilter) ([]*structs.Tuple, string, error)

	Modifiers(table Relation, token string, f ...*TupleFilter) ([]*structs.Modifier, string, error)
	SetModifiers(table Relation, in ...*structs.Modifier) error
	DeleteModifiers(table Relation, expires_before_tick int) error
}