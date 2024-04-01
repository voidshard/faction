package simutil

import (
	"github.com/voidshard/faction/pkg/economy"
	"github.com/voidshard/faction/pkg/structs"
)

func PlotValuation(p *structs.Plot, eco economy.Economy, tickOffset int) float64 {
	if p.Crop == nil {
		return 0
	}
	land := float64(p.Crop.Size) * eco.LandValue(p.AreaID, tickOffset)
	if p.Crop.Commodity == "" {
		return land
	}
	return land + float64(p.Crop.Yield)*eco.CommodityValue(p.Crop.Commodity, p.AreaID, tickOffset)
}
