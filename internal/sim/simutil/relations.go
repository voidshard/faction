package simutil

import (
	dg "github.com/voidshard/faction/internal/random/demographics"
	"github.com/voidshard/faction/pkg/structs"
)

func AddChildToFamily(dice *dg.Demographic, child *structs.Person, f *structs.Family) {
	child.Ethos = structs.EthosAverage(child.Ethos, f.Ethos)
	child.BirthFamilyID = f.ID
	child.Race = f.Race
	child.Culture = f.Culture
	f.NumberOfChildren++
	if f.NumberOfChildren >= dice.MaxFamilySize() {
		f.IsChildBearing = false
	}
}

func ProfessionalRelationByTrust(a int) structs.PersonalRelation {
	if a < structs.MaxTuple/10 && a > structs.MinTuple/10 {
		// the neutral zone
		return structs.PersonalRelation_Colleague
	} else if a < 0 {
		if a < structs.MinTuple/2 {
			return structs.PersonalRelation_HatedEnemy
		}
		return structs.PersonalRelation_Enemy
	} else { // a > 0
		if a > structs.MaxTuple/2 {
			return structs.PersonalRelation_CloseFriend
		}
		return structs.PersonalRelation_Friend
	}
}

func AddSkillAndFaith(dice *dg.Dice, mp *MetaPeople, person *structs.Person) ([]*structs.Tuple, []*structs.Tuple) {
	if !dice.IsValidDemographic(person.Race, person.Culture) {
		return []*structs.Tuple{}, []*structs.Tuple{}
	}

	demo := dice.MustDemographic(person.Race, person.Culture)

	skills := demo.RandomProfession(person.ID)
	if len(skills) > 0 {
		mp.Skills = append(mp.Skills, skills...)
		person.PreferredProfession = skills[0].Object
		person.Ethos = person.Ethos.AddEthos(dice.EthosWeightFromProfessions(skills))
	}

	faiths := demo.RandomFaith(person.ID)
	mp.Faith = append(mp.Faith, faiths...)

	return skills, faiths
}

func SiblingRelationship(dice *dg.Demographic, mp *MetaPeople, a, b *structs.Person) {
	if a.ID == b.ID {
		return
	}

	brel := structs.PersonalRelation_Sister
	if b.IsMale {
		brel = structs.PersonalRelation_Brother
	}

	arel := structs.PersonalRelation_Sister
	if a.IsMale {
		arel = structs.PersonalRelation_Brother
	}

	mp.Relations = append(
		mp.Relations,
		&structs.Tuple{Subject: a.ID, Object: b.ID, Value: int64(brel)},
		&structs.Tuple{Subject: b.ID, Object: a.ID, Value: int64(arel)},
	)
	mp.Trust = append(
		mp.Trust,
		&structs.Tuple{Subject: a.ID, Object: b.ID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: b.ID, Object: a.ID, Value: dice.RandomTrust()},
	)
}

// AddParentChildRelations adds parent & grandparent relationships
func AddParentChildRelations(dice *dg.Demographic, mp *MetaPeople, child *structs.Person, f *structs.Family) {
	rel := structs.PersonalRelation_Daughter
	grel := structs.PersonalRelation_Granddaughter
	if child.IsMale {
		rel = structs.PersonalRelation_Son
		grel = structs.PersonalRelation_Grandson
	}

	// parents
	mp.Relations = append(
		mp.Relations,
		&structs.Tuple{Subject: f.FemaleID, Object: child.ID, Value: int64(rel)},
		&structs.Tuple{Subject: f.MaleID, Object: child.ID, Value: int64(rel)},
		&structs.Tuple{Subject: child.ID, Object: f.FemaleID, Value: int64(structs.PersonalRelation_Mother)},
		&structs.Tuple{Subject: child.ID, Object: f.MaleID, Value: int64(structs.PersonalRelation_Father)},
	)
	mp.Trust = append(
		mp.Trust,
		&structs.Tuple{Subject: f.FemaleID, Object: child.ID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: f.MaleID, Object: child.ID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: child.ID, Object: f.FemaleID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: child.ID, Object: f.MaleID, Value: dice.RandomTrust()},
	)

	// grandparents (first generation families wont have grandparents)
	if f.MaGrandmaID != "" {
		mp.Relations = append(
			mp.Relations,
			&structs.Tuple{Subject: f.MaGrandmaID, Object: child.ID, Value: int64(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandmaID, Value: int64(structs.PersonalRelation_Grandmother)},
		)
		mp.Trust = append(
			mp.Trust,
			&structs.Tuple{Subject: f.MaGrandmaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandmaID, Value: dice.RandomTrust()},
		)
	}
	if f.MaGrandpaID != "" {
		mp.Relations = append(
			mp.Relations,
			&structs.Tuple{Subject: f.MaGrandpaID, Object: child.ID, Value: int64(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandpaID, Value: int64(structs.PersonalRelation_Grandfather)},
		)
		mp.Trust = append(
			mp.Trust,
			&structs.Tuple{Subject: f.MaGrandpaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandpaID, Value: dice.RandomTrust()},
		)
	}
	if f.PaGrandmaID != "" {
		mp.Relations = append(
			mp.Relations,
			&structs.Tuple{Subject: f.PaGrandmaID, Object: child.ID, Value: int64(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandmaID, Value: int64(structs.PersonalRelation_Grandmother)},
		)
		mp.Trust = append(
			mp.Trust,
			&structs.Tuple{Subject: f.PaGrandmaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandmaID, Value: dice.RandomTrust()},
		)
	}
	if f.PaGrandpaID != "" {
		mp.Relations = append(
			mp.Relations,
			&structs.Tuple{Subject: f.PaGrandpaID, Object: child.ID, Value: int64(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandpaID, Value: int64(structs.PersonalRelation_Grandfather)},
		)
		mp.Trust = append(
			mp.Trust,
			&structs.Tuple{Subject: f.PaGrandpaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandpaID, Value: dice.RandomTrust()},
		)
	}
}
