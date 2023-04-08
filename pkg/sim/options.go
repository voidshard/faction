package sim

// simOption controls internal settings of the simulation beyond a
// what is contained in a config file.
//
// The goal here is to support what might be highly volatile, differing based
// on location, time or whatever to be dictated by the user as the sim goes on.
// Ie. over time more technologies might become available, or more commodities
// as either is discovered.
// These might also differ based on where the caller is.
type simOption func(Simulation) error

// SetTechnology registers the given technology tree with the simulation.
func SetTechnology(tech Technology) simOption {
	return func(s Simulation) error {
		s.(*simulationImpl).tech = tech
		return nil
	}
}

// SetEconomy registers the given economy with the simulation.
func SetEconomy(eco Economy) simOption {
	return func(s Simulation) error {
		s.(*simulationImpl).eco = eco
		return nil
	}
}
