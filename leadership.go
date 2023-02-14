package faction

type LeaderType int

// LeaderTypes represent how a faction is run.
const (
	LeaderTypeSingle  LeaderType = iota // there is a single leader
	LeaderTypeDual                      // there are two leaders (typically to counter each other)
	LeaderTypeTriad                     // there are three equal leaders (majority rules)
	LeaderTypeCouncil                   // all those with "great" to "excellent" (75+) affiliation
	LeaderTypeAll                       // everyone in the faction can vote
)
