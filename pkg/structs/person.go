package structs

// Person roughly outlines someone that can belong to / work for a faction.
type Person struct {
	Ethos // rough outlook

	ID string `db:"id"`

	AreaID string `db:"area_id"` // area person lives in
	JobID  string `db:"job_id"`  // ie. current job id (if any)

	BirthTick int `db:"birth_tick"`
	DeathTick int `db:"death_tick"`

	IsMale bool `db:"is_male"`
}
