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
	// de-dupe stuff before we enter a transaction.
	// By definition anything we "Set" overwrites any existing data so we only need to apply
	// the last object in the slice for each ID.
	// Technically the DB could resolve this (depending on the engine) but this is busy work
	// we can do out-of-transaction
	p.Adults = unique(p.Adults)
	p.Children = unique(p.Children)
	p.Skills = unique(p.Skills)
	p.Faith = unique(p.Faith)
	p.Trust = unique(p.Trust)
	p.Relations = unique(p.Relations)
	p.Families = unique(p.Families)
	p.Events = unique(p.Events)
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
