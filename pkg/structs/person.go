package structs

// Person roughly outlines someone that can belong to / work for a faction.
type Person struct {
	Ethos // rough outlook

	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`

	ID            string `db:"id"`
	BirthFamilyID string `db:"birth_family_id"` // family person was born into

	Race string `db:"race"`

	AreaID string `db:"area_id"` // area person lives in
	JobID  string `db:"job_id"`  // ie. current job id (if any)

	BirthTick int `db:"birth_tick"`
	DeathTick int `db:"death_tick"`

	IsMale  bool `db:"is_male"`
	IsChild bool `db:"is_child"` // person is too young for much of anything

	DeathMetaReason string  `db:"death_meta_reason"`
	DeathMetaKey    MetaKey `db:"death_meta_key"`
	DeathMetaVal    string  `db:"death_meta_val"`

	PreferredProfession string `db:"preferred_profession"` // ie. what they want to do for a living
	PreferredFactionID  string `db:"preferred_faction_id"` // ie. who they want to work for
}
