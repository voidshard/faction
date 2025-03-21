package api

import (
	v1 "github.com/voidshard/faction/pkg/structs/v1"
)

type SearchRequest struct {
	v1.Query
}

type SearchResponse struct {
	Data  []interface{}  `json:"Data"`
	Error *ErrorResponse `json:"Error"`
}
