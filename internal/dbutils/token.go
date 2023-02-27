package dbutils

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	limitDefault = 5000
	limitMax     = 10000
)

var (
	// ErrInvalidToken probably implies the token was garbled
	ErrInvalidToken = fmt.Errorf("iter token is invalid")
)

// IterToken is a simple struct representing a state
// in an iteration of a large set
type IterToken struct {
	Limit  int
	Offset int
}

// String turns our token into a string
func (i *IterToken) String() string {
	return fmt.Sprintf("%d,%d", i.Limit, i.Offset)
}

// NewIterToken returns a default starting token
func NewIterToken() *IterToken {
	return &IterToken{
		Limit:  5000,
		Offset: 0,
	}
}

// ParseIterToken returns our limit / offset from a token
func ParseIterToken(in string) (*IterToken, error) {
	if in == "" {
		// special case for those starting an iteration
		return NewIterToken(), nil
	}

	bits := strings.SplitN(in, ",", 2)

	limit, err := strconv.Atoi(bits[0])
	if err != nil {
		return nil, err
	}
	if limit < 1 {
		return nil, fmt.Errorf("%w found limit %d", ErrInvalidToken, limit)
	}
	if limit > limitMax {
		limit = limitMax
	}

	offset, err := strconv.Atoi(bits[1])
	if err != nil {
		return nil, err
	}
	if offset < 0 {
		return nil, fmt.Errorf("%w found offset %d", ErrInvalidToken, offset)
	}

	return &IterToken{
		Limit:  limit,
		Offset: offset,
	}, nil
}
