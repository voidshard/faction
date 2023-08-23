package structs

import (
	"fmt"
	"math"
)

const (
	// MaxEthos is the highest possible ethos value
	MaxEthos = 10000

	// MinEthos is the lowest possible ethos value
	MinEthos = -10000
)

var (
	// max possible distance between two ethos values
	maxDist = ethosDistance((&Ethos{}).Add(MinEthos), (&Ethos{}).Add(MaxEthos))
)

// Ethos is a set of guiding principles that someone (or a faction) abides by,
// and how strongly the do so (or do not).
type Ethos struct {
	// Altruism is unselfish concern for the welfare of others.
	//
	// High altruism implies selflessness, self sacrifice etc.
	// Low altruism implies selfishness, the complete lack of concern for others.
	Altruism int `db:"ethos_altruism"`

	// Ambition is the desire to get ahead in society, to obtain riches, honors, power etc.
	//
	// High ambition implies the willingness to go the extra mile, to work hard, to strive upwards.
	// Low ambition implies the lack of desire to improve ones station.
	Ambition int `db:"ethos_ambition"`

	// Tradition is a measure of ones desire to stay within the laws and traditions of one's society, culture, laws etc.
	//
	// High tradition implies a (generally) law abiding outlook, great value placed on shared culture & values.
	// Low tradition implies a more chaotic, devil-may-care outlook, considering tradition(s) too confining and
	// (possibly even) laws too restrictive.
	Tradition int `db:"ethos_tradition"`

	// Pacifism is dedication to peace, eschewing violence & conflict.
	//
	// High pacifism implies one takes great pains avoid harming others, possibly even preferring death.
	// Low pacifism implies a strong propensity to violence.
	Pacifism int `db:"ethos_pacifism"`

	// Piety is faith in the divine, religious devotion.
	//
	// High piety implies strict adherence to ones faith & it's tenants.
	// Low piety implies no adherence to a faith.
	Piety int `db:"ethos_piety"`

	// Caution is the propensity is calculate carefully & weigh up risks before acting.
	//
	// High caution implies very deliberate, well thought out choices, multiple safeguards and counter strategies.
	// Low caution implies recklessness, the propensity to act without thinking; "there is no plan"
	Caution int `db:"ethos_caution"`
}

func (e *Ethos) String() string {
	return fmt.Sprintf("Ethos{Amb:%d, Alt:%d, Tra:%d, Pac:%d, Pie:%d, Cau:%d}", e.Ambition, e.Altruism, e.Tradition, e.Pacifism, e.Piety, e.Caution)
}

// Add v to Ethos values returning a new Ethos
func (e *Ethos) Add(v int) *Ethos {
	return (&Ethos{
		Ambition:  int(e.Ambition + v),
		Altruism:  int(e.Altruism + v),
		Tradition: int(e.Tradition + v),
		Pacifism:  int(e.Pacifism + v),
		Piety:     int(e.Piety + v),
		Caution:   int(e.Caution + v),
	}).Clamp()
}

// AddEthos adds ethos values returning a new Ethos
func (e *Ethos) AddEthos(v *Ethos) *Ethos {
	return (&Ethos{
		Ambition:  int(e.Ambition + v.Ambition),
		Altruism:  int(e.Altruism + v.Altruism),
		Tradition: int(e.Tradition + v.Tradition),
		Pacifism:  int(e.Pacifism + v.Pacifism),
		Piety:     int(e.Piety + v.Piety),
		Caution:   int(e.Caution + v.Caution),
	}).Clamp()
}

// Multiply ethos values by v returning a new Ethos
func (e *Ethos) Multiply(v float64) *Ethos {
	return (&Ethos{
		Ambition:  int(float64(e.Ambition) * v),
		Altruism:  int(float64(e.Altruism) * v),
		Tradition: int(float64(e.Tradition) * v),
		Pacifism:  int(float64(e.Pacifism) * v),
		Piety:     int(float64(e.Piety) * v),
		Caution:   int(float64(e.Caution) * v),
	}).Clamp()
}

// Clamp ethos values to min / max values
func (e *Ethos) Clamp() *Ethos {
	e.Ambition = int(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Ambition))))
	e.Altruism = int(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Altruism))))
	e.Tradition = int(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Tradition))))
	e.Pacifism = int(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Pacifism))))
	e.Piety = int(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Piety))))
	e.Caution = int(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Caution))))
	return e
}

// EthosDistance returns the distance between two ethos values as a value 0.0 - 1.0
func EthosDistance(a, b *Ethos) float64 {
	return ethosDistance(a, b) / maxDist
}

// ethosDistance returns the distance between two ethos values in absolute terms
func ethosDistance(a, b *Ethos) float64 {
	a.Clamp()
	b.Clamp()

	dist := 0.0
	count := 0

	for _, j := range [][2]int{
		{a.Altruism, b.Altruism},
		{a.Ambition, b.Ambition},
		{a.Tradition, b.Tradition},
		{a.Pacifism, b.Pacifism},
		{a.Piety, b.Piety},
		{a.Caution, b.Caution},
	} {
		if j[0] == 0 && j[1] == 0 {
			continue
		}
		count++
		if j[0] > j[1] {
			dist += float64(j[0] - j[1])
		} else {
			dist += float64(j[1] - j[0])
		}
	}

	if count == 0 {
		return 0.0
	}

	return dist / float64(count)
}

// EthosAverage returns the average of the given ethos values
func EthosAverage(in ...*Ethos) *Ethos {
	e := &Ethos{}
	if len(in) == 0 {
		return e
	}

	for _, i := range in {
		i.Clamp()
		e.Altruism += i.Altruism
		e.Ambition += i.Ambition
		e.Tradition += i.Tradition
		e.Pacifism += i.Pacifism
		e.Piety += i.Piety
		e.Caution += i.Caution
	}
	e.Altruism /= len(in)
	e.Ambition /= len(in)
	e.Tradition /= len(in)
	e.Pacifism /= len(in)
	e.Piety /= len(in)
	e.Caution /= len(in)
	return e
}
