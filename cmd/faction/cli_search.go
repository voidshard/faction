package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/voidshard/faction/pkg/client"
	"github.com/voidshard/faction/pkg/util/log"
)

type cliSearchCmd struct {
	optCliConn
	optGeneral
	optCliGlobal

	Object struct {
		Kind string `positional-arg-name:"object" description:"Object to get"`
	} `positional-args:"true" required:"true"`

	Limit        int64   `long:"limit" short:"l" default:"1" description:"Limit number of results"`
	RandomWeight float64 `long:"random-weight" short:"r" default:"0" description:"Weight for random selection"`

	All []string `long:"all" short:"a" description:"Match vs all docs. Expects field[=><]value pairs."`
	Any []string `long:"any" short:"o" description:"Match vs any docs. Expects field[=><]value pairs."`
	Not []string `long:"not" short:"n" description:"Exclude docs. Expects field[=><]value pairs."`

	Score []string `long:"score" short:"s" description:"Score docs. Expects weight:field[=><]value tuples."`
}

func (c *cliSearchCmd) Execute(args []string) error {
	c.Object.Kind = validKind(c.Object.Kind)
	if c.Object.Kind == "" {
		return fmt.Errorf("invalid object kind %s", c.Object.Kind)
	}

	if c.World == "" {
		return fmt.Errorf("world must be set for search operations")
	}

	conn, err := client.New(client.NewConfig())
	if err != nil {
		return err
	}

	search := conn.Search(c.World, c.Object.Kind, c.Limit)
	search.RandomWeight(c.RandomWeight)
	for _, a := range c.All {
		log.Debug().Str("input", a).Msg("[All] parsing match")
		field, op, value, _ := parseMatch(a)
		if field == "" {
			log.Warn().Str("input", a).Msg("failed to parse match")
			continue
		}
		search.All(field, value, op)
	}
	for _, a := range c.Any {
		log.Debug().Str("input", a).Msg("[Any] parsing match")
		field, op, value, _ := parseMatch(a)
		if field == "" {
			log.Warn().Str("input", a).Msg("failed to parse match")
			continue
		}
		search.Any(field, value, op)
	}
	for _, a := range c.Not {
		log.Debug().Str("input", a).Msg("[Not] parsing match")
		field, op, value, _ := parseMatch(a)
		if field == "" {
			log.Warn().Str("input", a).Msg("failed to parse match")
			continue
		}
		search.Not(field, value, op)
	}
	for _, a := range c.Score {
		log.Debug().Str("input", a).Msg("[Score] parsing match")
		field, op, value, weight := parseMatch(a)
		if field == "" {
			log.Warn().Str("input", a).Msg("failed to parse match")
			continue
		}
		search.Score(field, value, weight, op)
	}

	objs, err := search.Do()
	if err != nil {
		return err
	}

	yamlData, err := dumpYaml(objs)
	if yamlData != nil {
		fmt.Println(string(yamlData))
	}
	return err
}

func parseMatch(in string) (string, client.Operation, interface{}, float64) {
	bits := strings.SplitN(in, ":", 2)
	remainder := ""

	// determine the weight (in a string w:field=value we want the w)
	weight := 0.0
	if len(bits) > 1 {
		remainder = bits[1]
		weight, _ = strconv.ParseFloat(bits[0], 64)
	} else {
		remainder = bits[0]
	}

	// break the rest of the string up (nb: we don't know the operator)
	symbols := map[string]client.Operation{
		"=": client.Equal,
		">": client.GreaterThan,
		"<": client.LessThan,
	}

	var operation client.Operation
	var field string
	var value string

	for symbol, op := range symbols {
		bits = strings.SplitN(remainder, symbol, 2)
		if len(bits) > 1 {
			field = bits[0]
			value = bits[1]
			operation = op
			break
		}
	}

	if field == "" {
		// ??
		return "", client.Equal, "", 0
	}

	// parse the value so it's the correct type
	var finalValue interface{}
	if operation == client.GreaterThan || operation == client.LessThan {
		var err error
		finalValue, err = strconv.ParseFloat(value, 64)
		if err != nil {
			finalValue = 0.0
		}
		log.Debug().Str("input", in).Err(err).Msg("parsed value as float")
	} else if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
		// accept true/false in any case as a boolean
		finalValue = strings.ToLower(value) == "true"
		log.Debug().Str("input", in).Msg("parsed value as boolean")
	} else {
		// try to parse as a float, if it fails, assume it's a string
		var err error
		finalValue, err = strconv.ParseFloat(value, 64)
		if err == nil {
			log.Debug().Str("input", in).Msg("parsed value as float")
		} else {
			finalValue = value
			log.Debug().Str("input", in).Msg("parsed value as string")
		}
	}
	log.Debug().Str("input", in).Str("field", field).Str("value", value).Str("operation", string(operation)).Msg("parsed match")

	return field, operation, finalValue, weight
}
