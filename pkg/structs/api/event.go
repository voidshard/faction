package api

type StreamEvents struct {
	// Filters, where not set we mean "all" (match any)
	World      string `json:"World" validate:"alphanum-or-empty"`
	Kind       string `json:"Kind" validate:"alphanum-or-empty"`
	Controller string `json:"Controller" validate:"alphanum-or-empty"`
	Id         string `json:"Id" validate:"alphanum-or-empty"`

	// Queue name to listen on, if set implies durable subscription
	Queue string `json:"Queue" validate:"alphanum-or-empty"`
}

type DeferEventRequest struct {
	World      string `json:"World" validate:"alphanum"`
	Kind       string `json:"Kind" validate:"alphanum"`
	Controller string `json:"Controller" validate:"alphanum-or-empty"`
	Id         string `json:"Id" validate:"alphanum"`

	// technically this allows both of these to be set, if so we honor ToTick
	// first, as it is more specific and requires less computation.
	ToTick uint64 `json:"ToTick" validate:"gte=0,required_without=ByTick"`
	ByTick uint64 `json:"ByTick" validate:"gte=0,required_without=ToTick"`
}

type DeferEventResponse struct {
	ToTick uint64         `json:"ToTick"`
	Error  *ErrorResponse `json:"Error"`
}
