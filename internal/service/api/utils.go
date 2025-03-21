package api

import (
	"encoding/json"
	"net/http"
)

func readJson(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(v)
}
