package config

import (
	"github.com/voidshard/faction/pkg/structs"
)

// SocialStructure denotes how "classes" are decided & social mobility (if any).
// We do this by calculating some number of "social points" for "good" (nb. meaning 'considered good
// by the social elite'). Families (not individuals) are then assigned a social class based on their score.
//
// We can score by:
// - skill in profession(s)
// - faith in some religion
// - rank(s) in faction(s)
// - gender
// - race
// - culture
// - area
// - wealth ** TODO: we don't store family wealth yet
type SocialStructure struct {
	// Base number of points for members of this culture (ie. "every Greek gets 10 points"),
	Base float64

	// Social classes are assigned based on the number of "social points" a family has.
	// A familiy is awarded the highest class they qualify for.
	//
	// Ie. {"upper": 100, "middle": 50, "lower": 0}
	// Thus a family with 87 points is "middle" and 20 would be "lower"
	//
	// The default class is the one with the lowest number.
	//
	// Nb. Each class should have a unique name & value.
	Classes map[string]int

	// Assign points by profession. This is a map
	//
	// Map: Profession name -> ratio
	// Points awarded: skill_at_profession * ratio
	//
	// As a special case we accept "" as "any profession"
	Profession map[string]float64

	// Multiply profession points by whether the profession is legal or not.
	ProfessionMultLegal   float64
	ProfessionMultIllegal float64

	// Assign points by faith.
	//
	// Map: Faith name -> ratio
	// Points awarded: faith_level * ratio
	//
	// As a special case we accept "" as "any faith"
	Faith map[string]float64

	// Multiply faith points by whether the faith is legal or not.
	FaithMultLegal   float64
	FaithMultIllegal float64

	// Assign points by faction rank.
	//
	// faction_id -> rank -> ratio
	// Points awarded: rank * ratio
	//
	// As a special case we accept "" as "any faction"
	Faction map[string]map[structs.FactionRank]float64

	// Multiply faction points by whether the faction is legal or not.
	FactionMultLegal   float64
	FactionMultIllegal float64

	// Assign points by gender.
	//
	// Here the bool means "is_male"
	//
	// Map: is_male -> points
	// Points awarded: points
	Gender map[bool]float64

	// Area grants points based on where the family is based (at least, at the time of assignment).
	Area map[string]float64 // area_id -> points
}
