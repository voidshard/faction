package main

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"

	"github.com/voidshard/faction/internal/db"
)

func main() {
	// inserts at least one of each struct into the DB so we call
	// all of the SET sql statements.

	cfg := &config.Database{
		Driver:   config.DatabaseSQLite3,
		Name:     "test.sqlite",
		Location: "/tmp",
	}
	conn, err := db.New(cfg)
	if err != nil {
		panic(err)
	}

	govt1 := &structs.Government{ID: structs.NewID()}

	// RelationLawGovernmentToCommodidty
	law1 := &structs.Tuple{Subject: govt1.ID, Object: "apples", Value: 1} // illegal
	law2 := &structs.Tuple{Subject: govt1.ID, Object: "fur", Value: 0}    // legal

	// RelationLawGovernmentToAction
	law3 := &structs.Tuple{Subject: govt1.ID, Object: string(structs.ActionTypeWar), Value: 1} // illegal

	area1 := &structs.Area{ID: structs.NewID(), GoverningFactionID: govt1.ID}
	area2 := &structs.Area{ID: structs.NewID()}
	landright1 := &structs.LandRight{
		ID: structs.NewID(), AreaID: area2.ID, GoverningFactionID: govt1.ID, Resource: "apples",
	}
	landright2 := &structs.LandRight{
		ID: structs.NewID(), AreaID: area1.ID, GoverningFactionID: govt1.ID, Resource: "fur",
	}

	route1 := &structs.Route{SourceAreaID: area1.ID, TargetAreaID: area2.ID, TravelTime: 3}  // downhill
	route2 := &structs.Route{SourceAreaID: area2.ID, TargetAreaID: area1.ID, TravelTime: 50} // uphill

	faction1 := &structs.Faction{
		ID:           structs.NewID(),
		Name:         "Genesis",
		GovernmentID: govt1.ID,
		IsGovernment: true,
	}

	husband1 := &structs.Person{
		ID:            structs.NewID(),
		FirstName:     "Adam",
		BirthFamilyID: structs.NewID(),
		AreaID:        area1.ID,
		IsMale:        true,
		BirthTick:     1,
	}
	wife1 := &structs.Person{
		ID:            structs.NewID(),
		FirstName:     "Eve",
		BirthFamilyID: structs.NewID(),
		AreaID:        area1.ID,
		IsMale:        false,
		BirthTick:     1,
	}
	family1 := &structs.Family{
		ID: structs.NewID(), AreaID: area1.ID, IsChildBearing: true, MaleID: husband1.ID, FemaleID: wife1.ID,
	}

	// RelationPersonPersonRelationships
	relation1 := &structs.Tuple{Subject: husband1.ID, Object: wife1.ID, Value: int(structs.PersonalRelationWife)}
	relation2 := &structs.Tuple{Subject: wife1.ID, Object: husband1.ID, Value: int(structs.PersonalRelationHusband)}

	// RelationPersonPersonTrust
	trust1 := &structs.Tuple{Subject: husband1.ID, Object: wife1.ID, Value: 100}
	trust2 := &structs.Tuple{Subject: wife1.ID, Object: husband1.ID, Value: 100}

	mod1 := new(structs.Modifier)
	mod1.Subject = husband1.ID
	mod1.Object = wife1.ID
	mod1.Value = -10
	mod1.TickExpires = 10
	mod1.MetaKey = structs.MetaKeyPerson
	mod1.MetaVal = wife1.ID
	mod1.MetaReason = "listening to serpents"

	// RelationPersonFactionAffiliation
	affiliation1 := &structs.Tuple{Subject: husband1.ID, Object: faction1.ID, Value: 100}
	affiliation2 := &structs.Tuple{Subject: wife1.ID, Object: faction1.ID, Value: 50}

	// A house for our family
	plot1 := &structs.Plot{
		ID: structs.NewID(), AreaID: area1.ID, OwnerFactionID: faction1.ID, Size: 10,
	}

	job1 := &structs.Job{
		ID:              structs.NewID(),
		SourceFactionID: faction1.ID,
		SourceAreaID:    area1.ID,
		Action:          structs.ActionTypeFestival,
		TargetAreaID:    area1.ID,
		TargetMetaKey:   structs.MetaKeyFamily,
		TargetMetaVal:   family1.ID,
		PeopleMin:       2,
		PeopleMax:       100,
		TickCreated:     0,
		TickStarts:      5,
		TickEnds:        15,
		Secrecy:         0,
		IsIllegal:       false,
		State:           structs.JobStatePending,
	}

	tx, err := conn.Transaction()
	if err != nil {
		panic(err)
	}
	errRollback := func(err error) {
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.SetTick(1)
	errRollback(err)

	err = tx.SetAreas(area1, area2)
	errRollback(err)

	err = tx.SetLandRights(landright1, landright2)
	errRollback(err)

	err = tx.SetRoutes(route1, route2)
	errRollback(err)

	err = tx.SetGovernments(govt1)
	errRollback(err)

	err = tx.SetFactions(faction1)
	errRollback(err)

	err = tx.SetPeople(husband1, wife1)
	errRollback(err)

	err = tx.SetFamilies(family1)
	errRollback(err)

	err = tx.SetPlots(plot1)
	errRollback(err)

	err = tx.SetTuples(db.RelationLawGovernmentToCommodidty, law1, law2)
	errRollback(err)

	err = tx.SetTuples(db.RelationLawGovernmentToAction, law3)
	errRollback(err)

	err = tx.SetTuples(db.RelationPersonPersonRelationship, relation1, relation2)
	errRollback(err)

	err = tx.SetTuples(db.RelationPersonPersonTrust, trust1, trust2)
	errRollback(err)

	err = tx.SetTuples(db.RelationPersonFactionAffiliation, affiliation1, affiliation2)
	errRollback(err)

	err = tx.SetModifiers(db.RelationPersonPersonTrust, mod1)
	errRollback(err)

	err = tx.SetJobs(job1)
	errRollback(err)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
