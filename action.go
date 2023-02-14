package action

type ActionType int

const (
	// Friendly actions
	ActionTypeTrade ActionType = iota
	ActionTypeBribe
	ActionTypeCharity
	ActionTypeFestival

	// Neutral actions
	ActionTypePropoganda
	ActionTypeRecruit
	ActionTypeExpand
	ActionTypeDownsize
	ActionTypeCraft
	ActionTypeHarvest
	ActionTypeConsolidate
	ActionTypeResearchScience
	ActionTypeResearchTheology
	ActionTypeResearchMagic
	ActionTypeResearchOccult
	ActionTypePassLaw
	ActionTypeLobbyGovernment
	ActionTypeHoldRitual
	ActionTypeExcommunicate
	ActionTypeGatherSecrets
	ActionTypeConcealSecrets

	// Unfriendly actions
	ActionTypeSpreadRumors
	ActionTypeAssassinate
	ActionTypeFrame
	ActionTypeRaid
	ActionTypePillage
	ActionTypeBlackmail
	ActionTypeKidnap

	// Hostile actions
	ActionTypeWar
	ActionTypeShadowWar
	ActionTypeCrusade
)
