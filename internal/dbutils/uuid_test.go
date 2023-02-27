package dbutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidID(t *testing.T) {
	cases := []struct {
		In    string
		Valid bool
	}{
		{">.<", false},
		{RandomID(), true},
		{"123e4567-e89b-12d3-a456-426655440000", true},
		{"123e4-12fdz3-a4ff56-4d26655440000", false},
		{NewID(1, 2, "aawdw"), true},
	}

	for _, tt := range cases {
		assert.Equal(t, IsValidID(tt.In), tt.Valid)
	}
}
