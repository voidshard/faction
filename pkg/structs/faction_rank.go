package structs

type FactionRank int

// FactionRank denotes the rank of a person in a faction
const (
	FactionRankAssociate FactionRank = iota
	FactionRankInitiate
	FactionRankApprentice
	FactionRankNovice
	FactionRankJourneyman
	FactionRankAdept
	FactionRankExpert
	FactionRankMaster
	FactionRankGrandMaster
	FactionRankElder
	FactionRankRuler
)
