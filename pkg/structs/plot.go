package structs

func (p *Plot) ObjectID() string {
	return p.ID
}

func (ls *LandSummary) Add(p *Plot) {
	ls.addCommodity(p)

	areaSum, ok := ls.Areas[p.AreaID]
	if !ok {
		areaSum = NewLandSummary()
	}
	areaSum.addCommodity(p)
	ls.Areas[p.AreaID] = areaSum
}

func (ls *LandSummary) addCommodity(p *Plot) {
	if p.Crop == nil {
		return
	}

	crop := ls.Commodities[p.Crop.Commodity]
	if crop == nil {
		crop = &Crop{}
	}

	crop.Size += p.Crop.Size
	crop.Yield += (p.Crop.Size * p.Crop.Yield)

	ls.TotalSize += p.Crop.Size
	ls.Count += 1
	ls.Commodities[p.Crop.Commodity] = crop
}

func NewLandSummary() *LandSummary {
	return &LandSummary{
		Commodities: map[string]*Crop{},
		Areas:       map[string]*LandSummary{},
	}
}
