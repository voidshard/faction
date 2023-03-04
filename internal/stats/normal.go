package stats

import (
	"math/rand"
	"time"
)

// Normalised is an interface for a random number generation that returns
// an element index randomly but with a given probability spread.
//
// The input spread is normalised so you don't have to ensure they add up to 100.
// At least one item must be given for this to make sense.
type Normalised interface {
	Random() int
	SetSeed(seed int64)
}

// normalised implements the Normalised interface
type normalised struct {
	rng   *rand.Rand
	total float64
	given []float64
}

// SetSeed sets our internal RNG seed
func (n *normalised) SetSeed(seed int64) {
	n.rng = rand.New(rand.NewSource(seed))
}

// Random returns a random index, with the probability
// given by the normalised input on creation.
func (n *normalised) Random() int {
	roll := n.rng.Float64()
	for k, v := range n.given {
		if roll <= v/n.total {
			return k
		}
	}
	return len(n.given) - 1
}

// NewNormalised creates a new normalised random number generator
func NewNormalised(in []float64) Normalised {
	total := 0.0
	for _, v := range in {
		total += v
	}
	return &normalised{
		total: total,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
		given: in,
	}
}
