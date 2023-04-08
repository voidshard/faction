package sim

// simOption controls internal settings of the simulation beyond a
// what is contained in a config file.
// Ie. we can register interfaces, live services etc.
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
