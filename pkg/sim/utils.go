package sim

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/voidshard/faction/internal/log"
	base "github.com/voidshard/faction/internal/sim/base"
	"github.com/voidshard/faction/pkg/config"
	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
)

const (
	envCfgLocation = "FACTION_SIMULATION_CFG_FILE"
)

func LoadConfig(cfgpath string) (*config.Simulation, error) {
	f, err := os.Open(cfgpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &config.Simulation{}
	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func SetDefaults(cfg *config.Simulation) {
	if cfg.Database == nil {
		cfg.Database = config.DefaultDatabase()
	}

	if cfg.Queue == nil {
		cfg.Queue = config.DefaultQueue()
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
	if cfg.Settings == nil {
		cfg.Settings = config.DefaultSimulationSettings()
	}
}

func merge(a, b *config.Simulation) {
	if a.Database == nil {
		a.Database = b.Database
	}
	if a.Queue == nil {
		a.Queue = b.Queue
	}
	if a.Races == nil {
		a.Races = b.Races
	}
	if a.Cultures == nil {
		a.Cultures = b.Cultures
	}
	if a.Actions == nil {
		a.Actions = b.Actions
	}
	if a.Settings == nil {
		a.Settings = b.Settings
	}
}

func New(cfg *config.Simulation, opts ...simOption) (Simulation, error) {
	if cfg == nil {
		cfg = &config.Simulation{}
	}

	// read config file from disk, this will form our "defaults" for
	// anything not set in our given cfg
	diskcfg := &config.Simulation{}
	if cfgpath := os.Getenv(envCfgLocation); cfgpath != "" {
		// if we have a config file, try load it
		var err error
		diskcfg, err = LoadConfig(cfgpath)
		if err != nil {
			return nil, err
		}
	}
	SetDefaults(diskcfg)

	// merge defaults with given cfg where required
	merge(cfg, diskcfg)

	// build a sim struct
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

	log.Debug().Str("driver", string(cfg.Database.Driver)).Msg()("sim configured with database")
	log.Debug().Str("driver", string(cfg.Queue.Driver)).Msg()("sim configured with queue")
	log.Debug().Int("count", len(cfg.Actions)).Msg()("loaded actions")
	log.Debug().Int("count", len(cfg.Races)).Msg()("loaded races")
	log.Debug().Int("count", len(cfg.Cultures)).Msg()("loaded cultures")

	return me, nil
}
