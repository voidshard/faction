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

func (e *Ethos) ToString() string {
	return fmt.Sprintf("Ethos{Amb:%d, Alt:%d, Tra:%d, Pac:%d, Pie:%d, Cau:%d}", e.Ambition, e.Altruism, e.Tradition, e.Pacifism, e.Piety, e.Caution)
}

// Add v to Ethos values returning a new Ethos
func (e *Ethos) Add(v int64) *Ethos {
	return (&Ethos{
		Ambition:  int64(e.Ambition + v),
		Altruism:  int64(e.Altruism + v),
		Tradition: int64(e.Tradition + v),
		Pacifism:  int64(e.Pacifism + v),
		Piety:     int64(e.Piety + v),
		Caution:   int64(e.Caution + v),
	}).Clamp()
}

// AddEthos adds ethos values returning a new Ethos
func (e *Ethos) AddEthos(v *Ethos) *Ethos {
	return (&Ethos{
		Ambition:  int64(e.Ambition + v.Ambition),
		Altruism:  int64(e.Altruism + v.Altruism),
		Tradition: int64(e.Tradition + v.Tradition),
		Pacifism:  int64(e.Pacifism + v.Pacifism),
		Piety:     int64(e.Piety + v.Piety),
		Caution:   int64(e.Caution + v.Caution),
	}).Clamp()
}

// Multiply ethos values by v returning a new Ethos
func (e *Ethos) Multiply(v float64) *Ethos {
	return (&Ethos{
		Ambition:  int64(float64(e.Ambition) * v),
		Altruism:  int64(float64(e.Altruism) * v),
		Tradition: int64(float64(e.Tradition) * v),
		Pacifism:  int64(float64(e.Pacifism) * v),
		Piety:     int64(float64(e.Piety) * v),
		Caution:   int64(float64(e.Caution) * v),
	}).Clamp()
}

// Clamp ethos values to min / max values
func (e *Ethos) Clamp() *Ethos {
	e.Ambition = int64(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Ambition))))
	e.Altruism = int64(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Altruism))))
	e.Tradition = int64(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Tradition))))
	e.Pacifism = int64(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Pacifism))))
	e.Piety = int64(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Piety))))
	e.Caution = int64(math.Min(MaxEthos, math.Max(MinEthos, float64(e.Caution))))
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

	for _, j := range [][2]int64{
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
	e.Altruism /= int64(len(in))
	e.Ambition /= int64(len(in))
	e.Tradition /= int64(len(in))
	e.Pacifism /= int64(len(in))
	e.Piety /= int64(len(in))
	e.Caution /= int64(len(in))
	return e
}
