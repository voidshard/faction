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
	tablePeople      = "people"
	tablePlots       = "plots"
	tableEvents      = "events"

	// metaClock is the current simulation tick, stored in the metadata table
	metaClock = "tick"

	// chunks (rows) we write a time
	chunksPerWrite = 1200
)

type Row interface {
	*structs.Plot | *structs.Person | *structs.Job | *structs.Government | *structs.Family | *structs.Faction | *structs.Area | *structs.Tuple | *structs.Modifier | *structs.Event
}

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
	if t < 1 {
		t = 1
	}
	return t, err
}

func (s *sqlDB) FactionChildren(id string, rs []structs.FactionRelation) ([]*structs.Faction, error) {
	return factionChildren(s.conn, id, rs)
}

func (s *sqlDB) FactionParents(id string, rs []structs.FactionRelation) ([]*structs.Faction, error) {
	return factionParents(s.conn, id, rs)
}

func (s *sqlDB) FactionLeadership(limit int, ids ...string) (map[string]*FactionLeadership, error) {
	return factionLeadership(s.conn, limit, ids...)
}

func (s *sqlDB) FactionPlots(limit int, ids ...string) (map[string][]*structs.Plot, error) {
	return factionPlots(s.conn, limit, ids...)
}

func (s *sqlDB) Plots(token string, in *Query) ([]*structs.Plot, string, error) {
	return plots(s.conn, token, in)
}

func (s *sqlDB) People(token string, in *Query) ([]*structs.Person, string, error) {
	return people(s.conn, token, in)
}

func (s *sqlDB) Events(token string, in *Query) ([]*structs.Event, string, error) {
	return events(s.conn, token, in)
}

func (s *sqlDB) Jobs(token string, in *Query) ([]*structs.Job, string, error) {
	return jobs(s.conn, token, in)
}

func (s *sqlDB) Governments(token string, in *Query) ([]*structs.Government, string, error) {
	return governments(s.conn, token, in)
}

func (s *sqlDB) Families(token string, in *Query) ([]*structs.Family, string, error) {
	return families(s.conn, token, in)
}

func (s *sqlDB) Factions(token string, in *Query) ([]*structs.Faction, string, error) {
	return factions(s.conn, token, in)
}

func (s *sqlDB) Areas(token string, in *Query) ([]*structs.Area, string, error) {
	return areas(s.conn, token, in)
}

func (s *sqlDB) Tuples(r Relation, token string, in *Query) ([]*structs.Tuple, string, error) {
	return tuples(s.conn, r, token, in)
}

func (s *sqlDB) CountTuples(r Relation, f *Query) (int, error) {
	return countTuples(s.conn, r, f)
}

func (s *sqlDB) Modifiers(r Relation, token string, in *Query) ([]*structs.Modifier, string, error) {
	return modifiers(s.conn, r, token, in)
}

func (s *sqlDB) ModifiersSum(table Relation, token string, f *Query) ([]*structs.Tuple, string, error) {
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

func (s *sqlTx) FactionChildren(id string, rs []structs.FactionRelation) ([]*structs.Faction, error) {
	return factionChildren(s.tx, id, rs)
}

func (s *sqlTx) FactionParents(id string, rs []structs.FactionRelation) ([]*structs.Faction, error) {
	return factionParents(s.tx, id, rs)
}

func (t *sqlTx) FactionLeadership(limit int, ids ...string) (map[string]*FactionLeadership, error) {
	return factionLeadership(t.tx, limit, ids...)
}

func (t *sqlTx) FactionPlots(limit int, ids ...string) (map[string][]*structs.Plot, error) {
	return factionPlots(t.tx, limit, ids...)
}

func (t *sqlTx) Tick() (int, error) {
	_, tick, err := t.Meta(metaClock)
	if tick < 1 {
		tick = 1
	}
	return tick, err
}

func (t *sqlTx) SetTick(tick int) error {
	if tick <= 1 {
		return nil
	}
	return t.SetMeta(metaClock, "", tick)
}

func (t *sqlTx) Plots(token string, in *Query) ([]*structs.Plot, string, error) {
	return plots(t.tx, token, in)
}

func (t *sqlTx) People(token string, in *Query) ([]*structs.Person, string, error) {
	return people(t.tx, token, in)
}

func (t *sqlTx) Events(token string, in *Query) ([]*structs.Event, string, error) {
	return events(t.tx, token, in)
}

func (t *sqlTx) Jobs(token string, in *Query) ([]*structs.Job, string, error) {
	return jobs(t.tx, token, in)
}

func (t *sqlTx) Governments(token string, in *Query) ([]*structs.Government, string, error) {
	return governments(t.tx, token, in)
}

func (t *sqlTx) Families(token string, in *Query) ([]*structs.Family, string, error) {
	return families(t.tx, token, in)
}

func (t *sqlTx) Factions(token string, in *Query) ([]*structs.Faction, string, error) {
	return factions(t.tx, token, in)
}

func (t *sqlTx) Areas(token string, in *Query) ([]*structs.Area, string, error) {
	return areas(t.tx, token, in)
}

func (t *sqlTx) Tuples(r Relation, token string, in *Query) ([]*structs.Tuple, string, error) {
	return tuples(t.tx, r, token, in)
}

func (t *sqlTx) CountTuples(r Relation, f *Query) (int, error) {
	return countTuples(t.tx, r, f)
}

func (t *sqlTx) Modifiers(r Relation, token string, in *Query) ([]*structs.Modifier, string, error) {
	return modifiers(t.tx, r, token, in)
}

func (t *sqlTx) ModifiersSum(table Relation, token string, f *Query) ([]*structs.Tuple, string, error) {
	return modifiersSum(t.tx, table, token, f)
}

// chunkWrite stops us from trying to write too many rows at once; triggering
// errors about too many SQL variables
func chunkWrite[R Row](fn func(sqlOperator, []R) error, tx sqlOperator, in []R) error {
	for i := 0; i < len(in); i += chunksPerWrite {
		j := i + chunksPerWrite
		if j > len(in) {
			return fn(tx, in[i:])
		}
		err := fn(tx, in[i:j])
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *sqlTx) SetPlots(in ...*structs.Plot) error {
	return chunkWrite(setPlots, t.tx, in)
}

func (t *sqlTx) SetPeople(in ...*structs.Person) error {
	return chunkWrite(setPeople, t.tx, in)
}

func (t *sqlTx) SetEvents(in ...*structs.Event) error {
	return chunkWrite(setEvents, t.tx, in)
}

func (t *sqlTx) SetJobs(in ...*structs.Job) error {
	return chunkWrite(setJobs, t.tx, in)
}

func (t *sqlTx) SetGovernments(in ...*structs.Government) error {
	return chunkWrite(setGovernments, t.tx, in)
}

func (t *sqlTx) SetFamilies(in ...*structs.Family) error {
	return chunkWrite(setFamilies, t.tx, in)
}

func (t *sqlTx) SetFactions(in ...*structs.Faction) error {
	return chunkWrite(setFactions, t.tx, in)
}

func (t *sqlTx) SetAreas(in ...*structs.Area) error {
	return chunkWrite(setAreas, t.tx, in)
}

func (t *sqlTx) SetTuples(r Relation, in ...*structs.Tuple) error {
	return setTuples(t.tx, r, in)
}

func (t *sqlTx) SetModifiers(r Relation, in ...*structs.Modifier) error {
	return setModifiers(t.tx, r, in)
}

func (t *sqlTx) IncrTuples(r Relation, v int, f *Query) error {
	return incrTuples(t.tx, r, v, f)
}

func (t *sqlTx) IncrModifiers(r Relation, v int, f *Query) error {
	return incrModifiers(t.tx, r, v, f)
}

func (t *sqlTx) DeleteModifiers(r Relation, expires_before_tick int) error {
	return deleteModifiers(t.tx, r, expires_before_tick)
}
