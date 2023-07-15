package structs

type EventType string

const (
	EventPersonDeath      EventType = "person_death"
	EventPersonBirth      EventType = "person_birth"
	EventPersonMove       EventType = "person_move"
	EventPersonChangeProf EventType = "person_change_profession"

	EventFamilyMarriage EventType = "family_marriage"
	EventFamilyDivorce  EventType = "family_divorce"
	EventFamilyPregnant EventType = "family_pregnant"
	EventFamilyAdoption EventType = "family_adoption"
	EventFamilyMove     EventType = "family_move"
)

// Event is something we want to report to the caller
type Event struct {
	ID   string    `db:"id"`
	Type EventType `db:"type"`
	Tick int       `db:"tick"`

	// Key is some user defined key in an Action config (see pkg/config/action.go)
	MetaKey MetaKey `db:"meta_key"`
	MetaVal string  `db:"meta_val"`

	// If we know the src event
	SourceEvent string `db:"source_event"`

	// Human readable message
	Message string `db:"message"`
}
