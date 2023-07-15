package base

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"
)

func newBirthEvent(tick int, personID, familyID string) *structs.Event {
	return &structs.Event{
		ID:             dbutils.NewID(),
		Type:           structs.EventPersonBirth,
		Tick:           tick,
		SubjectMetaKey: structs.MetaKeyPerson,
		SubjectMetaVal: personID,
		CauseMetaKey:   structs.MetaKeyFamily,
		CauseMetaVal:   familyID,
	}
}

// already have the families to hand).
func (s *Base) applyDeath(tick int, in ...*structs.Person) error {
	if len(in) == 0 {
		return nil
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

	token := (&dbutils.IterToken{Limit: 1000, Offset: 0}).String()
	var (
		families []*structs.Family
		err      error
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

		update := []*structs.Family{}
		rels := []*structs.Tuple{}

		for _, f := range families {
			if f.WidowedTick > 0 {
				// implies the other partner is already dead
				continue
			}

			f.IsChildBearing = false
			f.WidowedTick = tick

			// if a mother has died, the child has too :(
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
