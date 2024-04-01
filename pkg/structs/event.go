package structs

import (
	"encoding/json"
)

func AllEventTypes() []EventType {
	all := []EventType{}
	for _, r := range EventType_value {
		all = append(all, EventType(r))
	}
	return all
}

func (e *Event) MarshalJson() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) UnmarshalJson(b []byte) error {
	return json.Unmarshal(b, e)
}

func (e *Event) ObjectID() string {
	return e.ID
}
