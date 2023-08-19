package base

import (
	"fmt"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"

	mapset "github.com/deckarep/golang-set/v2"
)

// applyBirthSiblingRelations will find all people born in the same family and
// create a sibling relationship between them.
func (s *Base) applyBirthSiblingRelations(tick int, events []*structs.Event) error {
	query := db.Q(
		db.F(db.ID, db.In, eventSubjects(events)),
	).DisableSort()
	in, _, err := s.dbconn.People("", query)
	if err != nil {
		return err
	}

	newChildren := map[string][]*structs.Person{} // familyID -> []Person
	families := mapset.NewSet[string]()
	for _, p := range in {
		families.Add(p.BirthFamilyID)

		kids, ok := newChildren[p.BirthFamilyID]
		if !ok {
			kids = []*structs.Person{}
		}
		newChildren[p.BirthFamilyID] = append(kids, p)
	}

	pf := db.Q(db.F(db.BirthFamilyID, db.In, families.ToSlice()))
	var (
		people []*structs.Person
		token  string
	)
	for {
		people, token, err = s.dbconn.People(token, pf)
		if err != nil {
			return err
		}

		mp := newMetaPeople()
		for _, p := range people {
			kids, ok := newChildren[p.BirthFamilyID]
			if !ok {
				continue
			}

			for _, kid := range kids {
				if p.ID == kid.ID {
					continue
				}
				if !s.dice.IsValidDemographic(p.Race, p.Culture) {
					return fmt.Errorf("invalid demographic not found: [race] %s, [culture] %s", p.Race, p.Culture)
				}
				demo := s.dice.MustDemographic(p.Race, p.Culture)
				siblingRelationship(demo, mp, p, kid)
				mp.children = append(mp.children, p)

			}
		}

		err = writeMetaPeople(s.dbconn, mp)
		if err != nil {
			return err
		}

		if token == "" {
			break
		}
	}

	return nil
}
