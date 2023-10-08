package base

import (
	"github.com/voidshard/faction/pkg/structs"
)

func (s *Base) applyFamilyMarriage(tick int, events []*structs.Event) error {
	/*
		if s.cfg.Cultures == nil {
			return nil
		}

		query := db.Q(db.F(db.ID, db.In, eventSubjects(events)))
		in, _, err := s.dbconn.Families(dbutils.NewTokenWith(len(events), 0), query)
		if err != nil {
			return err
		}

		for _, fami := range in {
			continue
		}
	*/
	return nil
}

// defaultSocialClass
func defaultSocialClass(s map[string]int) string {
	if s == nil {
		return ""
	}
	lowest := ""
	lowestVal := 0
	for name, val := range s {
		if val < lowestVal {
			lowest = name
			lowestVal = val
		}
	}
	return lowest
}
