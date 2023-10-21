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
	err := conn.SetPeople(append(unique(p.Adults), unique(p.Children)...)...)
	if err != nil {
		return err
	}
	err = conn.SetEvents(unique(p.Events)...)
	if err != nil {
		return err
	}
	err = conn.SetTuples(db.RelationPersonProfessionSkill, unique(p.Skills)...)
	if err != nil {
		return err
	}
	err = conn.SetTuples(db.RelationPersonReligionFaith, unique(p.Faith)...)
	if err != nil {
		return err
	}
	err = conn.SetTuples(db.RelationPersonPersonTrust, unique(p.Trust)...)
	if err != nil {
		return err
	}
	err = conn.SetTuples(db.RelationPersonPersonRelationship, unique(p.Relations)...)
	if err != nil {
		return err
	}
	return conn.SetFamilies(p.Families...)
}
