package faction

type Database interface {
	Transaction() (*Editor, error)
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

	Areas(token string, f ...*AreaFilter) ([]*Area, string, error)
	SetAreas(in []*Area) error

	Factions(token string, f ...*FactionFilter) ([]*Faction, string, error)
	SetFactions(in []*Faction) error

	Governments(token string, ids ...string) ([]*Government, string, error)
	SetGovernments(in []*Government) error

	Jobs(token string, f ...*JobFilter) ([]*Job, string, error)
	SetJobs(in []*Job) error

	LandRights(token string, f ...*LandRightFilter) ([]*LandRight, string, error)
	SetLandRighs(in []*LandRight) error

	People(token string, f ...*PersonFilter) ([]*Person, string, error)
	SetPeople(in []*Person) error

	Tuples(table, token string, f ...*TupleFilter) ([]*Tuple, string, error)
	SetTuples(table string, in ...Tuple) error
}
