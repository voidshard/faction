package stats

import (
	"math"
	"math/rand"
	"time"
)

const (
	// candidateValues is the number of random values we'll make before picking
	// the one that seems most desirable to keeping to(wards) our deviation.
	candidateValues = 5
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
	deviation float64

	// https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance
	// translated from python snippet on "Computing shifted data"
	k     float64 // mean
	count float64 // number of values generated
	ex    float64
	ex2   float64
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
//
// Nb. this is best effort as we probably don't sit exactly on this value, we don't
// guarantee that the std deviation requested is honored exactly .. I mean .. we *are*
// returning random values here .. we just aim to supply slightly less randomness so
// the values huddle around some variance.
func (r *Rand) value() float64 {
	// The brute force approach; we'll choose some number of values, experimentally
	// add them to our running values & see how the std dev. changes.
	// Then we pick the best value that keeps us closest (in absolute terms) to our
	// desired std deviation.
	// There's probably a smarter way to do this, I had a short look for a nice library
	// but .. nothing jumped out at me
	// TODO: more math less YOLO
	bestVal := r.min + r.rng.Float64()*(r.max-r.min)
	r.add(bestVal)
	bestDev := math.Abs(r.deviation - r.runningStdDev())
	r.sub(bestVal)

	for i := 0; i < candidateValues-1; i++ {
		v := r.min + r.rng.Float64()*(r.max-r.min)
		r.add(bestVal)
		d := math.Abs(r.deviation - r.runningStdDev())
		r.sub(bestVal)

		if d < bestDev {
			bestDev = d
			bestVal = v
		}
	}

	r.add(bestVal)
	return bestVal
}

// Add a value to the running totals
func (r *Rand) add(v float64) {
	r.count += 1
	r.ex += v - r.k
	r.ex2 += math.Pow(v-r.k, 2)
}

// Sub (remove) a value from the running totals
func (r *Rand) sub(v float64) {
	if r.count == 0 {
		return
	}
	r.count -= 1
	r.ex -= v - r.k
	r.ex2 -= math.Pow(v-r.k, 2)
}

// runningMean returns the running mean value (average)
func (r *Rand) runningMean() float64 {
	return r.k + r.ex/r.count
}

// runningStdDev returns the estimated std deviation squared.
//
// We could sqrt this but we're going to use this a lot internally - it's cheaper to square
// the desired std deviation (passed in on creation) than the sqrt everything whenever we
// call this.
func (r *Rand) runningStdDev() float64 {
	return (r.ex2 - math.Pow(r.ex, 2)/r.count) / (r.count - 1)
}

// SetSeed sets our internal RNG seed
func (r *Rand) SetSeed(seed int64) {
	r.rng = rand.New(rand.NewSource(seed))
}

// NewRand creates a new random number generator that attempts to yield
// values with some desired average (mean) and std deviation.
func NewRand(min, max, average, deviation float64) *Rand {
	return &Rand{
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
		min:       min,
		max:       max,
		deviation: math.Pow(deviation, 2),
		k:         average,
	}
}
