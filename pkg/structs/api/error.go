package api

type ErrorResponse struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}
