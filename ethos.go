package ethos

import (
	"math"
)

// Ethic is some principal one can abide by.
type Ethic int

const (
	// number of ethics
	ethicCount = 6
)

const (
	// Altruism is unselfish concern for the welfare of others.
	//
	// High altruism implies selflessness, self sacrifice etc.
	// Low altruism implies selfishness, the complete lack of concern for others.
	EthicAltruism Ethic = iota

	// Ambition is the desire to get ahead in society, to obtain riches, honors, power etc.
	//
	// High ambition implies the willingness to go the extra mile, to work hard, to strive upwards.
	// Low ambition implies the lack of desire to improve ones station.
	EthicAmbition

	// Tradition is a measure of ones desire to stay within the laws and traditions of one's society, culture, laws etc.
	//
	// High tradition implies a (generally) law abiding outlook, great value placed on shared culture & values.
	// Low tradition implies a more chaotic, devil-may-care outlook, considering tradition(s) too confining and
	// (possibly even) laws too restrictive.
	EthicTradition

	// Pacifism is dedication to peace, eschewing violence & conflict.
	//
	// High pacifism implies one takes great pains avoid harming others, possibly even preferring death.
	// Low pacifism implies a strong propensity to violence.
	EthicPacifisim

	// Piety is faith in the divine, religious devotion.
	//
	// High piety implies strict adherence to ones faith & it's tenants.
	// Low piety implies no adherence to a faith.
	EthicPiety

	// Caution is the propensity is calculate carefully & weigh up risks before acting.
	//
	// High caution implies very deliberate, well thought out choices, multiple safeguards and counter strategies.
	// Low caution implies recklessness, the propensity to act without thinking; "there is no plan"
	EthicCaution
)

// Ethos is a general outlook for individuals or groups (averages).
type Ethos [ethicCount]int

// Distance between two ethos values
func (e *Ethos) Distance(o *Ethos) float64 {
	sum := 0.0
	for i := range e {
		sum += math.Pow(float64(e[i])-float64(o[i]), 2)
	}
	return math.Sqrt(sum)
}

// Average
func Average(in ...*Ethos) *Ethos {
	// TODO
	return &Ethos{}
}

// EthosDeviation represents variation (std. deviation) from an average in a populace.
// https://en.wikipedia.org/wiki/Standard_deviation
type EthosDeviation [ethicCount]float64
