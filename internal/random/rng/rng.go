package rng

import (
	"math"
	"math/rand"
	"time"
)

// Rand yields random numbers that follow some desired distribution.
// This isn't super pretty, but we only need something to get us started .. we could use
// straight up random values but it would make our societies seem a bit too random.
//
// We don't keep all yielded values in memory, rather we calculate running average &
// std deviation values as we go.
type Rand struct {
	rng *rand.Rand

	min       float64
	max       float64
	mean      float64
	deviation float64
}

// Float64 returns a new random value between min & max such that we stay reasonably
// close to the desired std deviation.
func (r *Rand) Float64() float64 {
	return r.value()
}

// Int is sugar over Float64; it returns a new random value between min & max such
// that we stay reasonably close to the desired std deviation.
func (r *Rand) Int() int {
	return int(r.value())
}

// value returns a new random value
func (r *Rand) value() float64 {
	if r.max <= r.min {
		return r.min
	}
	// go-recipes.dev/generating-random-numbers-with-go-616d30ccc926?gi=471026f18bf6
	// Can't believe I didn't see this NormFloat64() func before :smh:
	return math.Max(math.Min(r.rng.NormFloat64()*r.deviation+r.mean, r.max), r.min)
}

// SetSeed sets our internal RNG seed
func (r *Rand) SetSeed(seed int64) {
	r.rng = rand.New(rand.NewSource(seed))
}

// NewRand creates a new random number generator that attempts to yield
// values with some desired average (mean) and std deviation.
func NewRand(min, max, mean, deviation float64) *Rand {
	return &Rand{
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
		min:       min,
		max:       max,
		mean:      mean,
		deviation: deviation,
	}
}
