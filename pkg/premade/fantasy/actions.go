package premade

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// Actions returns a default action configuration for fantasy type settings.
//
// This could be useful as a starting point for your own configuration.
func Actions() map[structs.ActionType]*config.Action {
	minor := structs.MaxEthos / 10
	modest := structs.MaxEthos / 2
	major := structs.MaxEthos

	return map[structs.ActionType]*config.Action{
		structs.ActionTypeTrade: {
			MinPeople: 3,
			MaxPeople: 150,
			Ethos:     structs.Ethos{Altruism: -1 * minor},
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
			RefundOnCancel: 0.5,
			SecrecyWeight:  1,
		},
		structs.ActionTypeBribe: {
			MinPeople: 1,
			MaxPeople: 15,
			Ethos:     structs.Ethos{Tradition: -2 * minor, Ambition: 2 * minor, Caution: 2 * minor},
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
			SecrecyWeight: 1.5,
		},
		structs.ActionTypeFestival: {
			MinPeople: 100,
			MaxPeople: -1,
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
		},
		structs.ActionTypeGrantLand: {
			MinPeople: 1,
			MaxPeople: 15,
			Ethos:     structs.Ethos{Altruism: modest, Tradition: 2 * minor},
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
		},
		structs.ActionTypeCharity: {
			MinPeople: 5,
			MaxPeople: 1500,
			Ethos:     structs.Ethos{Altruism: major, Piety: minor},
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
		},
		structs.ActionTypePropoganda: {
			MinPeople: 30,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Tradition: -1 * minor, Ambition: modest},
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
			SecrecyWeight: 1.1,
		},
		structs.ActionTypeRecruit: {
			MinPeople:    1,
			MaxPeople:    60,
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
		},
		structs.ActionTypeExpand: { // cost determined by land prices in economy
			MinPeople:    1,
			MaxPeople:    15,
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
		},
		structs.ActionTypeDownsize: {
			MinPeople:    1,
			MaxPeople:    15,
			Ethos:        structs.Ethos{Caution: -1 * modest},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				NOBLE:  5,
				SCRIBE: 4,
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      1 * DEFAULT_TICKS_PER_DAY,
				Deviation: 2 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      7 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1,
		},
		structs.ActionTypeCraft: {
			MinPeople:    1,
			MaxPeople:    -1,
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
		},
		structs.ActionTypeHarvest: {
			MinPeople:    1,
			MaxPeople:    -1,
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
		},
		structs.ActionTypeConsolidate: {
			MinPeople:    1,
			MaxPeople:    15,
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
		},
		structs.ActionTypeResearch: {
			MinPeople: 1,
			MaxPeople: 100,
			Ethos:     structs.Ethos{Piety: -1 * minor, Tradition: -2 * minor},
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
		},
		structs.ActionTypeExcommunicate: {
			MinPeople:    1,
			MaxPeople:    100,
			Ethos:        structs.Ethos{Piety: major, Caution: -1 * minor},
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
			SecrecyWeight: 1.2,
		},
		structs.ActionTypeConcealSecrets: {
			MinPeople:    1,
			MaxPeople:    100,
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
		},
		structs.ActionTypeGatherSecrets: {
			MinPeople:    2,
			MaxPeople:    14,
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
		},
		structs.ActionTypeRevokeLand: {
			MinPeople: 1,
			MaxPeople: 15,
			Ethos:     structs.Ethos{Tradition: 2 * minor, Caution: -1 * minor, Altruism: -1 * modest},
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
		},
		structs.ActionTypeHireMercenaries: {
			MinPeople: 1,
			MaxPeople: 5,
			Ethos:     structs.Ethos{Ambition: 3 * minor, Tradition: 1 * minor, Pacifism: 1 * minor},
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
				CLERK:   4,
				SCRIBE:  2,
				NOBLE:   2,
				SOLDIER: 5,
			},
			SecrecyWeight: 1.2,
		},
		structs.ActionTypeSpreadRumors: {
			MinPeople: 1,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Ambition: 3 * minor, Tradition: -1 * minor},
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
			SecrecyWeight: 1.3,
		},

		structs.ActionTypeAssassinate: {
			MinPeople: 1,
			MaxPeople: 21,
			Ethos:     structs.Ethos{Ambition: modest, Tradition: -1 * minor, Altruism: -1 * modest, Pacifism: -1 * minor},
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
			SecrecyWeight: 1.2,
		},
		structs.ActionTypeFrame: {
			MinPeople: 1,
			MaxPeople: 15,
			Ethos:     structs.Ethos{Ambition: modest, Tradition: -1 * minor, Pacifism: minor},
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
			SecrecyWeight: 1.4,
		},
		structs.ActionTypeRaid: {
			MinPeople: 14,
			MaxPeople: 120,
			Ethos:     structs.Ethos{Ambition: 2 * minor, Pacifism: -1 * modest, Caution: -1 * modest},
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
			RefundOnCancel:  0.5,
			RefundOnFailure: 0.2,
			SecrecyWeight:   1.1,
		},
		structs.ActionTypeSteal: {
			MinPeople: 1,
			MaxPeople: 15,
			Ethos:     structs.Ethos{Altruism: -2 * minor, Caution: 2 * minor},
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
			SecrecyWeight: 1.3,
		},
		structs.ActionTypePillage: {
			MinPeople: 14,
			MaxPeople: 120,
			Ethos:     structs.Ethos{Ambition: 2 * minor, Pacifism: -1 * modest, Caution: -1 * modest},
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
			RefundOnCancel:  0.5,
			RefundOnFailure: 0.2,
			SecrecyWeight:   1.1,
		},
		structs.ActionTypeBlackmail: {
			MinPeople: 1,
			MaxPeople: 7,
			Ethos:     structs.Ethos{Ambition: modest, Caution: 2 * minor, Altruism: -1 * minor},
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
			SecrecyWeight: 2,
		},
		structs.ActionTypeKidnap: {
			MinPeople: 3,
			MaxPeople: 21,
			Ethos:     structs.Ethos{Ambition: modest, Caution: 2 * minor, Altruism: -1 * minor},
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
			SecrecyWeight: 1.2,
		},
		structs.ActionTypeShadowWar: {
			MinPeople: 100,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Ambition: major, Caution: -1 * minor, Pacifism: -1 * major},
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
			SecrecyWeight: 1,
		},
		structs.ActionTypeCrusade: {
			MinPeople: 500,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Ambition: major, Piety: major, Caution: -1 * major, Pacifism: -1 * major},
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
		},
		structs.ActionTypeWar: {
			MinPeople: 500,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Ambition: major, Caution: -1 * major, Pacifism: -1 * major},
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
		},
	}
}
