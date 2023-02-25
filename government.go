package faction

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
//
// = vassals =
// Note that a Government may have a subordinate faction that is also a
// Government. ie Kingdom of Morovia (Government) run by the Morovian Royal
// Family (Faction) might have a vassal, Kingdom of Duria (Government) (run
// by the Durian Royal Family).
// - Factions under Morovia (but not Duria) obey the laws of Morovia.
// - Factions under Duria obey the laws of Duria.
type Government struct {

	// IllegalActions are actions this government considers illegal (obviously).
	// Factions that perform these may face punishment if they're caught.
	//
	// Any faction performing an illegal action automatically attempts to do
	// so secretly (unless it's in full revolt ..).
	IllegalActions []ActionType

	// IllegalCrafts are Commodities (or even processes) that the government
	// has outlawed.
	//
	// It is considered illegal to craft, import or export these crafts, however
	// they can still be very lucrative to people less concerned with the law.
	IllegalCrafts []string

	// Every tick the governing faction will collect funds from law abiding
	// factions under it.
	//
	// Higher tax rates make factions increasingly unhappy. Obviously.
	TaxRate float64

	// LandRights are (Area, Commodity) pairs that the government has the right to.
	// Ie.
	//   Area FooTown - Iron Ore Mine
	//   Area BobVillage - Wheat Field
	// The government doesn't work these itself but sells them off to faction(s)
	// to work them (or spawns divisions of the government to handle them).
	// Well, technically a government *could* work them .. if they awarded
	// themselves the right. This would be like .. a state run mine or something.
	LandRights []*LandRight
}

// LandRight is the legal right to work some land for some resource
type LandRight struct {
	// Area ID
	Area string

	// Resource (commodity name)
	Resource string
}
