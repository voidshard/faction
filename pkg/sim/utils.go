package sim

import (
	base "github.com/voidshard/faction/internal/sim/base"
	"github.com/voidshard/faction/pkg/config"
	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
)

func New(cfg *config.Simulation, opts ...simOption) (Simulation, error) {
	// apply default settings
	if cfg == nil {
		cfg = &config.Simulation{}
	}
	if cfg.Database == nil {
		cfg.Database = config.DefaultDatabase()
	}

	// if not told otherwise, we'll assume human demographics
	if cfg.Races == nil {
		cfg.Races = map[string]*config.Race{
			"human": fantasy.RaceHuman(),
		}
	}
	if cfg.Cultures == nil {
		cfg.Cultures = map[string]*config.Culture{
			"human": fantasy.CultureHuman(),
		}
	}
	if cfg.Actions == nil {
		cfg.Actions = fantasy.Actions()
	}

	me, err := base.New(cfg)
	if err != nil {
		return nil, err
	}

	// apply given options
	for _, opt := range opts {
		err = opt(me)
		if err != nil {
			return nil, err
		}
	}

	return me, nil
}
