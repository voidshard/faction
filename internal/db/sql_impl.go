package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/voidshard/faction/pkg/structs"
)

const (
	// SQL table names
	tableMeta        = "meta"
	tableAreas       = "areas"
	tableFactions    = "factions"
	tableFamilies    = "families"
	tableGovernments = "governments"
	tableLaws        = "laws"
	tableJobs        = "jobs"
	tableLandRights  = "landrights"
	tablePeople      = "people"
	tableRoutes      = "routes"
	tablePlots       = "plots"

	// metaClock is the current simulation tick, stored in the metadata table
	metaClock = "tick"
)

// sqlDB represents a generic DB wrapper -- this allows SQLite & Postgres to run
// the same code & execute queries in the same manner.
//
// We do have to be careful as the two underlying DBs are not exactly the same,
// but for our simple query requirements it's fine.
//
// We also divide the structs into
// - read only (can be outside of a transaction)
// - read / write (must be inside a transaction)
type sqlDB struct {
	conn *sqlx.DB
}

// sqlTx represents a transaction with read/write powers
type sqlTx struct {
	tx *sqlx.Tx
}

// Meta returns saved metadata, outside of a transaction
func (s *sqlDB) Meta(id string) (string, int, error) {
	return meta(s.conn, id)
}

func (s *sqlDB) Tick() (int, error) {
	_, t, err := s.Meta(metaClock)
	return t, err
}

func (s *sqlDB) Plots(token string, in ...*PlotFilter) ([]*structs.Plot, string, error) {
	return plots(s.conn, token, in)
}

func (s *sqlDB) Routes(token string, in ...*RouteFilter) ([]*structs.Route, string, error) {
	return routes(s.conn, token, in)
}

func (s *sqlDB) People(token string, in ...*PersonFilter) ([]*structs.Person, string, error) {
	return people(s.conn, token, in)
}

func (s *sqlDB) LandRights(token string, in ...*LandRightFilter) ([]*structs.LandRight, string, error) {
	return landRights(s.conn, token, in)
}

func (s *sqlDB) Jobs(token string, in ...*JobFilter) ([]*structs.Job, string, error) {
	return jobs(s.conn, token, in)
}

func (s *sqlDB) Governments(token string, in ...*GovernmentFilter) ([]*structs.Government, string, error) {
	return governments(s.conn, token, in)
}

func (s *sqlDB) Families(token string, in ...*FamilyFilter) ([]*structs.Family, string, error) {
	return families(s.conn, token, in)
}

func (s *sqlDB) Factions(token string, in ...*FactionFilter) ([]*structs.Faction, string, error) {
	return factions(s.conn, token, in)
}

func (s *sqlDB) Areas(token string, in ...*AreaFilter) ([]*structs.Area, string, error) {
	return areas(s.conn, token, in)
}

func (s *sqlDB) Tuples(r Relation, token string, in ...*TupleFilter) ([]*structs.Tuple, string, error) {
	return tuples(s.conn, r, token, in)
}

func (s *sqlDB) Modifiers(r Relation, token string, in ...*ModifierFilter) ([]*structs.Modifier, string, error) {
	return modifiers(s.conn, r, token, in)
}

func (s *sqlDB) ModifiersSum(table Relation, token string, f ...*ModifierFilter) ([]*structs.Tuple, string, error) {
	return modifiersSum(s.conn, table, token, f)
}

// Close connection to DB
func (s *sqlDB) Close() error {
	return s.conn.Close()
}

// Begin a new transaction.
// You *must* either call Commit() or Rollback() after calling this
func (s *sqlDB) Transaction() (Transaction, error) {
	tx, err := s.conn.Beginx()

	if err != nil {
		return nil, err
	}

	return &sqlTx{tx: tx}, nil
}

// Commit finish & commit changes to the database
func (t *sqlTx) Commit() error {
	return t.tx.Commit()
}

// Rollback abort transaction without changing anything
func (t *sqlTx) Rollback() error {
	return t.tx.Rollback()
}

// Meta returns saved metadata
func (t *sqlTx) Meta(id string) (string, int, error) {
	return meta(t.tx, id)
}

// SetMeta sets some metadata within a transaction
func (t *sqlTx) SetMeta(id, strv string, intv int) error {
	return setMeta(t.tx, id, strv, intv)
}

func (t *sqlTx) Tick() (int, error) {
	_, tick, err := t.Meta(metaClock)
	return tick + 1, err
}

func (t *sqlTx) SetTick(tick int) error {
	if tick <= 1 {
		return nil
	}
	return t.SetMeta(metaClock, "", tick-1)
}

func (t *sqlTx) Plots(token string, in ...*PlotFilter) ([]*structs.Plot, string, error) {
	return plots(t.tx, token, in)
}

func (t *sqlTx) Routes(token string, in ...*RouteFilter) ([]*structs.Route, string, error) {
	return routes(t.tx, token, in)
}

func (t *sqlTx) People(token string, in ...*PersonFilter) ([]*structs.Person, string, error) {
	return people(t.tx, token, in)
}

func (t *sqlTx) LandRights(token string, in ...*LandRightFilter) ([]*structs.LandRight, string, error) {
	return landRights(t.tx, token, in)
}

func (t *sqlTx) Jobs(token string, in ...*JobFilter) ([]*structs.Job, string, error) {
	return jobs(t.tx, token, in)
}

func (t *sqlTx) Governments(token string, in ...*GovernmentFilter) ([]*structs.Government, string, error) {
	return governments(t.tx, token, in)
}

func (t *sqlTx) Families(token string, in ...*FamilyFilter) ([]*structs.Family, string, error) {
	return families(t.tx, token, in)
}

func (t *sqlTx) Factions(token string, in ...*FactionFilter) ([]*structs.Faction, string, error) {
	return factions(t.tx, token, in)
}

func (t *sqlTx) Areas(token string, in ...*AreaFilter) ([]*structs.Area, string, error) {
	return areas(t.tx, token, in)
}

func (t *sqlTx) Tuples(r Relation, token string, in ...*TupleFilter) ([]*structs.Tuple, string, error) {
	return tuples(t.tx, r, token, in)
}

func (t *sqlTx) Modifiers(r Relation, token string, in ...*ModifierFilter) ([]*structs.Modifier, string, error) {
	return modifiers(t.tx, r, token, in)
}

func (t *sqlTx) ModifiersSum(table Relation, token string, f ...*ModifierFilter) ([]*structs.Tuple, string, error) {
	return modifiersSum(t.tx, table, token, f)
}

func (t *sqlTx) SetPlots(plots ...*structs.Plot) error {
	return setPlots(t.tx, plots)
}

func (t *sqlTx) SetRoutes(routes ...*structs.Route) error {
	return setRoutes(t.tx, routes)
}

func (t *sqlTx) SetPeople(people ...*structs.Person) error {
	return setPeople(t.tx, people)
}

func (t *sqlTx) SetLandRights(landRights ...*structs.LandRight) error {
	return setLandRights(t.tx, landRights)
}

func (t *sqlTx) SetJobs(jobs ...*structs.Job) error {
	return setJobs(t.tx, jobs)
}

func (t *sqlTx) SetGovernments(governments ...*structs.Government) error {
	return setGovernments(t.tx, governments)
}

func (t *sqlTx) SetFamilies(families ...*structs.Family) error {
	return setFamilies(t.tx, families)
}

func (t *sqlTx) SetFactions(factions ...*structs.Faction) error {
	return setFactions(t.tx, factions)
}

func (t *sqlTx) SetAreas(areas ...*structs.Area) error {
	return setAreas(t.tx, areas)
}

func (t *sqlTx) SetTuples(r Relation, tuples ...*structs.Tuple) error {
	return setTuples(t.tx, r, tuples)
}

func (t *sqlTx) SetModifiers(r Relation, mod ...*structs.Modifier) error {
	return setModifiers(t.tx, r, mod)
}

func (t *sqlTx) IncrTuples(r Relation, v int, f ...*TupleFilter) error {
	return incrTuples(t.tx, r, v, f)
}

func (t *sqlTx) IncrModifiers(r Relation, v int, f ...*ModifierFilter) error {
	return incrModifiers(t.tx, r, v, f)
}

func (t *sqlTx) DeleteModifiers(r Relation, expires_before_tick int) error {
	return deleteModifiers(t.tx, r, expires_before_tick)
}
