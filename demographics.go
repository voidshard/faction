package faction

// Demographics roughly describes a large population.
//
// For randomly making societies that look "sort of like this."
type Demographics struct {
	// Average outlook of members of the population
	Ethos *Ethos

	// EthosDeviation is the standard deviation of the populace
	EthosDeviation *EthosDeviation
}
