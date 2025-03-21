package v1

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/voidshard/faction/internal/util/flatten"
)

// unmarshalObject takes a byte slice or object and unmarshals it into
// the given Object
func unmarshalObject(in interface{}, obj Object) error {
	var err error
	data, ok := in.([]byte)
	if !ok { // if it's not a byte slice, try to marshal it
		data, err = yaml.Marshal(in)
		if err != nil {
			return err
		}
	}
	err = yaml.Unmarshal(data, obj)
	return err
}

// GetFields turns a struct into a flat map of key value pairs that
// can be searched.
func GetFields(o Object) (map[string]interface{}, error) {
	// turn object into map[string]interface{} so we can flatten
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	// flatten the map (making map keys into dot separated strings)
	flat, err := flatten.Flatten(m, "", flatten.DotStyle)
	if err != nil {
		return nil, err
	}

	// lower case keys & remove dot or _ prefixes
	final := map[string]interface{}{}
	for k, v := range flat {
		if strings.HasPrefix(k, ".") || strings.HasPrefix(k, "_") {
			k = k[1:]
		}
		final[strings.ToLower(k)] = v
	}

	// make sure default is set
	final["controller"] = o.GetController()

	return final, nil
}
