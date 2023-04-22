package premade

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

func Affiliation() *config.Affiliation {
	return &config.Affiliation{
		Affiliation: config.Distribution{
			Min:       10,
			Max:       structs.MaxTuple,
			Mean:      structs.MaxTuple / 2,
			Deviation: structs.MaxTuple / 2,
		},
		Trust: config.Distribution{
			Min:       structs.MinTuple,
			Max:       structs.MaxTuple,
			Mean:      0,
			Deviation: structs.MaxTuple * 3 / 10,
		},
		EthosDistance:  0.3,
		OutlawedWeight: 0.8,
		Members: config.Distribution{
			Min:       50,
			Max:       2050,
			Mean:      1050,
			Deviation: 1000,
		},
	}
}