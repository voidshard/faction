package base

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"
)

// metaPeople is a working data set when operating on people & associated data
type metaPeople struct {
	adults    []*structs.Person
	children  []*structs.Person
	skills    []*structs.Tuple
	faith     []*structs.Tuple
	trust     []*structs.Tuple
	relations []*structs.Tuple
	families  []*structs.Family
	events    []*structs.Event
}

func newMetaPeople() *metaPeople {
	return &metaPeople{
		adults:    []*structs.Person{},
		children:  []*structs.Person{},
		skills:    []*structs.Tuple{},
		faith:     []*structs.Tuple{},
		trust:     []*structs.Tuple{},
		relations: []*structs.Tuple{},
		families:  []*structs.Family{},
		events:    []*structs.Event{},
	}
}

func writeMetaPeople(conn *db.FactionDB, p *metaPeople) error {
	return conn.InTransaction(func(tx db.ReaderWriter) error {
		err := tx.SetPeople(append(p.adults, p.children...)...)
		if err != nil {
			return err
		}
		err = tx.SetEvents(p.events...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonProfessionSkill, p.skills...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonReligionFaith, p.faith...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonPersonTrust, p.trust...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonPersonRelationship, p.relations...)
		if err != nil {
			return err
		}
		return tx.SetFamilies(p.families...)
	})
}
