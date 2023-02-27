package dbutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidName(t *testing.T) {
	cases := []struct {
		In    string
		Valid bool
	}{
		{">.<", false},
		{RandomID(), true},
		{"123e4567-e89b-12d3-a456-426655440000", true},
		{"123e4-12fdz3-a4ff56-4d26655440000", true},
		{NewID(1, 2, "aawdw"), true},
		{"-8132-142-1241", false},
		{"8132-142-1241-", false},
		{"123e4567-e89b-12d3-a456-426655440000-heightmap", true},
	}

	for _, tt := range cases {
		assert.Equal(t, IsValidName(tt.In), tt.Valid)
	}
}
