package sim

import (
	"os"

	"gopkg.in/yaml.v3"

	base "github.com/voidshard/faction/internal/sim/base"
	"github.com/voidshard/faction/pkg/config"
	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
)

const (
	envDBLocation    = "FACTION_DB_LOCATION"
	envRedisLocation = "FACTION_REDIS_LOCATION"
	envCfgLocation   = "FACTION_SIMULATION_CFG_FILE"
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
	if dbLocation := os.Getenv(envDBLocation); dbLocation != "" { // if set, override
		cfg.Database.Location = dbLocation
	}

	if cfg.Queue == nil {
		cfg.Queue = config.DefaultQueue()
	}
	if redisLocation := os.Getenv(envRedisLocation); redisLocation != "" { // if set, override
		cfg.Queue.Location = redisLocation
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

func New(cfg *config.Simulation, opts ...simOption) (Simulation, error) {
	// apply default settings
	if cfg == nil {
		if cfgpath := os.Getenv(envCfgLocation); cfgpath != "" {
			// if we have a config file, try load it
			var err error
			cfg, err = LoadConfig(cfgpath)
			if err != nil {
				return nil, err
			}
		} else {
			// otherwise, set us a blank
			cfg = &config.Simulation{}
		}
	}

	SetDefaults(cfg)

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
