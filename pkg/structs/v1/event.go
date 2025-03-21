package v1

type Event struct {
	World      string `json:"world"`
	Kind       string `json:"kind"`
	Controller string `json:"controller"`
	Id         string `json:"id"`

	AckId string `json:"ack_id,omitempty"`
}
