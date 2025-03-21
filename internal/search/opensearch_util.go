package search

import (
	"encoding/json"

	v1 "github.com/voidshard/faction/pkg/structs/v1"
)

// toOpensearchQuery converts a v1.Query to an opensearch query.
// Example of final query
/*
{
  "size": 10,
  "query": {
    "function_score": {
      "query": {
        "bool": {
          "must": [
            {"term": {"race": "human"}}
          ],
          "should": [
            {"range": {"ethos.law": {"gte": 80}}},
            {"range": {"ethos.good": {"gte": 80}}}
          ]
        }
      },
      "functions": [
        {
          "filter": { "match_all": { } },
          "random_score": {},
          "weight": 2
        },
        {
          "filter": {
            "term": {
              "labels.human/variant": "fox"
            }
          },
          "weight": 5
        },
        {
          "filter": {
            "term": {
              "labels.green/variant": "wood"
            }
          },
          "weight": 2
        },
        {
          "filter": {
            "range": {
              "professions.mage.level": {"gte": 1}
            }
          },
          "weight": 0
        }
      ]
    }
  }
}
*/
func toOpensearchQuery(q *v1.Query) (string, error) {
	// build query filters
	must := []map[string]interface{}{}
	for _, f := range q.Filter.All {
		must = append(must, toOsFilter(&f))
	}

	should := []map[string]interface{}{}
	for _, f := range q.Filter.Any {
		should = append(should, toOsFilter(&f))
	}

	must_not := []map[string]interface{}{}
	for _, f := range q.Filter.Not {
		must_not = append(must_not, toOsFilter(&f))
	}

	// assemble filters into query
	query := map[string]interface{}{}
	if len(must) > 0 {
		query["must"] = must
	}
	if len(should) > 0 {
		query["should"] = should
	}
	if len(must_not) > 0 {
		query["must_not"] = must_not
	}
	if len(query) == 0 {
		// if no filters are given, match all
		query = map[string]interface{}{"match_all": map[string]interface{}{}}
	} else {
		query = map[string]interface{}{"bool": query}
	}

	// build scoring filters
	score := []map[string]interface{}{}
	if q.RandomWeight > 0 {
		// add random scoring
		score = append(score, map[string]interface{}{
			"filter":       map[string]interface{}{"match_all": map[string]interface{}{}},
			"random_score": map[string]interface{}{},
			"weight":       q.RandomWeight,
		})
	}
	for _, s := range q.Score {
		// score docs based on given scoring filters
		score = append(score, map[string]interface{}{
			"filter": toOsFilter(&s.Match),
			"weight": s.Weight,
		})
	}

	// build final query
	document := map[string]interface{}{
		"size": q.Limit,
		"query": map[string]interface{}{
			"function_score": map[string]interface{}{
				"query":     query,
				"functions": score,
			},
		},
	}

	data, err := json.Marshal(document)
	return string(data), err
}

func toOsFilter(f *v1.Match) map[string]interface{} {
	switch f.Op {
	case "eq", "":
		return map[string]interface{}{
			"term": map[string]interface{}{f.Field: f.Value},
		}
	case "lt", "gt", "lte", "gte":
		return map[string]interface{}{
			"range": map[string]interface{}{
				f.Field: map[string]interface{}{f.Op: f.Value},
			},
		}
	}
	return nil
}
