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
