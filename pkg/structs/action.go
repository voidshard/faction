package structs

type ActionType int

// trust
// affiliation

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
const (
	// Friendly actions (most of these target another faction)
	ActionTypeTrade     ActionType = iota // trade goods with another faction, everyone wins
	ActionTypeBribe                       // pay a faction to increase their favor
	ActionTypeFestival                    // hold a festival (usually religious in nature), increases general favor & affiliation
	ActionTypeGrantLand                   // government grants the use of some land
	ActionTypeCharity                     // donate money to people of an area(s), increases general favor

	// Neutral actions (these target "self")
	ActionTypePropoganda       // increases general favor towards you, cheaper than 'Charity' but less um, honest
	ActionTypeRecruit          // push recuitment (overtly or covertly), increases general affiliation
	ActionTypeExpand           // purchase more land / property, add base of operations in another city etc
	ActionTypeDownsize         // sell land / property to recoup funds
	ActionTypeCraft            // focus on crafting, increasing funds
	ActionTypeHarvest          // focus on harvesting, increasing crop yield, mining - increasing funds
	ActionTypeConsolidate      // consolidate; internal re-organisation, process streamlining etc (increases cohession)
	ActionTypeResearchScience  // adds some science & funds
	ActionTypeResearchTheology // adds some theology & funds
	ActionTypeResearchMagic    // adds some magic(al science) & funds
	ActionTypeResearchOccult   // adds some occult (research) & funds
	ActionTypeHoldMagicRitual  // culmination of magical research (effect varies)
	ActionTypeHoldOccultRitual // culmination of occult magical research (effect varies)
	ActionTypeExcommunicate    // religion explicitly excommunicates someone; person loses favour & affiliation
	ActionTypeConcealSecrets   // bribes are paid, people are silenced, documents burned (+secrecy)

	// Unfriendly actions (all of these have a target faction)
	ActionTypeGatherSecrets // attempt to discover secrets of target faction
	ActionTypeRevokeLand    // government retracts right to use some land
	ActionTypeSpreadRumors  // inverse of Propoganda; decrease general favor toward target
	ActionTypeAssassinate   // someone is selected for elimination
	ActionTypeFrame         // the non religious version of 'excommunicate,' but can involve legal .. entanglements
	ActionTypeRaid          // small armed conflict with the aim of destroying as much as possible
	ActionTypePillage       // small armed conflict with the aim of stealing wealth
	ActionTypeBlackmail     // trade in a secret for a pile of gold, or ruin a reputation
	ActionTypeKidnap        // similar to blackmail, but you trade back a person

	// Hostile actions
	ActionTypeShadowWar // shadow war is a full armed conflict carried out between one or more covert factions
	ActionTypeCrusade   // excommunication on a grander and more permanent scale
	ActionTypeWar       // full armed conflict
)
