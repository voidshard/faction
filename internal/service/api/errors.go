package api

import (
	"errors"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	ErrInvalid  = errors.New("invalid")
	ErrNotFound = errors.New("not found")
)

func toError(err error, code ...structs.ErrorCode) *structs.Error {
	if err == nil {
		return nil
	}
	if len(code) > 0 {
		return &structs.Error{Code: code[0], Message: err.Error()}
	}
	if errors.Is(err, db.ErrNotFound) || errors.Is(err, ErrNotFound) {
		return &structs.Error{Code: structs.ErrorCode_NOT_FOUND, Message: err.Error()}
	}
	if errors.Is(err, db.ErrEtagMismatch) {
		return &structs.Error{Code: structs.ErrorCode_CONFLICT, Message: err.Error()}
	}
	if errors.Is(err, db.ErrInvalid) || errors.Is(err, ErrInvalid) {
		return &structs.Error{Code: structs.ErrorCode_INVALID, Message: err.Error()}
	}
	return &structs.Error{Code: structs.ErrorCode_INTERNAL, Message: err.Error()}
}
