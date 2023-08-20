package simutil

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"
)

// MetaPeople is a working data set when operating on people & associated data
type MetaPeople struct {
	Adults    []*structs.Person
	Children  []*structs.Person
	Skills    []*structs.Tuple
	Faith     []*structs.Tuple
	Trust     []*structs.Tuple
	Relations []*structs.Tuple
	Families  []*structs.Family
	Events    []*structs.Event
}

func NewMetaPeople() *MetaPeople {
	return &MetaPeople{
		Adults:    []*structs.Person{},
		Children:  []*structs.Person{},
		Skills:    []*structs.Tuple{},
		Faith:     []*structs.Tuple{},
		Trust:     []*structs.Tuple{},
		Relations: []*structs.Tuple{},
		Families:  []*structs.Family{},
		Events:    []*structs.Event{},
	}
}

func WriteMetaPeople(conn *db.FactionDB, p *MetaPeople) error {
	return conn.InTransaction(func(tx db.ReaderWriter) error {
		err := tx.SetPeople(append(p.Adults, p.Children...)...)
		if err != nil {
			return err
		}
		err = tx.SetEvents(p.Events...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonProfessionSkill, p.Skills...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonReligionFaith, p.Faith...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonPersonTrust, p.Trust...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonPersonRelationship, p.Relations...)
		if err != nil {
			return err
		}
		return tx.SetFamilies(p.Families...)
	})
}
