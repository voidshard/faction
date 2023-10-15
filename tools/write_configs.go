package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/voidshard/faction/pkg/config"
	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
	"github.com/voidshard/faction/pkg/sim"
)

func main() {
	cfg := &config.Simulation{}
	sim.SetDefaults(cfg)

	for key, value := range map[string]interface{}{
		"race_fantasy_human":    fantasy.RaceHuman(),
		"culture_fantasy_human": fantasy.CultureHuman(),
		"actions_fantasy":       fantasy.Actions(),
		"affiliation_fantasy":   fantasy.Affiliation(),
		"faction_fantasy":       fantasy.Faction(),
		"government_fantasy":    fantasy.Government(),
		"simulation_fantasy":    cfg,
	} {
		data, err := yaml.Marshal(value)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(fmt.Sprintf("example_%s.yaml", key), data, 0644)
		if err != nil {
			panic(err)
		}
	}
}
