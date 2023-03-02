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

	IsMale bool `db:"is_male"`
}
