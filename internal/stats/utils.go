package stats

import (
	"math/rand"
)

// ChooseIndexes picks indexes in a random order from 0 to numIndexes-1.
// If choices is greater than numIndexes, all indexes are returned, though
// their order is random.
// If no indexes are given, an empty list is returned (since there are
// no choices).
func ChooseIndexes(numIndexes, choices int) []int {
	indexes := []int{}

	if numIndexes <= 0 {
		return indexes
	}

	for i := 0; i < numIndexes; i++ {
		indexes = append(indexes, i)
	}

	rand.Shuffle(len(indexes), func(i, j int) {
		indexes[i], indexes[j] = indexes[j], indexes[i]
	})

	if choices >= len(indexes) {
		return indexes
	}

	return indexes[:choices]
}
