package structs

type LeaderType int

// LeaderTypes represent how a faction is run.
const (
	LeaderTypeSingle  LeaderType = iota // there is a single leader
	LeaderTypeDual                      // there are two leaders (typically to counter each other)
	LeaderTypeTriad                     // there are three equal leaders (majority rules)
	LeaderTypeCouncil                   // all those with "great" to "excellent" (75+) affiliation
	LeaderTypeAll                       // everyone in the faction can vote
)

type LeaderStructure int

// LeaderStructure represents how a faction is structured
const (
	// LeaderStructurePyramid is a pyramid structure, with each rank up the structure getting
	// successively smaller in size at a constant rate.
	// Ie. 1 -> 3 -> 9 -> 27 -> 81
	// Example use: an organised army, government
	LeaderStructurePyramid LeaderStructure = iota

	// LeaderStructureLoose is a loose structure, where people are given rank based on
	// affiliation alone. Ie there isn't a limited number of slots per rank.
	// Example use: a mob, where rank is defined by how many kills / battles you've been in
	LeaderStructureLoose

	// LeaderStructureCell is a cell structure, each group of members functions autonomously
	// with only the cell leader knowing the next highest in rank.
	// In addition a few ranks are simply go betweens to the next level, hiding leadership.
	// Example use: a criminal organisation, cult
	LeaderStructureCell
)

// Rulers returns the number of rulers for a given leader type.
// For 'All' this will always return 1 (technically, everyone).
func (t LeaderType) Rulers() int {
	switch t {
	case LeaderTypeSingle:
		return 1
	case LeaderTypeDual:
		return 2
	case LeaderTypeTriad:
		return 3
	case LeaderTypeCouncil:
		return 7
	case LeaderTypeAll:
		return 1
	}
	return 0
}
