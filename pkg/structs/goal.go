package structs

type Goal string

const (
	GoalWealth    Goal = "wealth"    // all the glitters is gold
	GoalGrowth    Goal = "growth"    // expand the population
	GoalTerritory Goal = "territory" // expand the empire
	GoalPower     Goal = "power"     // destroy rivals
	GoalPiety     Goal = "piety"     // spread the faith
	GoalResearch  Goal = "research"  // research
	GoalMilitary  Goal = "military"  // build an army
	GoalDiplomacy Goal = "diplomacy" // make friends
	GoalEspionage Goal = "espionage" // spy on rivals
	GoalStability Goal = "stability" // keep the peace
)
