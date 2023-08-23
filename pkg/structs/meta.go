package structs

// MetaKey is a key used to store additional metadata about something
// referenced by an accompanying MetaVal.
type MetaKey string

const (
	MetaKeyPerson     MetaKey = "person"     // targets specific person Val:PERSON_ID
	MetaKeyPlot       MetaKey = "plot"       // targets specific plot (ie. building_a) Val:PLOT_ID
	MetaKeyResearch   MetaKey = "research"   // targets specific research (ie. physics) Val:RESEARCH_TOPIC
	MetaKeyFaction    MetaKey = "faction"    // targets specific faction Val:FACTION_ID
	MetaKeyReligion   MetaKey = "religion"   // targets specific religion Val:RELIGION_ID
	MetaKeyGovernment MetaKey = "government" // targets specific government Val:GOVERNMENT_ID
	MetaKeyFamily     MetaKey = "family"     // targets specific family Val:FAMILY_ID
	MetaKeyCommodity  MetaKey = "commodity"  // targets specific commodity (ie. spice) Val:COMMODITY_ID
	MetaKeyAction     MetaKey = "action"     // targets specific action (ie. assassination) Val:ACTION_TYPE
	MetaKeyJob        MetaKey = "job"        // targets specific job Val:JOB_ID
	MetaKeyArea       MetaKey = "area"       // targets specific area Val:AREA_ID
	MetaKeyRoute      MetaKey = "route"      // targets specific route Val:ROUTE_ID
	MetaKeyEvent      MetaKey = "event"      // targets specific event (ie. person_birth) Val:EVENT_ID
	MetaKeyGoal       MetaKey = "goal"       // targets specific goal (ie. wealth) Val:GOAL_ID
)
