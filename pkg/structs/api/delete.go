package api

type DeleteRequest struct {
	Ids   []string `json:"Id" validate:"required,min=1,max=5000,dive,valid_id"`
	World string   `json:"World" validate:"alphanum-if-non-global"`
}

func NewDeleteRequest() *DeleteRequest {
	return &DeleteRequest{
		Ids: []string{},
	}
}

type DeleteResponse struct {
	Error *ErrorResponse `json:"Error"`
}
