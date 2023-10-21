package structs

const (
	PersonRandomMax = 1000000
)

// Person roughly outlines someone that can belong to / work for a faction.
type Person struct {
	Ethos // rough outlook

	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`

	ID            string `db:"id"`
	BirthFamilyID string `db:"birth_family_id"` // family person was born into

	Race    string `db:"race"`
	Culture string `db:"culture"`

	AreaID string `db:"area_id"` // area person lives in
	JobID  string `db:"job_id"`  // ie. current job id (if any)

	BirthTick int `db:"birth_tick"`
	DeathTick int `db:"death_tick"`

	IsMale bool `db:"is_male"`

	AdulthoodTick int `db:"adulthood_tick"` // tick person becomes an adult

	DeathMetaReason string  `db:"death_meta_reason"`
	DeathMetaKey    MetaKey `db:"death_meta_key"`
	DeathMetaVal    string  `db:"death_meta_val"`

	PreferredProfession string `db:"preferred_profession"` // ie. what they want to do for a living
	PreferredFactionID  string `db:"preferred_faction_id"` // ie. who they want to work for

	// used for blind randomisation without needing to read a record from the DB.
	// ie. if a disease kills 1 in 10000 people we can just roll a random
	// number and set the death cause with a single update without needing to read any records.
	//
	// A number 0 -> 1000000
	// This number never changes.
	Random int `db:"random"`

	// NaturalDeathTick is the tick someone is slated to die of old age.
	//
	// Set on birth. (Let's call it .. "fated death"?)
	//
	// Computationally it's easier to store the tick someone is slated to die on
	// of "old age" than to randomly roll for each person down the line each tick to
	// see if they die.
	NaturalDeathTick int `db:"natural_death_tick"`
}

func (p *Person) ObjectID() string {
	return p.ID
}

// SetBirthTick updates the birthtick, and moves the natural death tick and adulthood tick
// accordingly.
func (p *Person) SetBirthTick(t int) {
	lifespan := p.NaturalDeathTick - p.BirthTick
	childhood := p.AdulthoodTick - p.BirthTick

	p.BirthTick = t
	p.NaturalDeathTick = t + lifespan
	p.AdulthoodTick = t + childhood
}
