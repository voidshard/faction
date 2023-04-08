package premade

/*
Constants here form the basis of our premade configuarations.

These are only examples, your configuration can be as simple or as complex as you like,
The engine doesn't actually care what these are, it only needs to understand how they relate.
*/

const (
	DEFAULT_TICKS_PER_YEAR = 365
	DEFAULT_TICKS_PER_DAY  = 1
)

const (
	// commodities
	IRON_ORE         = "iron-ore"
	IRON_INGOT       = "iron-ingot"
	WHEAT            = "wheat"
	FLOUR_WHEAT      = "flour-wheat"
	STEEL_INGOT      = "steel-ingot"
	HIDE             = "hide"
	MEAT             = "meat"
	LEATHER          = "leather"
	TIMBER           = "timber"
	FLAX             = "flax"  // used for linen / textiles
	LINEN            = "linen" // used for clothing
	FISH             = "fish"
	LINEN_CLOTHING   = "linen-clothing"
	LEATHER_CLOTHING = "leather-clothing"
	IRON_WEAPON      = "iron-weapon"
	IRON_ARMOUR      = "iron-armour"
	STEEL_WEAPON     = "steel-weapon"
	STEEL_ARMOUR     = "steel-armour"
	WOODEN_FURNITURE = "wooden-furniture"
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
	FARMER        = "farmer"
	MINER         = "miner"
	FISHERMAN     = "fisherman"
	HUNTER        = "hunter"
	FORESTER      = "forester"
	WEAVER        = "weaver"
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
	// - more fantasy style professions
	MAGE      = "mage"
	SPY       = "spy"
	ASSASSIN  = "assassin"
	ALCHEMIST = "alchemist"
)
