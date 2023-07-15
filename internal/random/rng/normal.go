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
	Int() int
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

// Int returns a random index, with the probability
// given by the normalised input on creation.
func (n *normalised) Int() int {
	if len(n.given) <= 0 {
		return -1
	}

	roll := n.rng.Float64()
	sofar := 0.0
	for k, v := range n.given {
		sofar += v / n.total
		if roll <= sofar {
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
