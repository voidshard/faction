package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/voidshard/faction/pkg/premade"
)

func main() {
	for key, value := range map[string]interface{}{
		"demographics_fantasy_human": premade.DemographicsFantasyHuman(),
		"actions_fantasy":            premade.ActionsFantasy(),
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
