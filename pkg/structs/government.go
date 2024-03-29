package structs

// Government for our purposes is a set of laws, rights, obligations
// and similar.
//
// This doesn't represent an actual faction, merely the laws / edicts.
// To represent the ruling faction that makes up the Government we use
// a Faction (*gasp* .. I know right).
//
// The ruling faction can add / remove rules from this and all factions
// under it are expected to follow the rules.
//
// Ahem. At least, in theory.
//
// Most factions under a government qualify as Tributaries (they're
// allowed to act on their own, but pay tax to their overlord).
// (See relations.go).
//
// = vassals =
// Note that a Government may have a subordinate faction that is also a
// Government. ie Kingdom of Morovia (Government) run by the Morovian Royal
// Family (Faction) might have a vassal, Kingdom of Duria (Government) run
// by the Durian Royal Family (Faction).
// - Factions under Morovia (but not Duria) obey the laws of Morovia.
// - Factions under Duria obey the laws of Duria.
type Government struct {
	ID string `db:"id"`

	// Every `TaxFrequency` tick(s) the governing faction will collect
	// funds from law abiding factions under it.
	// Covert factions do not pay tax .. since that would require them to exist
	// openly.
	//
	// Higher tax rates make factions increasingly unhappy. Obviously.
	TaxRate      float64 `db:"tax_rate"`
	TaxFrequency int     `db:"tax_frequency"`

	// Outlawed is a set of laws marking factions, actions, commodities as
	// illegal (or not).
	Outlawed *Laws `db:"-"`
}

func (g *Government) ObjectID() string {
	return g.ID
}

type Laws struct {
	// We track these as map[<id>]bool meaning: <id> -> is_illegal
	// Technically we could store that items are specifically legal too, but
	// in practice we regard anything not marked illegal as legal.

	// IllegalFactions are factions that are illegal to be a part of.
	// Any jobs from these factions are illegal, even if the action is otherwise legal,
	// if within the jurisdiction of the government.
	Factions map[string]bool `db:"-"`

	// IllegalActions are actions that are illegal to perform. Any job with such an action
	// is illegal, if within the jurisdiction of the government.
	Actions map[string]bool `db:"-"`

	// IllegalCommodities are commodities that are illegal to trade, harvest or produce
	// if within the jurisdiction of the government.
	Commodities map[string]bool `db:"-"`

	// IllegalResearch is research that is illegal to perform.
	Research map[string]bool `db:"-"`

	// Religions that are outright banned (heretical, generally considered harmful or maybe
	// the ruler(s) are just .. bigots?).
	Religions map[string]bool `db:"-"`
}

// NewLaws returns a new Laws struct with all maps initialised.
func NewLaws() *Laws {
	return &Laws{
		Factions:    map[string]bool{},
		Actions:     map[string]bool{},
		Commodities: map[string]bool{},
		Research:    map[string]bool{},
	}
}
