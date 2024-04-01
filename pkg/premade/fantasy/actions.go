package fantasy

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// Actions returns a default action configuration for fantasy type settings.
//
// This could be useful as a starting point for your own configuration.
func Actions() map[string]*config.Action {
	var minor int64 = structs.MaxEthos / 10
	var modest int64 = structs.MaxEthos / 2
	var major int64 = structs.MaxEthos

	return map[string]*config.Action{
		Trade: {
			MinPeople:   3,
			MaxPeople:   15,
			Probability: 0.3,
			Category:    config.ActionCategoryFriendly,
			Ethos:       structs.Ethos{Altruism: -1 * minor},
			Cost: config.Distribution{
				Min:       1000000,  // 100gp
				Max:       50000000, // 5000gp
				Mean:      10000000, // 1000gp
				Deviation: 40000000, // 4000gp
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       5 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 4 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       10 * DEFAULT_TICKS_PER_DAY,
				Max:       90 * DEFAULT_TICKS_PER_DAY,
				Mean:      30 * DEFAULT_TICKS_PER_DAY,
				Deviation: 60 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				MERCHANT: 5,
				SCRIBE:   1,
				CLERK:    1,
				SOLDIER:  2,
				SAILOR:   3,
			},
			SecrecyWeight:  1,
			Goals:          []structs.Goal{structs.Goal_Wealth},
			TargetMinTrust: structs.MaxEthos / 10,
			TargetMaxTrust: structs.MaxEthos,
		},
		Bribe: {
			MinPeople:   1,
			MaxPeople:   15,
			Category:    config.ActionCategoryFriendly,
			Target:      structs.Meta_KeyPerson,
			Probability: 0.05,
			Ethos:       structs.Ethos{Tradition: -2 * minor, Ambition: 2 * minor, Caution: 2 * minor},
			Cost: config.Distribution{
				Min:       20000000,  // 2000gp
				Max:       100000000, // 10000gp
				Mean:      25000000,  // 2500gp
				Deviation: 70000000,  // 7000gp
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      7 * DEFAULT_TICKS_PER_DAY,
				Deviation: 7 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 2 * DEFAULT_TICKS_PER_DAY,
			},
			ProfessionWeights: map[string]float64{
				CLERK:    1,
				SCRIBE:   1,
				SPY:      5,
				ASSASSIN: 2,
				NOBLE:    2,
			},
			SecrecyWeight:  1.5,
			Goals:          []structs.Goal{structs.Goal_Diplomacy},
			TargetMinTrust: structs.MinEthos / 4,
			TargetMaxTrust: structs.MaxEthos / 10,
		},
		Festival: {
			MinPeople: 100,
			MaxPeople: -1,
			Category:  config.ActionCategoryFriendly,
			Ethos:     structs.Ethos{Altruism: modest, Tradition: minor, Piety: modest},
			Cost: config.Distribution{
				Min:       1000000,  // 100gp
				Max:       10000000, // 1000gp
				Mean:      5000000,  // 500gp
				Deviation: 5000000,  // 500gp
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      21 * DEFAULT_TICKS_PER_DAY,
				Deviation: 15 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 2 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				PRIEST: 5,
				CLERK:  2,
			},
			Goals: []structs.Goal{structs.Goal_Piety, structs.Goal_Diplomacy},
		},
		Ritual: {
			MinPeople:   15,
			MaxPeople:   -1,
			Target:      structs.Meta_KeyPlot,
			JobPriority: 50,
			Ethos:       structs.Ethos{Tradition: modest, Piety: major},
			Restricted: [][]config.Condition{
				{config.ConditionSrcFactionIsReligion},
			},
			Cost: config.Distribution{
				Min:       100000,   // 10gp
				Max:       10000000, // 1000gp
				Mean:      250000,   // 500gp
				Deviation: 2000000,  // 200gp
			},
			TimeToPrepare: config.Distribution{
				Min:       3 * DEFAULT_TICKS_PER_DAY,
				Max:       5 * DEFAULT_TICKS_PER_DAY,
				Mean:      4 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 3 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				PRIEST: 5,
				CLERK:  2,
			},
			Goals: []structs.Goal{structs.Goal_Piety},
		},
		RequestLand: {
			MinPeople:   1,
			MaxPeople:   12,
			Target:      structs.Meta_KeyPlot,
			Probability: 0.1,
			Ethos:       structs.Ethos{Altruism: minor, Tradition: 2 * minor, Ambition: minor},
			Restricted: [][]config.Condition{
				{config.ConditionSrcFactionIsNotGovernment, config.ConditionSrcFactionIsNotCovert},
			},
			TimeToPrepare: config.Distribution{
				Min:       15 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      20 * DEFAULT_TICKS_PER_DAY,
				Deviation: 15 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       30 * DEFAULT_TICKS_PER_DAY,
				Max:       90 * DEFAULT_TICKS_PER_DAY,
				Mean:      60 * DEFAULT_TICKS_PER_DAY,
				Deviation: 60 * DEFAULT_TICKS_PER_DAY,
			},
			Goals:          []structs.Goal{structs.Goal_Growth},
			TargetMinTrust: structs.MaxEthos / 2,
			TargetMaxTrust: structs.MaxEthos,
		},
		Charity: {
			MinPeople:   5,
			MaxPeople:   1500,
			Category:    config.ActionCategoryFriendly,
			Probability: 0.05,
			Ethos:       structs.Ethos{Altruism: major, Piety: minor},
			Cost: config.Distribution{
				Min:       100000,   // 10gp
				Max:       10000000, // 1000gp
				Mean:      5000000,  // 500gp
				Deviation: 5000000,  // 500gp
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      21 * DEFAULT_TICKS_PER_DAY,
				Deviation: 15 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      10 * DEFAULT_TICKS_PER_DAY,
				Deviation: 7 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 2,
			Goals:        []structs.Goal{structs.Goal_Piety, structs.Goal_Diplomacy},
		},
		Propoganda: {
			MinPeople:   30,
			MaxPeople:   -1,
			Probability: 0.05,
			Ethos:       structs.Ethos{Tradition: -1 * minor, Ambition: modest},
			Cost: config.Distribution{
				Min:       50000000,  // 5000gp
				Max:       100000000, // 10000gp
				Mean:      75000000,  // 7500gp
				Deviation: 20000000,  // 2000gp
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      21 * DEFAULT_TICKS_PER_DAY,
				Deviation: 14 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       90 * DEFAULT_TICKS_PER_DAY,
				Max:       180 * DEFAULT_TICKS_PER_DAY,
				Mean:      135 * DEFAULT_TICKS_PER_DAY,
				Deviation: 50 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				SCRIBE: 4,
				PRIEST: 2,
				NOBLE:  5,
			},
			SecrecyWeight:  1.1,
			Goals:          []structs.Goal{structs.Goal_Power},
			TargetMinTrust: structs.MinEthos / 2,
			TargetMaxTrust: structs.MaxEthos / 10,
		},
		Recruit: {
			MinPeople:    1,
			MaxPeople:    60,
			Probability:  0.2,
			JobPriority:  80,
			Ethos:        structs.Ethos{Caution: modest},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				SCRIBE: 4,
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 5 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       30 * DEFAULT_TICKS_PER_DAY,
				Max:       90 * DEFAULT_TICKS_PER_DAY,
				Mean:      60 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1,
			Goals:         []structs.Goal{structs.Goal_Growth},
		},
		Expand: { // cost determined by land prices in economy
			MinPeople:    1,
			MaxPeople:    15,
			Probability:  0.2,
			Ethos:        structs.Ethos{Caution: 3 * minor},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				NOBLE:  5,
				SCRIBE: 4,
			},
			TimeToPrepare: config.Distribution{
				Min:       3 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      14 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       3 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      14 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1,
			Goals:         []structs.Goal{structs.Goal_Territory},
		},
		Craft: {
			MinPeople:    1,
			MaxPeople:    -1,
			Probability:  0.3,
			JobPriority:  50,
			Ethos:        structs.Ethos{Altruism: -1 * minor},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				WEAVER:        3,
				CLOTHIER:      3,
				TANNER:        1,
				SMELTER:       1,
				LEATHERWORKER: 3,
				CARPENTER:     3,
				SMITH:         3,
				ALCHEMIST:     3,
				MAGE:          3,
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       3 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      4 * DEFAULT_TICKS_PER_DAY,
				Deviation: 4 * DEFAULT_TICKS_PER_DAY,
			},
			Goals: []structs.Goal{structs.Goal_Wealth},
		},
		Harvest: {
			MinPeople:    1,
			MaxPeople:    -1,
			Restricted:   [][]config.Condition{{config.ConditionSrcFactionHasHarvestablePlot}},
			Probability:  0.5,
			JobPriority:  55,
			Ethos:        structs.Ethos{},
			PersonWeight: 2,
			ProfessionWeights: map[string]float64{
				FARMER:    4,
				MINER:     4,
				FISHERMAN: 4,
				HUNTER:    4,
				FORESTER:  4,
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 4 * DEFAULT_TICKS_PER_DAY,
			},
			Goals: []structs.Goal{structs.Goal_Wealth},
		},
		Consolidate: {
			MinPeople:    1,
			MaxPeople:    15,
			Probability:  0.1,
			JobPriority:  80,
			Ethos:        structs.Ethos{Caution: modest},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				SCRIBE: 4,
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       3 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      5 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1,
			Goals:         []structs.Goal{structs.Goal_Stability},
		},
		Research: {
			MinPeople:   1,
			MaxPeople:   100,
			Target:      structs.Meta_KeyResearch,
			Probability: 0.3,
			JobPriority: 50,
			Ethos:       structs.Ethos{Piety: -1 * minor, Tradition: -2 * minor},
			ProfessionWeights: map[string]float64{
				MAGE:      5,
				ALCHEMIST: 3,
				PRIEST:    3,
				SCRIBE:    2,
				SCHOLAR:   4,
			},
			Cost: config.Distribution{
				Min:       1000000,  // 100gp
				Max:       20000000, // 2000gp
				Mean:      5000000,  // 500gp
				Deviation: 15000000, // 1500gp
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 7 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       91 * DEFAULT_TICKS_PER_DAY,
				Mean:      45 * DEFAULT_TICKS_PER_DAY,
				Deviation: 60 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1.4,
			Goals:         []structs.Goal{structs.Goal_Research},
		},
		Excommunicate: {
			MinPeople:   1,
			MaxPeople:   100,
			Category:    config.ActionCategoryUnfriendly,
			Target:      structs.Meta_KeyPerson,
			Probability: 0.1,
			Ethos:       structs.Ethos{Piety: major, Caution: -1 * minor},
			Restricted: [][]config.Condition{
				{config.ConditionSrcFactionIsReligion},
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				PRIEST: 5,
			},
			Cost: config.Distribution{},
			TimeToPrepare: config.Distribution{
				Min:       15 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      20 * DEFAULT_TICKS_PER_DAY,
				Deviation: 15 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       15 * DEFAULT_TICKS_PER_DAY,
				Max:       60 * DEFAULT_TICKS_PER_DAY,
				Mean:      45 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight:  1.2,
			Goals:          []structs.Goal{structs.Goal_Piety, structs.Goal_Power},
			TargetMinTrust: (structs.MinEthos * 3) / 4,
			TargetMaxTrust: structs.MaxEthos / 10,
		},
		ConcealSecrets: {
			MinPeople:    1,
			MaxPeople:    100,
			Probability:  0.20,
			JobPriority:  80,
			Ethos:        structs.Ethos{Caution: major},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				SPY:       5,
				ASSASSIN:  3,
				ALCHEMIST: 2,
				CLERK:     2,
				SCRIBE:    2,
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      1 * DEFAULT_TICKS_PER_DAY,
				Deviation: 7 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       3 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      7 * DEFAULT_TICKS_PER_DAY,
				Deviation: 10 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1.5,
			Goals:         []structs.Goal{structs.Goal_Stability, structs.Goal_Espionage},
		},
		GatherSecrets: {
			MinPeople:    2,
			MaxPeople:    14,
			Category:     config.ActionCategoryUnfriendly,
			Probability:  0.10,
			Ethos:        structs.Ethos{Ambition: modest, Caution: 2 * minor},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				SPY:       5,
				ASSASSIN:  3,
				ALCHEMIST: 2,
				CLERK:     2,
				SCRIBE:    2,
			},
			Cost: config.Distribution{
				Min:       1000000,  // 500gp
				Max:       20000000, // 2000gp
				Mean:      2000000,  // 200gp
				Deviation: 15000000, // 1500gp
			},
			TimeToPrepare: config.Distribution{
				Min:       3 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      7 * DEFAULT_TICKS_PER_DAY,
				Deviation: 7 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       60 * DEFAULT_TICKS_PER_DAY,
				Mean:      30 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1.4,
			Goals:         []structs.Goal{structs.Goal_Espionage, structs.Goal_Wealth},
		},
		RevokeLand: {
			MinPeople: 1,
			MaxPeople: 15,
			Target:    structs.Meta_KeyPlot,
			Ethos:     structs.Ethos{Tradition: 2 * minor, Caution: -1 * minor, Altruism: -1 * modest},
			Restricted: [][]config.Condition{
				{config.ConditionSrcFactionIsGovernment},
			},
			TimeToPrepare: config.Distribution{
				Min:       15 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      20 * DEFAULT_TICKS_PER_DAY,
				Deviation: 15 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       30 * DEFAULT_TICKS_PER_DAY,
				Max:       90 * DEFAULT_TICKS_PER_DAY,
				Mean:      60 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			Goals:          []structs.Goal{structs.Goal_Power},
			TargetMinTrust: structs.MinEthos / 2,
			TargetMaxTrust: structs.MaxEthos / 20,
		},
		HireSpies: {
			MinPeople: 1,
			MaxPeople: 5,
			Category:  config.ActionCategoryUnfriendly,
			MercenaryActions: []string{
				Frame,
				Blackmail,
				Assassinate,
				SpreadRumors,
			},
			Probability: 0.02,
			Ethos:       structs.Ethos{Ambition: 2 * minor, Caution: 2 * minor, Pacifism: 1 * minor},
			Cost: config.Distribution{
				Min:       20000000, // 2000gp
				Max:       25000000, // 4000gp
				Mean:      30000000, // 3000gp
				Deviation: 10000000, // 1000gp
			},
			TimeToPrepare: config.Distribution{
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      10 * DEFAULT_TICKS_PER_DAY,
				Deviation: 5 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       30 * DEFAULT_TICKS_PER_DAY,
				Max:       60 * DEFAULT_TICKS_PER_DAY,
				Mean:      45 * DEFAULT_TICKS_PER_DAY,
				Deviation: 15 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:    3,
				SCRIBE:   3,
				SPY:      5,
				SOLDIER:  1,
				ASSASSIN: 2,
				THIEF:    2,
			},
			SecrecyWeight:  1.2,
			Goals:          []structs.Goal{structs.Goal_Military, structs.Goal_Power},
			TargetMinTrust: structs.MinEthos,
			TargetMaxTrust: structs.MinEthos / 2,
		},
		HireMercenaries: {
			MinPeople: 1,
			MaxPeople: 5,
			Category:  config.ActionCategoryUnfriendly,
			MercenaryActions: []string{
				Raid,
				Pillage,
			},
			Probability: 0.02,
			Ethos:       structs.Ethos{Ambition: 3 * minor, Tradition: 1 * minor, Pacifism: 1 * minor},
			Cost: config.Distribution{
				Min:       20000000, // 2000gp
				Max:       25000000, // 4000gp
				Mean:      30000000, // 3000gp
				Deviation: 10000000, // 1000gp
			},
			TimeToPrepare: config.Distribution{
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      10 * DEFAULT_TICKS_PER_DAY,
				Deviation: 5 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       30 * DEFAULT_TICKS_PER_DAY,
				Max:       60 * DEFAULT_TICKS_PER_DAY,
				Mean:      45 * DEFAULT_TICKS_PER_DAY,
				Deviation: 15 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:   3,
				SCRIBE:  2,
				NOBLE:   2,
				SOLDIER: 5,
			},
			SecrecyWeight:  1.2,
			Goals:          []structs.Goal{structs.Goal_Military, structs.Goal_Power},
			TargetMinTrust: structs.MinEthos,
			TargetMaxTrust: structs.MinEthos / 2,
		},
		SpreadRumors: {
			MinPeople:   1,
			MaxPeople:   -1,
			Category:    config.ActionCategoryUnfriendly,
			Probability: 0.20,
			Ethos:       structs.Ethos{Ambition: 3 * minor, Tradition: -1 * minor},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      21 * DEFAULT_TICKS_PER_DAY,
				Deviation: 2 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       90 * DEFAULT_TICKS_PER_DAY,
				Max:       180 * DEFAULT_TICKS_PER_DAY,
				Mean:      135 * DEFAULT_TICKS_PER_DAY,
				Deviation: 70 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				SCRIBE: 4,
				PRIEST: 2,
				NOBLE:  5,
			},
			SecrecyWeight:  1.3,
			Goals:          []structs.Goal{structs.Goal_Power, structs.Goal_Espionage},
			TargetMinTrust: structs.MinEthos / 2,
			TargetMaxTrust: structs.MaxEthos / 2,
		},
		Assassinate: {
			MinPeople:   1,
			MaxPeople:   21,
			Category:    config.ActionCategoryUnfriendly,
			Target:      structs.Meta_KeyPerson,
			Probability: 0.02,
			Ethos:       structs.Ethos{Ambition: modest, Tradition: -1 * minor, Altruism: -1 * modest, Pacifism: -1 * minor},
			Cost: config.Distribution{
				Min:       50000000,  // 5000gp
				Max:       200000000, // 20000gp
				Mean:      60000000,  // 6000gp
				Deviation: 2000000,   // 200gp
			},
			TimeToPrepare: config.Distribution{
				Min:       7,
				Max:       60,
				Mean:      14,
				Deviation: 45,
			},
			TimeToExecute: config.Distribution{
				Min:       1,
				Max:       3,
				Mean:      1,
				Deviation: 1,
			},
			PersonWeight: 0.001,
			ProfessionWeights: map[string]float64{
				ASSASSIN:  8,
				SPY:       5,
				ALCHEMIST: 5,
				MAGE:      4,
				SOLDIER:   4,
			},
			SecrecyWeight:  1.2,
			Goals:          []structs.Goal{structs.Goal_Power, structs.Goal_Espionage},
			TargetMinTrust: (structs.MinEthos * 3) / 4,
			TargetMaxTrust: structs.MinEthos / 4,
		},
		Frame: {
			MinPeople:   1,
			MaxPeople:   15,
			Category:    config.ActionCategoryUnfriendly,
			Target:      structs.Meta_KeyPerson,
			Probability: 0.15,
			Ethos:       structs.Ethos{Ambition: modest, Tradition: -1 * minor, Pacifism: minor},
			Cost: config.Distribution{
				Min:       10000000, // 1000gp
				Max:       25000000, // 2500gp
				Mean:      15000000, // 1500gp
				Deviation: 10000000, // 1000gp
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      21 * DEFAULT_TICKS_PER_DAY,
				Deviation: 10 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       90 * DEFAULT_TICKS_PER_DAY,
				Max:       180 * DEFAULT_TICKS_PER_DAY,
				Mean:      135 * DEFAULT_TICKS_PER_DAY,
				Deviation: 60 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight:  1.4,
			Goals:          []structs.Goal{structs.Goal_Power, structs.Goal_Espionage},
			TargetMinTrust: (structs.MinEthos * 3) / 4,
			TargetMaxTrust: structs.MaxEthos / 20,
		},
		Raid: {
			MinPeople:   14,
			MaxPeople:   120,
			Category:    config.ActionCategoryUnfriendly,
			Target:      structs.Meta_KeyPlot,
			Probability: 0.05,
			Ethos:       structs.Ethos{Ambition: 2 * minor, Pacifism: -1 * modest, Caution: -1 * modest},
			Cost: config.Distribution{
				Min:       10000000, // 1000gp
				Max:       50000000, // 5000gp
				Mean:      25000000, // 2500gp
				Deviation: 25000000, // 2500gp
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				SOLDIER: 5,
				SPY:     2,
				MAGE:    3,
				SAILOR:  1,
				HUNTER:  1,
			},
			RefundOnFailure: 0.2,
			SecrecyWeight:   1.1,
			Goals:           []structs.Goal{structs.Goal_Military, structs.Goal_Wealth},
			TargetMinTrust:  (structs.MinEthos * 3) / 4,
			TargetMaxTrust:  structs.MinEthos / 2,
		},
		Enslave: {
			MinPeople:   14,
			MaxPeople:   60,
			Category:    config.ActionCategoryUnfriendly,
			Target:      structs.Meta_KeyPlot,
			Probability: 0.08,
			Ethos:       structs.Ethos{Altruism: -1 * minor, Caution: -1 * major},
			Cost: config.Distribution{
				Min:       10000000, // 1000gp
				Max:       50000000, // 5000gp
				Mean:      25000000, // 2500gp
				Deviation: 25000000, // 2500gp
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				SOLDIER: 4,
				SPY:     1,
				MAGE:    3,
				SAILOR:  1,
				HUNTER:  2,
			},
			RefundOnFailure: 0.2,
			SecrecyWeight:   1.05,
			Goals:           []structs.Goal{structs.Goal_Growth, structs.Goal_Power},
			TargetMinTrust:  (structs.MinEthos * 3) / 4,
			TargetMaxTrust:  structs.MinEthos / 2,
		},
		Steal: {
			MinPeople:   1,
			MaxPeople:   15,
			Category:    config.ActionCategoryUnfriendly,
			Target:      structs.Meta_KeyPlot,
			Probability: 0.1,
			Ethos:       structs.Ethos{Altruism: -2 * minor, Caution: 2 * minor},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 27 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 2 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				THIEF:  8,
				SPY:    5,
				MAGE:   2,
				SAILOR: 1,
				HUNTER: 1,
			},
			SecrecyWeight:  1.3,
			Goals:          []structs.Goal{structs.Goal_Wealth},
			TargetMinTrust: (structs.MinEthos * 4) / 5,
			TargetMaxTrust: structs.MinEthos / 5,
		},
		Pillage: {
			MinPeople:   14,
			MaxPeople:   120,
			Category:    config.ActionCategoryUnfriendly,
			Target:      structs.Meta_KeyPlot,
			Probability: 0.04,
			Ethos:       structs.Ethos{Ambition: 2 * minor, Pacifism: -1 * modest, Caution: -1 * modest},
			Cost: config.Distribution{
				Min:       10000000, // 1000gp
				Max:       50000000, // 5000gp
				Mean:      25000000, // 2500gp
				Deviation: 25000000, // 2500gp
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 27 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				THIEF:   2,
				SOLDIER: 5,
				SPY:     2,
				MAGE:    5,
				SAILOR:  1,
				HUNTER:  1,
			},
			RefundOnFailure: 0.2,
			SecrecyWeight:   1.1,
			Goals:           []structs.Goal{structs.Goal_Military, structs.Goal_Power},
			TargetMinTrust:  structs.MinEthos,
			TargetMaxTrust:  (structs.MinEthos * 2) / 3,
		},
		Blackmail: {
			MinPeople:   1,
			MaxPeople:   7,
			Category:    config.ActionCategoryUnfriendly,
			Target:      structs.Meta_KeyPerson,
			Probability: 0.10,
			Ethos:       structs.Ethos{Ambition: modest, Caution: 2 * minor, Altruism: -1 * minor},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 6 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 6 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				SCRIBE: 4,
				SPY:    5,
				NOBLE:  4,
				THIEF:  4,
			},
			SecrecyWeight:  2,
			Goals:          []structs.Goal{structs.Goal_Espionage, structs.Goal_Wealth},
			TargetMinTrust: structs.MinEthos / 2,
			TargetMaxTrust: structs.MaxEthos / 2,
		},
		Kidnap: {
			MinPeople:   3,
			MaxPeople:   21,
			Category:    config.ActionCategoryUnfriendly,
			Target:      structs.Meta_KeyPerson,
			Probability: 0.05,
			Ethos:       structs.Ethos{Ambition: modest, Caution: 2 * minor, Altruism: -1 * minor},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 20 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 6 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				THIEF:     5,
				SPY:       4,
				SOLDIER:   5,
				HUNTER:    3,
				ALCHEMIST: 2,
				ASSASSIN:  2,
			},
			SecrecyWeight:  1.2,
			Goals:          []structs.Goal{structs.Goal_Wealth},
			TargetMinTrust: structs.MinEthos / 3,
			TargetMaxTrust: structs.MaxEthos / 20,
		},
		ShadowWar: {
			MinPeople:   100,
			MaxPeople:   -1,
			Category:    config.ActionCategoryHostile,
			Probability: 0.05,
			JobPriority: 90,
			Ethos:       structs.Ethos{Ambition: major, Caution: -1 * minor, Pacifism: -1 * major},
			Restricted: [][]config.Condition{
				{config.ConditionSrcFactionIsCovert},
			},
			Conscription: [][]config.Condition{
				{config.ConditionSrcFactionIsCovert},
			},
			Cost: config.Distribution{
				Min:       20000000,  // 2000gp
				Max:       100000000, // 10000gp
				Mean:      25000000,  // 2500gp
				Deviation: 50000000,  // 5000gp
			},
			TimeToPrepare: config.Distribution{
				Min:       60 * DEFAULT_TICKS_PER_DAY,
				Max:       120 * DEFAULT_TICKS_PER_DAY,
				Mean:      90 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       DEFAULT_TICKS_PER_YEAR * DEFAULT_TICKS_PER_DAY,
				Max:       DEFAULT_TICKS_PER_YEAR * 5 * DEFAULT_TICKS_PER_DAY,
				Mean:      DEFAULT_TICKS_PER_YEAR * 3 * DEFAULT_TICKS_PER_DAY,
				Deviation: DEFAULT_TICKS_PER_YEAR * 4 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				SOLDIER:       4,
				SPY:           5,
				THIEF:         5,
				ASSASSIN:      6,
				MAGE:          5,
				ALCHEMIST:     4,
				SAILOR:        1,
				HUNTER:        1,
				SMITH:         2,
				LEATHERWORKER: 2,
			},
			SecrecyWeight:  0.9,
			Goals:          []structs.Goal{structs.Goal_Military, structs.Goal_Espionage, structs.Goal_Power},
			TargetMinTrust: structs.MinEthos,
			TargetMaxTrust: (structs.MinEthos * 9) / 10,
		},
		Crusade: {
			MinPeople:   500,
			MaxPeople:   -1,
			Category:    config.ActionCategoryHostile,
			Probability: 0.10,
			JobPriority: 100,
			Ethos:       structs.Ethos{Ambition: major, Piety: major, Caution: -1 * major, Pacifism: -1 * major},
			Restricted: [][]config.Condition{
				{config.ConditionSrcFactionIsReligion},
			},
			Conscription: [][]config.Condition{
				{config.ConditionSrcFactionIsCovert, config.ConditionSrcFactionIsReligion},
			},
			Cost: config.Distribution{
				Min:       500000000,  // 50000gp
				Max:       2500000000, // 250000gp
				Mean:      1500000000, // 150000gp
				Deviation: 1000000000, // 100000gp
			},
			TimeToPrepare: config.Distribution{
				Min:       60 * DEFAULT_TICKS_PER_DAY,
				Max:       120 * DEFAULT_TICKS_PER_DAY,
				Mean:      90 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       DEFAULT_TICKS_PER_YEAR * DEFAULT_TICKS_PER_DAY,
				Max:       DEFAULT_TICKS_PER_YEAR * 5 * DEFAULT_TICKS_PER_DAY,
				Mean:      DEFAULT_TICKS_PER_YEAR * 3 * DEFAULT_TICKS_PER_DAY,
				Deviation: DEFAULT_TICKS_PER_YEAR * 4 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				SOLDIER:       5,
				SPY:           2,
				ASSASSIN:      2,
				MAGE:          5,
				ALCHEMIST:     2,
				SAILOR:        1,
				HUNTER:        1,
				SMITH:         4,
				LEATHERWORKER: 2,
			},
			SecrecyWeight:  0.0, // you sort of have to announce this to get the pious to join ..
			Goals:          []structs.Goal{structs.Goal_Military, structs.Goal_Piety, structs.Goal_Power},
			TargetMinTrust: structs.MinEthos,
			TargetMaxTrust: (structs.MinEthos * 9) / 10,
		},
		War: {
			MinPeople:   300,
			MaxPeople:   -1,
			Category:    config.ActionCategoryHostile,
			Probability: 0.10,
			JobPriority: 100,
			Ethos:       structs.Ethos{Ambition: major, Caution: -1 * major, Pacifism: -1 * major},
			Conscription: [][]config.Condition{
				{config.ConditionSrcFactionIsGovernment},
			},
			Cost: config.Distribution{
				Min:       500000000,  // 50000gp
				Max:       2500000000, // 250000gp
				Mean:      1500000000, // 150000gp
				Deviation: 1000000000, // 100000gp
			},
			TimeToPrepare: config.Distribution{
				Min:       60 * DEFAULT_TICKS_PER_DAY,
				Max:       120 * DEFAULT_TICKS_PER_DAY,
				Mean:      90 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       DEFAULT_TICKS_PER_YEAR * DEFAULT_TICKS_PER_DAY,
				Max:       DEFAULT_TICKS_PER_YEAR * 5 * DEFAULT_TICKS_PER_DAY,
				Mean:      DEFAULT_TICKS_PER_YEAR * 3 * DEFAULT_TICKS_PER_DAY,
				Deviation: DEFAULT_TICKS_PER_YEAR * 4 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				SOLDIER:       5,
				SPY:           2,
				ASSASSIN:      2,
				MAGE:          5,
				ALCHEMIST:     2,
				SAILOR:        1,
				HUNTER:        1,
				SMITH:         4,
				LEATHERWORKER: 2,
			},
			SecrecyWeight:  0.1, // very hard to do quietly
			Goals:          []structs.Goal{structs.Goal_Military, structs.Goal_Territory, structs.Goal_Power},
			TargetMinTrust: structs.MinEthos,
			TargetMaxTrust: (structs.MinEthos * 9) / 10,
		},
	}
}
