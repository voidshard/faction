package structs

// MetaKey is a key used to store additional metadata about something
// referenced by an accompanying MetaVal.
type MetaKey string

const (
	MetaKeyPerson     MetaKey = "person"     // targets specific person (ie. assassination) Key:PERSON_ID
	MetaKeyPlot       MetaKey = "plot"       // targets specific plot (ie. raid) Key:PLOT_ID
	MetaKeyResearch   MetaKey = "research"   // targets specific research Key:RESEARCH_TOPIC
	MetaKeyFaction    MetaKey = "faction"    // targets specific faction (ie. war) Key:FACTION_ID
	MetaKeyReligion   MetaKey = "religion"   // targets specific religion (ie. conversion) Key:RELIGION_ID
	MetaKeyGovernment MetaKey = "government" // targets specific government (ie. coup) Key:GOVERNMENT_ID
	MetaKeyFamily     MetaKey = "family"     // targets specific family Key:FAMILY_ID
	MetaKeyCommodity  MetaKey = "commodity"  // targets specific commodity (ie. spice) Key:COMMODITY_ID
	MetaKeyAction     MetaKey = "action"     // targets specific action (ie. assassination) Key:ACTION_TYPE
)
