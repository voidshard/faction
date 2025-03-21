package api

type GetRequest struct {
	Ids    []string          `json:"Id" validate:"min=0,max=5000,dive,valid_id"`
	Limit  int64             `json:"Limit" validate:"gte=0,lte=5000"`
	Offset int64             `json:"Offset" validate:"gte=0"`
	Labels map[string]string `json:"Labels" validate:"max=10`
	World  string            `json:"World" validate:"alphanum-if-non-global"`
}

func NewGetRequest() *GetRequest {
	return &GetRequest{
		Ids:    []string{},
		Limit:  100,
		Offset: 0,
		Labels: map[string]string{},
	}
}

type GetResponse struct {
	Data  []interface{}  `json:"Data"`
	Error *ErrorResponse `json:"Error"`
}
