package faction

// TradeRoute represents a formal treaty to trade a specific good between
// two places.
// TODO: for v1 I think we'll leave this out, but leaving for future expansion.
type TradeRoute struct {
	AreaSource string  // area ID
	AreaTarget string  // area ID
	Commodity  string  // trade good to ship Source -> Target
	Tax        float64 // tax imposed by the importer
	IsIllicit  bool    // the commodity is illegal in one or both areas
}
