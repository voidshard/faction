package api

type SetRequest struct {
	Data  []interface{} `json:"Data" validate:"required,min=1,max=5000"`
	World string        `json:"World" validate:"alphanum-if-non-global"`
}

func NewSetRequest() *SetRequest {
	return &SetRequest{
		Data: []interface{}{},
	}
}

type SetResponse struct {
	Error *ErrorResponse `json:"Error"`
}
