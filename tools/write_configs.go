package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
)

func main() {
	for key, value := range map[string]interface{}{
		"demographics_fantasy_human": fantasy.DemographicsHuman(),
		"actions_fantasy":            fantasy.Actions(),
		"affiliation_fantasy":        fantasy.Affiliation(),
		"faction_fantasy":            fantasy.Faction(),
	} {
		data, err := yaml.Marshal(value)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(fmt.Sprintf("example_%s.yaml", key), data, 0644)
		if err != nil {
			panic(err)
		}
	}
}
