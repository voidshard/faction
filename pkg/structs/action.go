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
// - grant land, revoke land
//
// Religion Only Actions (faction.IsReligion must be true):
// - crusade, excommunicate
const (
	// Friendly actions (most of these target another faction)
	ActionTypeTrade     ActionType = "trade"     // trade goods with another faction, everyone wins
	ActionTypeBribe     ActionType = "bribe"     // pay a faction to increase their favor
	ActionTypeFestival  ActionType = "festival"  // hold a festival (usually religious in nature), increases general favor & affiliation
	ActionTypeGrantLand ActionType = "grantland" // government grants the use of some land
	ActionTypeCharity   ActionType = "charity"   // donate money to people of an area(s), increases general favor

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
	ActionTypeSpreadRumors    ActionType = "spread-rumors"    // inverse of Propoganda; decrease general favor toward target
	ActionTypeAssassinate     ActionType = "assassinate"      // someone is selected for elimination
	ActionTypeFrame           ActionType = "frame"            // the non religious version of 'excommunicate,' but can involve legal .. entanglements
	ActionTypeRaid            ActionType = "raid"             // small armed conflict with the aim of destroying as much as possible
	ActionTypeSteal           ActionType = "steal"            // similar to raid, but with less stabbing
	ActionTypePillage         ActionType = "pillage"          // small armed conflict with the aim of stealing wealth
	ActionTypeBlackmail       ActionType = "blackmail"        // trade in a secret for a pile of gold, or ruin a reputation
	ActionTypeKidnap          ActionType = "kidnap"           // similar to blackmail, but you trade back a person

	// Hostile actions
	ActionTypeShadowWar ActionType = "shadow-war" // shadow war is a full armed conflict carried out between one or more covert factions
	ActionTypeCrusade   ActionType = "crusade"    // excommunication on a grander and more permanent scale
	ActionTypeWar       ActionType = "war"        // full armed conflict
)
