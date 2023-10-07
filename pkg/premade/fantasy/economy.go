package fantasy

import (
	"log"
	"math"

	"github.com/voidshard/faction/pkg/structs"
)

// Economy is a simple economy that is used for testing and demonstration purposes.
//
// In general a user is expected to provide their own (this is just a janky demo)
type Economy struct {
	commodities      map[string]*structs.Commodity
	baseValues       map[string]float64
	baseLandValue    float64
	landFluctuation  float64
	areaFluctuation  float64
	valueFluctuation float64
	yieldFluctuation float64
}

func (e *Economy) Commodities(area string) []*structs.Commodity {
	var cs []*structs.Commodity
	for _, c := range e.commodities {
		cs = append(cs, c)
	}
	return cs
}

func (e *Economy) IsHarvestable(commodity string) bool {
	comm := e.Commodity(commodity)
	return comm.Requires == nil || len(comm.Requires) == 0
}

func (e *Economy) IsCraftable(commodity string) bool {
	return !e.IsHarvestable(commodity)
}

func (e *Economy) Commodity(name string) *structs.Commodity {
	c, _ := e.commodities[name]
	if c == nil {
		log.Println("warning; no commodity configured in the economy, probably this will blow up", name)
	}
	return c
}

func (e *Economy) CommodityValue(name, areaID string, ticks int) float64 {
	// the addtional Pi here makes Yield & Value opposite each other
	// ie. max yield -> min value, min yield -> max value
	flux := math.Sin(math.Pi+float64(ticks)+floatHash(areaID)) * e.valueFluctuation
	base := e.baseValues[name]
	return base + base*flux
}

func (e *Economy) CommodityYield(name, areaID string, ticks, skill int) float64 {
	flux := math.Sin(float64(ticks)+floatHash(areaID)) * e.valueFluctuation
	base := e.baseValues[name]
	return base + base*flux
}

func (e *Economy) LandValue(areaID string, ticks int) float64 {
	// some areas are just more / less than others irrespective of time
	fh := floatHash(areaID)
	areaWeight := math.Cos(fh) * e.areaFluctuation
	flux := math.Sin(float64(ticks)+fh) * e.landFluctuation
	return e.baseLandValue + e.baseLandValue*flux*areaWeight
}

func NewEconomy() *Economy {
	// www.dandwiki.com/wiki/5e_Trade_Goods
	// 100 cp -> 1 sp
	// 10 sp -> 1 gp
	return &Economy{
		valueFluctuation: 0.15,      // values fluctuate by +/- 15%
		yieldFluctuation: 0.30,      // yields fluctuate by +/- 30%
		landFluctuation:  0.05,      // land values fluctuate by +/- 10%
		areaFluctuation:  2.0,       // some areas are up to 2x more (or less) valuable than others
		baseLandValue:    100000000, // 10000gp per hectare (farmland) / city section (urban)
		// we could add arbitrarily many here, but just for examples
		commodities: map[string]*structs.Commodity{
			FISH: &structs.Commodity{
				Name:       FISH,
				Profession: FISHERMAN,
			},
			OPIUM: &structs.Commodity{
				Name:       OPIUM,
				Profession: FARMER,
			},
			FLAX: &structs.Commodity{
				Name:       FLAX,
				Profession: FARMER,
			},
			LINEN: &structs.Commodity{
				Name:       LINEN,
				Profession: WEAVER,
				Requires:   []map[string]float64{{FLAX: 5.0}},
			},
			LINEN_CLOTHING: &structs.Commodity{
				Name:       LINEN_CLOTHING,
				Profession: CLOTHIER,
				Requires:   []map[string]float64{{LINEN: 1.0}},
			},
			FODDER: &structs.Commodity{
				Name:       FODDER,
				Profession: FARMER,
			},
			WILD_GAME: &structs.Commodity{
				Name:       WILD_GAME,
				Profession: HUNTER,
			},
			HIDE: &structs.Commodity{
				Name:       HIDE,
				Profession: HUNTER,
				Requires:   []map[string]float64{{WILD_GAME: 1.0}},
			},
			MEAT: &structs.Commodity{
				Name:       MEAT,
				Profession: HUNTER,
				Requires:   []map[string]float64{{WILD_GAME: 1.0}},
			},
			LEATHER: &structs.Commodity{
				Name:       LEATHER,
				Profession: TANNER,
				Requires:   []map[string]float64{{HIDE: 1.0}},
			},
			WHEAT: &structs.Commodity{
				// www.theartofdoingstuff.com/im-growing-wheat-this-year-and-you-can-too/
				// 10 plants per m^2 -> ~.5 kg wheat -> 2 cups flour
				Name:       WHEAT,
				Profession: FARMER,
			},
			FLOUR_WHEAT: &structs.Commodity{
				// For every 0.5 units of wheat, we produce 1 unit of flour.
				Name:       FLOUR_WHEAT,
				Profession: FARMER,
				Requires:   []map[string]float64{{WHEAT: 0.5}},
			},
			IRON_ORE: &structs.Commodity{
				Name:       IRON_ORE,
				Profession: MINER,
			},
			IRON_INGOT: &structs.Commodity{
				Name:       IRON_INGOT,
				Profession: SMELTER,
				Requires:   []map[string]float64{{IRON_ORE: 10.0}},
			},
			IRON_TOOLS: &structs.Commodity{
				Name:       IRON_TOOLS,
				Profession: SMITH,
				Requires:   []map[string]float64{{IRON_INGOT: 1.0, TIMBER: 1.0}},
			},
			STEEL_INGOT: &structs.Commodity{
				Name:       STEEL_INGOT,
				Profession: SMELTER,
				Requires:   []map[string]float64{{IRON_INGOT: 2.0}},
			},
			TIMBER: &structs.Commodity{
				Name:       TIMBER,
				Profession: FORESTER,
			},
			STEEL_WEAPON: &structs.Commodity{
				Name:       STEEL_WEAPON,
				Profession: SMITH,
				Requires:   []map[string]float64{{STEEL_INGOT: 2.0, LEATHER: 1.0}},
			},
			STEEL_ARMOUR: &structs.Commodity{
				Name:       STEEL_ARMOUR,
				Profession: SMITH,
				Requires:   []map[string]float64{{STEEL_INGOT: 5.0, LEATHER: 2.0}},
			},
			STEEL_TOOLS: &structs.Commodity{
				Name:       STEEL_TOOLS,
				Profession: SMITH,
				Requires:   []map[string]float64{{STEEL_INGOT: 1.0, TIMBER: 1.0}},
			},
			WOODEN_TOOLS: &structs.Commodity{
				Name:       WOODEN_TOOLS,
				Profession: CARPENTER,
				Requires:   []map[string]float64{{TIMBER: 1.0}},
			},
			WOODEN_FURNITURE: &structs.Commodity{
				Name:       WOODEN_FURNITURE,
				Profession: CARPENTER,
				Requires:   []map[string]float64{{TIMBER: 4.0}},
			},
		},
		baseValues: map[string]float64{ // base value per unit (in copper pieces)
			// nb. these prices need revision. Also "price" here could be a misnomer as it implies
			// use of coinage & a healthy market, which depends on a lot of factors that may not
			// be at play here .. but then this is how our modern brains work so .. eh.
			FISH:             5.0,
			OPIUM:            1000.0, // 1sp
			FLAX:             3.0,
			LINEN:            5000.0, // 5 sp
			LINEN_CLOTHING:   7000.0, // 7 sp
			WILD_GAME:        400.0,  // .4 sp
			HIDE:             10.0,
			MEAT:             1000.0, // 1sp
			TIMBER:           5.0,
			WHEAT:            2.0,
			FLOUR_WHEAT:      5.0,
			IRON_ORE:         1000.0,    // 1sp
			IRON_INGOT:       200000.0,  // 20gp
			STEEL_INGOT:      1000000.0, // 100gp
			LEATHER:          30.0,
			IRON_TOOLS:       100000.0,  // 10gp
			STEEL_TOOLS:      500000.0,  // 50gp
			STEEL_WEAPON:     2000000.0, // 200gp
			STEEL_ARMOUR:     5000000.0, // 500gp
			WOODEN_TOOLS:     6.0,
			WOODEN_FURNITURE: 5000.0, // 5sp
		},
	}
}

// floatHash dumbly converts a string into a float64 deterministically.
// We use this to change results of valuations based on area ID.
func floatHash(s string) float64 {
	f := 0.0
	for i, c := range s {
		f += float64(c) * float64(i) * 0.5
	}
	return f
}
