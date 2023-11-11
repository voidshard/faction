package config

// Government is the configuration for randomly creating a government.
type Government struct {
	// Probability that the given action will be outlawed
	ProbabilityOutlawAction map[string]float64

	// Probability that the given commodity will be outlawed
	ProbabilityOutlawCommodity map[string]float64

	// Probability that the given research topic will be outlawed
	ProbabilityOutlawResearch map[string]float64

	// Probability that the given religion will be outlawed
	ProbabilityOutlawReligion map[string]float64

	// How often (in ticks) the government will collect taxes
	// Min: 1 (every tick)
	// Probably you want this to be every month, quarter or something
	// (however many ticks that is).
	TaxFrequency Distribution

	// Rate (converted to a %, so this should be 0-100).
	//
	// Non-covert factions will pay this rate from their wealth into the government's
	// coffers. High rates make factions increasingly angry.
	//
	// Min: 0
	// Max: 100
	TaxRate Distribution

	// Weight to add to "goverment" actions
	// - RevokeLand, GrantLand
	ActionWeight Distribution

	// Grant governments a buff to military / espionage
	// 0.25 = 25% bonus
	// -0.5 = 50% penalty
	// .. etc
	MilitaryOffenseBonus  float64
	MilitaryDefenseBonus  float64
	EspionageOffenseBonus float64
	EspionageDefenseBonus float64
}
