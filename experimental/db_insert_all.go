package main

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/premade/fantasy"
	"github.com/voidshard/faction/pkg/structs"

	"github.com/voidshard/faction/internal/db"
)

func main() {
	// inserts at least one of each struct into the DB so we call
	// all of the SET sql statements.

	conn, err := db.New(config.DefaultDatabase())
	if err != nil {
		panic(err)
	}

	govt1 := &structs.Government{ID: structs.NewID(), Outlawed: &structs.Laws{
		Commodities: map[string]bool{"apples": true, "fur": false},
		Actions:     map[string]bool{fantasy.War: true, fantasy.Harvest: false},
	}}

	area1 := &structs.Area{ID: structs.NewID(), GovernmentID: govt1.ID}
	area2 := &structs.Area{ID: structs.NewID()}

	faction1 := &structs.Faction{
		ID:           structs.NewID(),
		Name:         "Genesis",
		GovernmentID: govt1.ID,
		HomeAreaID:   area1.ID,
		IsGovernment: true,
	}

	plotland1 := &structs.Plot{
		ID: structs.NewID(), AreaID: area2.ID, FactionID: faction1.ID, Crop: structs.Crop{Commodity: "apples", Yield: 10},
	}
	plotland2 := &structs.Plot{
		ID: structs.NewID(), AreaID: area1.ID, FactionID: faction1.ID, Crop: structs.Crop{Commodity: "fur", Yield: 5},
	}

	husband1 := &structs.Person{
		ID:            structs.NewID(),
		FirstName:     "Adam",
		BirthFamilyID: structs.NewID(),
		AreaID:        area1.ID,
		IsMale:        true,
		BirthTick:     1,
		Race:          "human",
		Culture:       "human",
	}
	wife1 := &structs.Person{
		ID:            structs.NewID(),
		FirstName:     "Eve",
		BirthFamilyID: structs.NewID(),
		AreaID:        area1.ID,
		IsMale:        false,
		BirthTick:     1,
		Race:          "human",
		Culture:       "human",
	}
	steve1 := &structs.Person{
		ID:            structs.NewID(),
		FirstName:     "Steve",
		BirthFamilyID: structs.NewID(),
		AreaID:        area2.ID,
		IsMale:        true,
		BirthTick:     2,
		Race:          "human",
		Culture:       "human",
	}
	family1 := &structs.Family{
		ID:             structs.NewID(),
		AreaID:         area1.ID,
		IsChildBearing: true,
		MaleID:         husband1.ID,
		FemaleID:       wife1.ID,
		Race:           "human",
		Culture:        "human",
	}

	// RelationPersonPersonRelationships
	relation1 := &structs.Tuple{Subject: husband1.ID, Object: wife1.ID, Value: int(structs.PersonalRelationWife)}
	relation2 := &structs.Tuple{Subject: wife1.ID, Object: husband1.ID, Value: int(structs.PersonalRelationHusband)}

	// RelationPersonPersonTrust
	trust1 := &structs.Tuple{Subject: husband1.ID, Object: wife1.ID, Value: 100}
	trust2 := &structs.Tuple{Subject: wife1.ID, Object: husband1.ID, Value: 100}
	trust3 := &structs.Tuple{Subject: husband1.ID, Object: steve1.ID, Value: 10}
	trust4 := &structs.Tuple{Subject: wife1.ID, Object: steve1.ID, Value: 12}
	trust5 := &structs.Tuple{Subject: steve1.ID, Object: husband1.ID, Value: 50}
	trust6 := &structs.Tuple{Subject: steve1.ID, Object: wife1.ID, Value: 70}

	mod1 := &structs.Modifier{
		Tuple:       structs.Tuple{Subject: husband1.ID, Object: wife1.ID, Value: -10},
		TickExpires: 10,
		MetaKey:     structs.MetaKeyPerson,
		MetaVal:     wife1.ID,
		MetaReason:  "listening to serpents",
	}
	mod2 := &structs.Modifier{
		Tuple:       structs.Tuple{Subject: husband1.ID, Object: wife1.ID, Value: 15},
		TickExpires: 20,
		MetaKey:     structs.MetaKeyPerson,
		MetaVal:     wife1.ID,
		MetaReason:  "being wonderful",
	}
	mod3 := &structs.Modifier{
		Tuple:       structs.Tuple{Subject: wife1.ID, Object: husband1.ID, Value: 25},
		TickExpires: 15,
		MetaKey:     structs.MetaKeyPerson,
		MetaVal:     husband1.ID,
		MetaReason:  "naming animals",
	}
	mod4 := &structs.Modifier{
		Tuple:       structs.Tuple{Subject: husband1.ID, Object: steve1.ID, Value: -20},
		TickExpires: 12,
		MetaKey:     structs.MetaKeyPerson,
		MetaVal:     steve1.ID,
		MetaReason:  "being suspicious",
	}

	// RelationPersonFactionAffiliation
	affiliation1 := &structs.Tuple{Subject: husband1.ID, Object: faction1.ID, Value: 100}
	affiliation2 := &structs.Tuple{Subject: wife1.ID, Object: faction1.ID, Value: 50}

	// A house for our family
	plot1 := &structs.Plot{
		ID: structs.NewID(), AreaID: area1.ID, FactionID: faction1.ID, Crop: structs.Crop{Size: 10},
	}

	job1 := &structs.Job{
		ID:              structs.NewID(),
		SourceFactionID: faction1.ID,
		SourceAreaID:    area1.ID,
		Action:          fantasy.Festival,
		TargetAreaID:    area1.ID,
		TargetFactionID: faction1.ID,
		TargetMetaKey:   structs.MetaKeyFamily,
		TargetMetaVal:   family1.ID,
		PeopleMin:       2,
		PeopleMax:       100,
		TickCreated:     0,
		TickStarts:      5,
		TickEnds:        15,
		Secrecy:         1,
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
	fmt.Println("Tick set to 1")

	err = tx.SetAreas(area1, area2)
	errRollback(err)
	fmt.Println("Areas set")

	err = tx.SetGovernments(govt1)
	errRollback(err)
	fmt.Println("Governments set")

	err = tx.SetFactions(faction1)
	errRollback(err)
	fmt.Println("Factions set")

	err = tx.SetPeople(husband1, wife1, steve1)
	errRollback(err)
	fmt.Println("People set")

	err = tx.SetFamilies(family1)
	errRollback(err)
	fmt.Println("Families set")

	err = tx.SetPlots(plot1, plotland1, plotland2)
	errRollback(err)
	fmt.Println("Plots set")

	err = tx.SetTuples(db.RelationPersonPersonRelationship, relation1, relation2)
	errRollback(err)
	fmt.Println("Relations set")

	err = tx.SetTuples(db.RelationPersonPersonTrust, trust1, trust2, trust3, trust4, trust5, trust6)
	errRollback(err)
	fmt.Println("Trust set")

	err = tx.SetTuples(db.RelationPersonFactionAffiliation, affiliation1, affiliation2)
	errRollback(err)
	fmt.Println("Affiliations set")

	err = tx.SetModifiers(db.RelationPersonPersonTrust, mod1, mod2, mod3, mod4)
	errRollback(err)
	fmt.Println("Modifiers set")

	err = tx.SetJobs(job1)
	errRollback(err)
	fmt.Println("Jobs set")

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
