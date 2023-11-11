package fantasy

/*
Constants here form the basis of our premade configuarations.

These are only examples, your configuration can be as simple or as complex as you like,
The engine doesn't actually care what these are, it only needs to understand how they relate.
*/

const (
	DEFAULT_TICKS_PER_YEAR = 365
	DEFAULT_TICKS_PER_DAY  = 1
)

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
	// Friendly actions (some of these target another faction)
	Trade       = "trade"       // trade goods with another faction, everyone wins
	Bribe       = "bribe"       // pay a faction to increase their favor
	Festival    = "festival"    // hold a festival (usually religious in nature), increases general favor & affiliation
	Ritual      = "ritual"      // hold a public ritual, increases faith
	RequestLand = "requestland" // request the government grant the use of some land
	Charity     = "charity"     // donate money to people of an area(s), increases general favor
	// Neutral actions (these target "self")
	Propoganda     = "propoganda"      // increases general favor towards you, cheaper than 'Charity' but less um, honest
	Recruit        = "recruit"         // push recuitment (overtly or covertly), increases general affiliation
	Expand         = "expand"          // purchase more land / property, add base of operations in another city etc
	Craft          = "craft"           // focus on crafting, increasing funds
	Harvest        = "harvest"         // focus on harvesting, increasing crop yield, mining - increasing funds
	Consolidate    = "consolidate"     // consolidate; internal re-organisation, process streamlining etc (increases cohession)
	Research       = "research"        // adds some science & funds
	Excommunicate  = "excommunicate"   // religion explicitly excommunicates someone; person loses favour & affiliation
	ConcealSecrets = "conceal-secrets" // bribes are paid, people are silenced, documents burned (+secrecy)
	// Unfriendly actions (all of these have a target faction)
	GatherSecrets   = "gather-secrets"   // attempt to discover secrets of target faction
	RevokeLand      = "revoke-land"      // government retracts right to use some land
	HireMercenaries = "hire-mercenaries" // hire another faction to do something for you
	HireSpies       = "hire-spies"       // hire another faction to spy on another faction
	SpreadRumors    = "spread-rumors"    // inverse of Propoganda; decrease general favor toward target
	Assassinate     = "assassinate"      // someone is selected for elimination
	Frame           = "frame"            // the non religious version of 'excommunicate,' but can involve legal .. entanglements
	Raid            = "raid"             // small armed conflict with the aim of destroying as much as possible
	Enslave         = "enslave"          // similar to raid, but with the aim of capturing people
	Steal           = "steal"            // similar to raid, but with less stabbing
	Pillage         = "pillage"          // small armed conflict with the aim of stealing wealth
	Blackmail       = "blackmail"        // trade in a secret for a pile of gold, or ruin a reputation
	Kidnap          = "kidnap"           // similar to blackmail, but you trade back a person
	// Hostile actions
	ShadowWar = "shadow-war" // shadow war is a full armed conflict carried out between one or more covert factions
	Crusade   = "crusade"    // excommunication on a grander and more permanent scale
	War       = "war"        // full armed conflict
)

const (
	// commodities
	IRON_ORE         = "iron-ore"
	IRON_INGOT       = "iron-ingot" // or "bloom" I guess
	WHEAT            = "wheat"
	FLOUR_WHEAT      = "flour-wheat"
	STEEL_INGOT      = "steel-ingot"
	WILD_GAME        = "wild-game" // produces meat & hide
	FODDER           = "fodder"    // fodder "grass" suitable for grazing animals
	HIDE             = "hide"
	MEAT             = "meat"
	LEATHER          = "leather"
	TIMBER           = "timber" // place holder for wood / charcoal (which .. really aren't the same)
	FLAX             = "flax"   // used for linen / textiles
	LINEN            = "linen"  // used for clothing
	FISH             = "fish"
	LINEN_CLOTHING   = "linen-clothing"
	IRON_TOOLS       = "iron-tools"
	STEEL_WEAPON     = "steel-weapon" // nb. if we have steel, probably iron is sub-standard for weapons / armour
	STEEL_ARMOUR     = "steel-armour"
	STEEL_TOOLS      = "steel-tools"
	WOODEN_FURNITURE = "wooden-furniture"
	WOODEN_TOOLS     = "wooden-tools"
	OPIUM            = "opium"
)

const (
	// Research topics
	// - ancient
	AGRICULTURE = "agriculture"
	WARFARE     = "warfare"    // starting with Sargon of Akkad (c. 2334 BC)
	ASTRONOMY   = "astronomy"  // Greek, Babylonian, Egyptian ..
	METALLURGY  = "metallurgy" // starting with copper -> bronze -> iron -> steel

	// - classical
	PHILOSOPHY   = "philosophy"   // plato, aristotle, socrates
	MEDICINE     = "medicine"     // hippocrates, galen
	MATHEMATICS  = "mathematics"  // euclid, pythagoras
	LITERATURE   = "literature"   // homer, virgil
	LAW          = "law"          // plato, aristotle, socrates
	ARCHITECTURE = "architecture" // ionic / doric schools of greek architecture

	// - medieval
	THEOLOGY    = "theology" // ie. the formal study of theology distinct from philosophy / metaphysics
	PHYSICS     = "physics"
	ENGINEERING = "engineering"

	// - renessaince
	HISTORY   = "history"
	ECONOMICS = "economics"

	// - industrial
	CHEMISTRY = "chemistry"
	BIOLOGY   = "biology"

	// - fantasy
	MAGIC_ARCANA = "magic-arcana" // acceptable magic
	MAGIC_OCCULT = "magic-occult" // forbidden magic
	ALCHEMY      = "alchemy"
)

const (
	// Professions
	// - people who craft / harvest
	// Nb. for our simple example we're vastly simplifying / collapsing what was a variety of professions.
	// There's a lot more complexity here in a closer-to-real-world implementation.
	FARMER        = "farmer"
	MINER         = "miner"
	FISHERMAN     = "fisherman"
	HUNTER        = "hunter"
	FORESTER      = "forester" // + someone who makes charcoal
	WEAVER        = "weaver"   // + spinner
	CLOTHIER      = "clothier"
	TANNER        = "tanner"
	LEATHERWORKER = "leatherworker"
	CARPENTER     = "carpenter"
	SMELTER       = "smelter"
	SMITH         = "smith"
	// - people with specialized professions
	SAILOR   = "sailor"
	SOLDIER  = "soldier"
	CLERK    = "clerk"
	PRIEST   = "priest"
	SCRIBE   = "scribe"
	MERCHANT = "merchant"
	THIEF    = "thief"
	SCHOLAR  = "scholar"
	// - medieval
	NOBLE = "noble"
	// - more fantasy style professions (well, some of these were actually a professions I guess ..)
	MAGE      = "mage"
	SPY       = "spy"
	ASSASSIN  = "assassin"
	ALCHEMIST = "alchemist"
)

const (
	// Biomes
	// It's generally considered that biomes are a mixture of biomes but are dominated by one.
	WOODLAND = "woodland"
	FOREST   = "forest"
	JUNGLE   = "jungle"

	// biomes by altitude
	ALPINE       = "alpine" // high mountians, probably permanently snow capped
	MOUNTAINS    = "mountains"
	HILLS        = "hills"
	SUBTERRANIAN = "subterranean"

	// plains of various types depending on moisture / grass vs shrubs vs trees
	SAVANNAH  = "savannah"
	GRASSLAND = "grassland" // savannah, with few if any woody plants
	BADLANDS  = "badlands"  // land where all the good soil, clay and sedimentary rock is gone
	DESERT    = "desert"    // minimal to no water

	// water laden biomes
	LAKE  = "lake"
	SWAMP = "swamp"
	MARSH = "marsh"
	COAST = "coast"
	REEF  = "reef"

	// underwater biomes
	// www.whoi.edu/know-your-ocean/ocean-topics/how-the-ocean-works/ocean-zones/
	SUBAQUATIC_SUNLIT   = "subaquatic-sunlit"   // [0-200m] close enough to surface to get light
	SUBAQUATIC_TWILIGHT = "subaquatic-twilight" // [200m-1km] not much light
	SUBAQUATIC_MIDNIGHT = "subaquatic-midnight" // [1-3km] no light, cold, high pressure
	SUBAQUATIC_ABYSSAL  = "subaquatic-abyssal"  // [3-6km] no light, very cold, high pressures
	SUBAQUATIC_HADIC    = "subaquatic-hadic"    // [6km+] no light, freezing, extreme pressures
	SUBAQUATIC_VOLCANIC = "subaquatic-volcanic" // underwater volcanic vents

	// cold biomes
	SUBARCTIC = "subarctic" // think Inuits, cold but there's hunting & fishing
	TUNDRA    = "tundra"    // frozen grassland, probably with permafrost, no trees, maybe spare grass
	ARCTIC    = "arctic"    // penguins, ice, permafrost, no plants

	// volcanic regions
	VOLCANIC_FERTILE = "volcanic-fertile" // volcanic-origin soil that is rich in nutrients
	VOLCANIC_BARREN  = "volcanic-barren"  // area so steeped in volcanic activity that nothing grows
)
