package simutil

import (
	"github.com/voidshard/faction/pkg/economy"
	"github.com/voidshard/faction/pkg/structs"
)

func PlotValuation(p *structs.Plot, eco economy.Economy, tickOffset int) float64 {
	land := float64(p.Size) * eco.LandValue(p.AreaID, tickOffset)
	if p.Commodity == "" {
		return land
	}
	return land + float64(p.Yield)*eco.CommodityValue(p.Commodity, p.AreaID, tickOffset)
}
