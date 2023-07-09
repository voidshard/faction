package base

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/voidshard/faction/pkg/config"
)

// Dice holds random dice for a demographics (race, culture) struct.
//
// We need quite a few with various average / deviation values so
// this helps keep things tidy.
//
// Ok well, tidy-er.
type Dice struct {
	simCfg *config.Simulation
	rng    *rand.Rand
	cfgs   map[string]*Demographic
}

func rckey(race, culture string) string {
	return fmt.Sprintf("%s:%s", race, culture)
}

func New(simCfg *config.Simulation) *Dice {
	d := &Dice{
		simCfg: simCfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
		cfgs:   map[string]*Demographic{},
	}

	for race, rdata := range simCfg.Races {
		for culture, cdata := range simCfg.Cultures {
			key := rckey(race, culture)
			d.cfgs[key] = newDemographic(race, culture, rdata, cdata)
		}
	}

	return d
}

// MaxDeathAdultMortalityProbability returns the highest DeathAdultMortalityProbability
// amoung all Demographics.
func (d *Dice) MaxDeathAdultMortalityProbability() float64 {
	// value is set on Culture, we don't need to iterate race information
	v := 0.0
	for _, culture := range d.simCfg.Cultures {
		if culture.DeathAdultMortalityProbability > v {
			v = culture.DeathAdultMortalityProbability
		}
	}
	return v
}

func (d *Dice) Float64() float64 {
	return d.rng.Float64()
}

func (d *Dice) Intn(n int) int {
	return d.rng.Intn(n)
}

func (d *Dice) IsValidDemographic(race, culture string) bool {
	_, ok := d.cfgs[rckey(race, culture)]
	return ok
}

func (d *Dice) MustDemographic(race, culture string) *Demographic {
	k := rckey(race, culture)
	demo := d.cfgs[k] // panic is intentional, the caller should check race/cultures are valid first
	return demo
}
