package base

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"
)

func (s *Base) applyDeathFamilyEffect(tick int, events []*structs.Event) error {
	query := db.Q(
		db.F(db.ID, db.In, eventSubjects(events)),
	).DisableSort()
	in, _, err := s.dbconn.People("", query)
	if err != nil {
		return err
	}

	female := []string{}
	male := []string{}
	femaleDead := map[string]bool{}

	for _, d := range in {
		if d.IsMale {
			male = append(male, d.ID)
		} else {
			female = append(female, d.ID)
			femaleDead[d.ID] = true
		}
	}

	token := dbutils.NewTokenWith(1000, 0)
	var (
		families []*structs.Family
	)
	q := db.Q(
		db.F(db.MaleID, db.In, male),
	).Or(
		db.F(db.FemaleID, db.In, female),
	)
	for {
		families, token, err = s.dbconn.Families(token, q)
		if err != nil {
			return err
		}
		if len(families) == 0 {
			return nil
		}

		update := []*structs.Family{}
		rels := []*structs.Tuple{}

		for _, f := range families {
			if f.WidowedTick > 0 {
				// implies the other partner is already dead
				continue
			}

			f.IsChildBearing = false
			f.WidowedTick = tick

			// if a mother has died, the unborn child has too :(
			_, ok := femaleDead[f.FemaleID]
			if ok {
				f.PregnancyEnd = 0
			}

			rels = append(
				rels,
				&structs.Tuple{Subject: f.MaleID, Object: f.FemaleID, Value: int(structs.PersonalRelationExWife)},
				&structs.Tuple{Subject: f.FemaleID, Object: f.MaleID, Value: int(structs.PersonalRelationExHusband)},
			)
			update = append(update, f)
		}

		s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
			err = tx.SetFamilies(update...)
			if err != nil {
				return err
			}
			return tx.SetTuples(db.RelationPersonPersonRelationship, rels...)
		})

		if token == "" {
			break
		}
	}

	return nil
}
