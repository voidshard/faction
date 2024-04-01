package base

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
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

	ethicsByProfession map[string]*structs.Ethos
}

func rckey(race, culture string) string {
	return fmt.Sprintf("%s:%s", race, culture)
}

func New(simCfg *config.Simulation) *Dice {
	d := &Dice{
		simCfg:             simCfg,
		rng:                rand.New(rand.NewSource(time.Now().UnixNano())),
		cfgs:               map[string]*Demographic{},
		ethicsByProfession: calcEthosWeightsForProfessions(simCfg),
	}

	for race, rdata := range simCfg.Races {
		for culture, cdata := range simCfg.Cultures {
			key := rckey(race, culture)
			d.cfgs[key] = newDemographic(race, culture, rdata, cdata)
		}
	}

	return d
}

func calcEthosWeightsForProfessions(cfg *config.Simulation) map[string]*structs.Ethos {
	w := map[string][]*structs.Ethos{}

	ethicw := float64(structs.MaxEthos / 100)

	for _, act := range cfg.Actions {
		if act.ProfessionWeights == nil {
			continue
		}
		for prof, weight := range act.ProfessionWeights {
			cur, ok := w[prof]
			if !ok {
				cur = []*structs.Ethos{}
			}
			cur = append(cur, act.Ethos.Add(int64(weight*ethicw)))
			w[prof] = cur
		}
	}

	averages := map[string]*structs.Ethos{}

	for prof, ethics := range w {
		if len(ethics) == 0 {
			continue
		}

		// apply a 50% weight to the average so we can't have someone's ethos
		// *entirely* defined by their profession(s)
		averages[prof] = structs.EthosAverage(ethics...).Multiply(0.5).Clamp()
	}

	return averages
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

// EthosWeightFromProfession returns the ethos weight for a person based on their professions.
func (d *Dice) EthosWeightFromProfessions(prof []*structs.Tuple) *structs.Ethos {
	if len(prof) == 0 {
		return &structs.Ethos{}
	}
	eth := []*structs.Ethos{}
	for _, p := range prof {
		w, ok := d.ethicsByProfession[p.Subject]
		if !ok {
			continue
		}
		eth = append(eth, w)
	}
	return structs.EthosAverage(eth...).Clamp()
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
