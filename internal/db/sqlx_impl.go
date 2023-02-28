package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Relation here is a name we accept for Tuple(s) and Modifier(s).
type Relation string

const (
	// Relation Table Partitioning
	//
	// We could in theory have all Tuples and Modifiers in two giant tables
	// (possible even the same table where Tuples don't use all of the columns)
	// and partition by the Relation (in the case of Postgres at least?) but
	// in practice we disallow querying the table directly in our interface and
	// it's sort of easier conceptually / book keeping wise to have distinct tables
	// which hold distinct data.
	//
	// Ie. an inter-personal relationship (wife / son / ally) between Person X and
	// Person Y has no bearing on the trust between Faction A and Faction B .. and
	// we don't to be able to query them in the same breath .. although they're both
	// tuples and in theory we *could*
	//
	// Thus currently each Relation has a "tuples_" and a "modifiers_" table created
	// for it. Modifiers apply temporary alterations to Tuple values until their expire
	// time is reached.
	//
	// Then given:
	// Tuple from tuples_NAME
	//    Subject (A), Object (B), Value (V)
	// Modifier(s) from modifiers_NAME
	//    Subject (A), Object (B), Value (X), Expires 10
	//    Subject (A), Object (B), Value (Y), Expires 12
	//    Subject (A), Object (B), Value (Z), Expires 15
	// We compute the final value based on tick T as
	//    [Time]       =     [Final Value]
	//    T < 10       =     V + X + Y + Z
	//    10 <= T < 12 =     V + Y + Z
	//    12 <= T < 15 =     V + Z
	//    15 <= T      =     V

	// RelationPFAffiliation holds how closely people affiliate with factions
	RelationPFAffiliation Relation = "affiliation_person_to_faction"

	// RelationFFTrust holds how much factions trust each other
	RelationFFTrust Relation = "trust_faction_to_faction"

	// RelationPProfessionSkill holds how skilled people are in professions
	RelationPProfessionSkill Relation = "skill_person_to_profession"
)

const (
	// SQL table names
	tableMeta        = "meta"
	tableAreas       = "areas"
	tableFaction     = "factions"
	tableFamilies    = "families"
	tableGovernments = "governments"
	tableJobs        = "jobs"
	tableLandRights  = "landrights"
	tablePeople      = "people"

	// metaClock is the current simulation tick, stored in the metadata table
	metaClock = "tick"
)

func (r Relation) tupleTable() string { return fmt.Sprintf("tuples_%s", r) }

func (r Relation) modTable() string { return fmt.Sprintf("modifiers_%s", r) }

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
