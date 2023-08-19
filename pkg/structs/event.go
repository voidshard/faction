package structs

import (
	"encoding/json"
)

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

	EventFactionPromotion EventType = "faction_promotion"
)

// Event is something we want to report to the caller
type Event struct {
	ID   string    `db:"id"`
	Type EventType `db:"type"`
	Tick int       `db:"tick"`

	// what the event is talking about
	SubjectMetaKey MetaKey `db:"subject_meta_key"`
	SubjectMetaVal string  `db:"subject_meta_val"`

	// if we know the cause
	CauseMetaKey MetaKey `db:"cause_meta_key"`
	CauseMetaVal string  `db:"cause_meta_val"`

	// Human readable message
	Message string `db:"message"`
}

func (e *Event) MarshalJson() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) UnmarshalJson(b []byte) error {
	return json.Unmarshal(b, e)
}
