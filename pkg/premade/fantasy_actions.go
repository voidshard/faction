package premade

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// ActionsFantasy returns a default action configuration for fantasy type settings.
//
// This could be useful as a starting point for your own configuration.
func ActionsFantasy() map[structs.ActionType]*config.Action {
	return map[structs.ActionType]*config.Action{
		structs.ActionTypeTrade: {
			MinPeople: 3,
			MaxPeople: 150,
			Ethos:     structs.Ethos{Altruism: -10},
			Cost: config.Distribution{
				Min:       1000000,  // 100gp
				Max:       50000000, // 5000gp
				Mean:      10000000, // 1000gp
				Deviation: 500000,   // 50gp
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       5 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       10 * DEFAULT_TICKS_PER_DAY,
				Max:       90 * DEFAULT_TICKS_PER_DAY,
				Mean:      30 * DEFAULT_TICKS_PER_DAY,
				Deviation: 5 * DEFAULT_TICKS_PER_DAY,
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
			Ethos:     structs.Ethos{Tradition: -20, Ambition: 20, Caution: 20},
			Cost: config.Distribution{
				Min:       20000000,  // 2000gp
				Max:       100000000, // 10000gp
				Mean:      25000000,  // 2500gp
				Deviation: 500000,    // 50gp
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      7 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			ProfessionWeights: map[string]float64{
				CLERK:    1,
				SCRIBE:   1,
				SPY:      5,
				ASSASSIN: 2,
				NOBLE:    2,
			},
			SecrecyWeight:      1.5,
			IllegalProbability: 0.6,
		},
		structs.ActionTypeFestival: {
			MinPeople: 100,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Altruism: 60, Tradition: 10, Piety: 40},
			Cost: config.Distribution{
				Min:       1000000,  // 100gp
				Max:       10000000, // 1000gp
				Mean:      5000000,  // 500gp
				Deviation: 500000,   // 50gp
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      21 * DEFAULT_TICKS_PER_DAY,
				Deviation: 3 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       3 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
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
			Ethos:     structs.Ethos{Altruism: 60, Tradition: 20},
			TimeToPrepare: config.Distribution{
				Min:       15 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      20 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       30 * DEFAULT_TICKS_PER_DAY,
				Max:       90 * DEFAULT_TICKS_PER_DAY,
				Mean:      60 * DEFAULT_TICKS_PER_DAY,
				Deviation: 3 * DEFAULT_TICKS_PER_DAY,
			},
		},
		structs.ActionTypeCharity: {
			MinPeople: 5,
			MaxPeople: 1500,
			Ethos:     structs.Ethos{Altruism: 100, Piety: 20},
			Cost: config.Distribution{
				Min:       100000,   // 10gp
				Max:       10000000, // 1000gp
				Mean:      5000000,  // 500gp
				Deviation: 10000,    // 10gp
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      21 * DEFAULT_TICKS_PER_DAY,
				Deviation: 2 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      10 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 2,
		},
		structs.ActionTypePropoganda: {
			MinPeople: 30,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Tradition: -10, Ambition: 30},
			Cost: config.Distribution{
				Min:       50000000,  // 5000gp
				Max:       100000000, // 10000gp
				Mean:      75000000,  // 7500gp
				Deviation: 1000000,   // 100gp
			},
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
				Deviation: 3 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				SCRIBE: 4,
				PRIEST: 2,
				NOBLE:  5,
			},
			SecrecyWeight:      1.1,
			IllegalProbability: 0.1,
		},
		structs.ActionTypeRecruit: {
			MinPeople: 1,
			MaxPeople: 60,
			Ethos:     structs.Ethos{Caution: 40},
			Cost: config.Distribution{
				Min:       100000,  // 10gp
				Max:       1000000, // 100gp
				Mean:      10000,   // 1gp
				Deviation: 10000,   // 10gp
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				SCRIBE: 4,
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       30 * DEFAULT_TICKS_PER_DAY,
				Max:       90 * DEFAULT_TICKS_PER_DAY,
				Mean:      60 * DEFAULT_TICKS_PER_DAY,
				Deviation: 3 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1,
		},
		structs.ActionTypeExpand: {
			MinPeople:    1,
			MaxPeople:    15,
			Ethos:        structs.Ethos{Caution: 30},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				NOBLE:  5,
				SCRIBE: 4,
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1,
		},
		structs.ActionTypeDownsize: {
			MinPeople:    1,
			MaxPeople:    15,
			Ethos:        structs.Ethos{Caution: -40},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				NOBLE:  5,
				SCRIBE: 4,
			},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1,
		},
		structs.ActionTypeCraft: {
			MinPeople:    1,
			MaxPeople:    -1,
			Ethos:        structs.Ethos{Altruism: -15},
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
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
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
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
		},
		structs.ActionTypeConsolidate: {
			MinPeople:    1,
			MaxPeople:    15,
			Ethos:        structs.Ethos{Caution: 40},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				SCRIBE: 4,
			},
			TimeToPrepare: config.Distribution{
				Min:       7 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight: 1,
		},
		structs.ActionTypeResearch: {
			MinPeople: 1,
			MaxPeople: 100,
			Ethos:     structs.Ethos{Piety: -10},
			ProfessionWeights: map[string]float64{
				MAGE:      5,
				ALCHEMIST: 3,
				PRIEST:    3,
				SCRIBE:    2,
			},
			Cost: config.Distribution{
				Min:       5000000,  // 500gp
				Max:       20000000, // 2000gp
				Mean:      1000000,  // 100gp
				Deviation: 100000,   // 10gp
			},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      1 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       21 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight:      1.4,
			IllegalProbability: 0.01,
		},
		structs.ActionTypeExcommunicate: {
			MinPeople:    1,
			MaxPeople:    100,
			Ethos:        structs.Ethos{Piety: 75, Caution: -10},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				PRIEST: 5,
			},
			Cost: config.Distribution{},
			TimeToPrepare: config.Distribution{
				Min:       15 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      20 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       15 * DEFAULT_TICKS_PER_DAY,
				Max:       60 * DEFAULT_TICKS_PER_DAY,
				Mean:      45 * DEFAULT_TICKS_PER_DAY,
				Deviation: 2 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight:      1.2,
			IllegalProbability: 0.05,
		},
		structs.ActionTypeConcealSecrets: {
			MinPeople:    1,
			MaxPeople:    100,
			Ethos:        structs.Ethos{Caution: 90},
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
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       3 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      7 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight:      1.5,
			IllegalProbability: 0.05,
		},
		structs.ActionTypeGatherSecrets: {
			MinPeople:    2,
			MaxPeople:    14,
			Ethos:        structs.Ethos{Ambition: 50},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				SPY:       5,
				ASSASSIN:  3,
				ALCHEMIST: 2,
				CLERK:     2,
				SCRIBE:    2,
			},
			Cost: config.Distribution{
				Min:       5000000,  // 500gp
				Max:       20000000, // 2000gp
				Mean:      1000000,  // 100gp
				Deviation: 100000,   // 10gp
			},
			TimeToPrepare: config.Distribution{
				Min:       3 * DEFAULT_TICKS_PER_DAY,
				Max:       14 * DEFAULT_TICKS_PER_DAY,
				Mean:      7 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       60 * DEFAULT_TICKS_PER_DAY,
				Mean:      30 * DEFAULT_TICKS_PER_DAY,
				Deviation: 2 * DEFAULT_TICKS_PER_DAY,
			},
			SecrecyWeight:      1.4,
			IllegalProbability: 0.2,
		},
		structs.ActionTypeRevokeLand: {
			MinPeople: 1,
			MaxPeople: 15,
			Ethos:     structs.Ethos{Tradition: 20, Caution: -10, Altruism: -50},
			TimeToPrepare: config.Distribution{
				Min:       15 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      20 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       30 * DEFAULT_TICKS_PER_DAY,
				Max:       90 * DEFAULT_TICKS_PER_DAY,
				Mean:      60 * DEFAULT_TICKS_PER_DAY,
				Deviation: 3 * DEFAULT_TICKS_PER_DAY,
			},
		},
		structs.ActionTypeSpreadRumors: {
			MinPeople: 1,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Ambition: 30, Tradition: -10},
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
				Deviation: 3 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				SCRIBE: 4,
				PRIEST: 2,
				NOBLE:  5,
			},
			SecrecyWeight:      1.3,
			IllegalProbability: 0.2,
		},
		structs.ActionTypeAssassinate: {
			MinPeople: 1,
			MaxPeople: 21,
			Ethos:     structs.Ethos{Ambition: 60, Tradition: -10, Altruism: -50, Pacifism: -50},
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
				Deviation: 2,
			},
			TimeToExecute: config.Distribution{
				Min:       1,
				Max:       3,
				Mean:      1,
				Deviation: 1,
			},
			PersonWeight: 0.001,
			ProfessionWeights: map[string]float64{
				ASSASSIN:  0.1,
				SPY:       0.06,
				ALCHEMIST: 0.05,
				MAGE:      0.04,
				SOLDIER:   0.02,
			},
			SecrecyWeight:      1.2,
			IllegalProbability: 0.95,
		},
		structs.ActionTypeFrame: {
			MinPeople: 1,
			MaxPeople: 15,
			Ethos:     structs.Ethos{Ambition: 40, Tradition: -10, Pacifism: -20},
			Cost: config.Distribution{
				Min:       10000000, // 1000gp
				Max:       25000000, // 2500gp
				Mean:      15000000, // 1500gp
				Deviation: 500000,   // 50gp
			},
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
				Deviation: 3 * DEFAULT_TICKS_PER_DAY,
			},
			IllegalProbability: 0.96,
			SecrecyWeight:      1.4,
		},
		structs.ActionTypeRaid: {
			MinPeople: 14,
			MaxPeople: 120,
			Ethos:     structs.Ethos{Ambition: 20, Pacifism: -70, Caution: -50},
			Cost: config.Distribution{
				Min:       10000000, // 1000gp
				Max:       50000000, // 5000gp
				Mean:      25000000, // 2500gp
				Deviation: 500000,   // 50gp
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
			IllegalProbability: 0.94,
			RefundOnCancel:     0.5,
			RefundOnFailure:    0.2,
			SecrecyWeight:      1.1,
		},
		structs.ActionTypePillage: {
			MinPeople: 14,
			MaxPeople: 120,
			Ethos:     structs.Ethos{Ambition: 20, Pacifism: -70, Caution: -50},
			Cost: config.Distribution{
				Min:       10000000, // 1000gp
				Max:       50000000, // 5000gp
				Mean:      25000000, // 2500gp
				Deviation: 500000,   // 50gp
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
				MAGE:    5,
				SAILOR:  1,
				HUNTER:  1,
			},
			IllegalProbability: 0.98,
			RefundOnCancel:     0.5,
			RefundOnFailure:    0.2,
			SecrecyWeight:      1.1,
		},
		structs.ActionTypeBlackmail: {
			MinPeople: 1,
			MaxPeople: 7,
			Ethos:     structs.Ethos{Ambition: 30, Caution: 20},
			TimeToPrepare: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				CLERK:  2,
				SCRIBE: 4,
				SPY:    5,
			},
			IllegalProbability: 0.9,
			SecrecyWeight:      2,
		},
		structs.ActionTypeKidnap: {
			MinPeople: 3,
			MaxPeople: 21,
			Ethos:     structs.Ethos{Ambition: 30, Caution: 20},
			TimeToPrepare: config.Distribution{
				Min:       14 * DEFAULT_TICKS_PER_DAY,
				Max:       30 * DEFAULT_TICKS_PER_DAY,
				Mean:      3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       1 * DEFAULT_TICKS_PER_DAY,
				Max:       7 * DEFAULT_TICKS_PER_DAY,
				Mean:      2 * DEFAULT_TICKS_PER_DAY,
				Deviation: 1 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				SPY:     5,
				SOLDIER: 5,
			},
			IllegalProbability: 0.95,
			SecrecyWeight:      1.2,
		},
		structs.ActionTypeShadowWar: {
			MinPeople: 100,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Ambition: 80, Caution: 80, Pacifism: -100},
			Cost: config.Distribution{
				Min:       20000000,  // 2000gp
				Max:       100000000, // 10000gp
				Mean:      25000000,  // 2500gp
				Deviation: 500000,    // 50gp
			},
			TimeToPrepare: config.Distribution{
				Min:       60 * DEFAULT_TICKS_PER_DAY,
				Max:       120 * DEFAULT_TICKS_PER_DAY,
				Mean:      90 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       365 * DEFAULT_TICKS_PER_DAY,
				Max:       365 * 5 * DEFAULT_TICKS_PER_DAY,
				Mean:      365 * 3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			PersonWeight: 1,
			ProfessionWeights: map[string]float64{
				SOLDIER:       4,
				SPY:           5,
				ASSASSIN:      5,
				MAGE:          5,
				ALCHEMIST:     4,
				SAILOR:        1,
				HUNTER:        1,
				SMITH:         2,
				LEATHERWORKER: 2,
			},
			IllegalProbability: 0.98,
			SecrecyWeight:      1,
		},
		structs.ActionTypeCrusade: {
			MinPeople: 500,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Ambition: 20, Piety: 100, Caution: -10, Pacifism: -100},
			Cost: config.Distribution{
				Min:       50000000,  // 5000gp
				Max:       250000000, // 25000gp
				Mean:      150000000, // 15000gp
				Deviation: 1000000,   // 100gp
			},
			TimeToPrepare: config.Distribution{
				Min:       60 * DEFAULT_TICKS_PER_DAY,
				Max:       120 * DEFAULT_TICKS_PER_DAY,
				Mean:      90 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       365 * DEFAULT_TICKS_PER_DAY,
				Max:       365 * 5 * DEFAULT_TICKS_PER_DAY,
				Mean:      365 * 3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
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
			IllegalProbability: 0.7,
		},
		structs.ActionTypeWar: {
			MinPeople: 500,
			MaxPeople: -1,
			Ethos:     structs.Ethos{Ambition: 80, Caution: -10, Pacifism: -100},
			Cost: config.Distribution{
				Min:       100000000, // 10000gp
				Max:       400000000, // 40000gp
				Mean:      200000000, // 20000gp
				Deviation: 1000000,   // 100gp
			},
			TimeToPrepare: config.Distribution{
				Min:       60 * DEFAULT_TICKS_PER_DAY,
				Max:       120 * DEFAULT_TICKS_PER_DAY,
				Mean:      90 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
			},
			TimeToExecute: config.Distribution{
				Min:       365 * DEFAULT_TICKS_PER_DAY,
				Max:       365 * 5 * DEFAULT_TICKS_PER_DAY,
				Mean:      365 * 3 * DEFAULT_TICKS_PER_DAY,
				Deviation: 30 * DEFAULT_TICKS_PER_DAY,
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
			IllegalProbability: 0.98,
		},
	}
}
