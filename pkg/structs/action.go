package structs

type ActionType string

// Actions that factions may perform. During a given tick each faction queues up some number of actions to be
// carried out.
//
// Actions aren't absolute, they imply *focus*
// Ie. 'Craft' implies a faction pushes more than usual this tick to make *more* goods over and above their
// general operations (ie. survival). Recruitment implies a focus this tick on canvasing people to
// join / signup over and above the usual; not that people aren't considering joining / recruiting at any
// given tick.
//
// Each Action alters some of the main faction variables, some just on the user, some on
// user & target, and some on the people in the given area(s) in which the action is used.
// Ie, one or more of;
// - weatlh
// - cohession
// - corruption
// - property
// - favor / trust (of factions)
// - favor / trust (of people)
// - attack / defense (espionage and/or military)
// - affiliation (of people)
// - research
//
// Government Only Actions (faction.IsGovernment must be true)
// - revoke land
//
// Religion Only Actions (faction.IsReligion must be true):
// - crusade, excommunicate
//
// Legal (non covert) non-Government factions:
// - requestland
const (
	// Friendly actions (some of these target another faction)
	ActionTypeTrade       ActionType = "trade"       // trade goods with another faction, everyone wins
	ActionTypeBribe       ActionType = "bribe"       // pay a faction to increase their favor
	ActionTypeFestival    ActionType = "festival"    // hold a festival (usually religious in nature), increases general favor & affiliation
	ActionTypeRitual      ActionType = "ritual"      // hold a public ritual, increases faith
	ActionTypeRequestLand ActionType = "requestland" // request the government grant the use of some land
	ActionTypeCharity     ActionType = "charity"     // donate money to people of an area(s), increases general favor

	// Neutral actions (these target "self")
	ActionTypePropoganda     ActionType = "propoganda"      // increases general favor towards you, cheaper than 'Charity' but less um, honest
	ActionTypeRecruit        ActionType = "recruit"         // push recuitment (overtly or covertly), increases general affiliation
	ActionTypeExpand         ActionType = "expand"          // purchase more land / property, add base of operations in another city etc
	ActionTypeDownsize       ActionType = "downsize"        // sell land / property to recoup funds
	ActionTypeCraft          ActionType = "craft"           // focus on crafting, increasing funds
	ActionTypeHarvest        ActionType = "harvest"         // focus on harvesting, increasing crop yield, mining - increasing funds
	ActionTypeConsolidate    ActionType = "consolidate"     // consolidate; internal re-organisation, process streamlining etc (increases cohession)
	ActionTypeResearch       ActionType = "research"        // adds some science & funds
	ActionTypeExcommunicate  ActionType = "excommunicate"   // religion explicitly excommunicates someone; person loses favour & affiliation
	ActionTypeConcealSecrets ActionType = "conceal-secrets" // bribes are paid, people are silenced, documents burned (+secrecy)

	// Unfriendly actions (all of these have a target faction)
	ActionTypeGatherSecrets   ActionType = "gather-secrets"   // attempt to discover secrets of target faction
	ActionTypeRevokeLand      ActionType = "revoke-land"      // government retracts right to use some land
	ActionTypeHireMercenaries ActionType = "hire-mercenaries" // hire another faction to do something for you
	ActionTypeHireSpies       ActionType = "hire-spies"       // hire another faction to spy on another faction
	ActionTypeSpreadRumors    ActionType = "spread-rumors"    // inverse of Propoganda; decrease general favor toward target
	ActionTypeAssassinate     ActionType = "assassinate"      // someone is selected for elimination
	ActionTypeFrame           ActionType = "frame"            // the non religious version of 'excommunicate,' but can involve legal .. entanglements
	ActionTypeRaid            ActionType = "raid"             // small armed conflict with the aim of destroying as much as possible
	ActionTypeEnslave         ActionType = "enslave"          // similar to raid, but with the aim of capturing people
	ActionTypeSteal           ActionType = "steal"            // similar to raid, but with less stabbing
	ActionTypePillage         ActionType = "pillage"          // small armed conflict with the aim of stealing wealth
	ActionTypeBlackmail       ActionType = "blackmail"        // trade in a secret for a pile of gold, or ruin a reputation
	ActionTypeKidnap          ActionType = "kidnap"           // similar to blackmail, but you trade back a person

	// Hostile actions
	ActionTypeShadowWar ActionType = "shadow-war" // shadow war is a full armed conflict carried out between one or more covert factions
	ActionTypeCrusade   ActionType = "crusade"    // excommunication on a grander and more permanent scale
	ActionTypeWar       ActionType = "war"        // full armed conflict
)

var (
	// Factions with a IsReligion & ReligionID may perform these actions
	ActionsReligionOnly = []ActionType{}

	// Factions with a IsGovernment & GovernmentID may perform these actions
	ActionsGovernmentOnly = []ActionType{}

	// Factions that are not Governments nor convert may perform these actions
	ActionsLegalFactionOnly = []ActionType{}

	// TargetType, if not given default to no target (ie. "self" or "we don't need to choose a target")
	// nb. Jobs always include a Faction and Area target, so we only do more work in picking a target if it's something else.
	//     That is we do not include MetaKeyFaction or MetaKeyTarget here, since it's always required, this is for actions
	//     that need *another* more specific target.
	ActionTarget = map[ActionType]MetaKey{
		ActionTypeExcommunicate:   MetaKeyPerson,   // excommunicate someone
		ActionTypeBribe:           MetaKeyPerson,   // pay someone to buff favor with them & their faction
		ActionTypeAssassinate:     MetaKeyPerson,   // assassinate someone
		ActionTypeFrame:           MetaKeyPerson,   // frame someone
		ActionTypeBlackmail:       MetaKeyPerson,   // blackmail someone
		ActionTypeKidnap:          MetaKeyPerson,   // kidnap someone
		ActionTypeRequestLand:     MetaKeyPlot,     // ask government to grant land
		ActionTypeRevokeLand:      MetaKeyPlot,     // revoke land from a faction
		ActionTypeDownsize:        MetaKeyPlot,     // sell a plot in some area
		ActionTypeSteal:           MetaKeyPlot,     // steal from some building
		ActionTypeRaid:            MetaKeyPlot,     // raid an area, centered on a plot
		ActionTypePillage:         MetaKeyPlot,     // pillage an area, centered on some plot
		ActionTypeEnslave:         MetaKeyPlot,     // enslave people in some area
		ActionTypeRitual:          MetaKeyPlot,     // perform a ritual in given building
		ActionTypeGatherSecrets:   MetaKeyFaction,  // gather secrets about some faction from people in some area / plot
		ActionTypeResearch:        MetaKeyResearch, // research some tech
		ActionTypeHireMercenaries: MetaKeyJob,      // create a job (MetaTarget) for some other faction, to attack *another* faction (TargetFactionID)
		ActionTypeHireSpies:       MetaKeyJob,      // create a job (MetaTarget) for some other faction, to spy on *another* faction (TargetFactionID)
	}

	// Actions that may be performed by mercenaries (eg. HireMercenaries action)
	ActionsForMercenaries = []ActionType{
		ActionTypeRaid,
		ActionTypePillage,
	}

	//
	ActionsForSpies = []ActionType{
		ActionTypeFrame,
		ActionTypeAssassinate,
		ActionTypeBlackmail,
	}

	// All actions known to us
	AllActions = []ActionType{
		ActionTypeTrade,
		ActionTypeBribe,
		ActionTypeFestival,
		ActionTypeRitual,
		ActionTypeRequestLand,
		ActionTypeCharity,
		ActionTypePropoganda,
		ActionTypeRecruit,
		ActionTypeExpand,
		ActionTypeDownsize,
		ActionTypeCraft,
		ActionTypeHarvest,
		ActionTypeConsolidate,
		ActionTypeResearch,
		ActionTypeExcommunicate,
		ActionTypeConcealSecrets,
		ActionTypeGatherSecrets,
		ActionTypeRevokeLand,
		ActionTypeHireMercenaries,
		ActionTypeHireSpies,
		ActionTypeSpreadRumors,
		ActionTypeAssassinate,
		ActionTypeFrame,
		ActionTypeRaid,
		ActionTypeEnslave,
		ActionTypeSteal,
		ActionTypePillage,
		ActionTypeBlackmail,
		ActionTypeKidnap,
		ActionTypeShadowWar,
		ActionTypeCrusade,
		ActionTypeWar,
	}
)

func init() {
	// Factions with a IsReligion & ReligionID may perform these actions
	ActionsReligionOnly = []ActionType{ActionTypeCrusade, ActionTypeExcommunicate, ActionTypeRitual}

	// Factions with a IsGovernment & GovernmentID may perform these actions
	ActionsGovernmentOnly = []ActionType{ActionTypeRevokeLand}

	// Factions that are not Governments nor convert may perform these actions
	ActionsLegalFactionOnly = []ActionType{ActionTypeRequestLand}
}
