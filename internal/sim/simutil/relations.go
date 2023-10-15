package simutil

import (
	dg "github.com/voidshard/faction/internal/random/demographics"
	"github.com/voidshard/faction/pkg/structs"
)

func AddChildToFamily(dice *dg.Demographic, child *structs.Person, f *structs.Family) {
	child.Ethos = *structs.EthosAverage(&child.Ethos, &f.Ethos)
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
		return structs.PersonalRelationColleague
	} else if a < 0 {
		if a < structs.MinTuple/2 {
			return structs.PersonalRelationHatedEnemy
		}
		return structs.PersonalRelationEnemy
	} else { // a > 0
		if a > structs.MaxTuple/2 {
			return structs.PersonalRelationCloseFriend
		}
		return structs.PersonalRelationFriend
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
		person.Ethos = *person.Ethos.AddEthos(dice.EthosWeightFromProfessions(skills))
	}

	faiths := demo.RandomFaith(person.ID)
	mp.Faith = append(mp.Faith, faiths...)

	return skills, faiths
}

func SiblingRelationship(dice *dg.Demographic, mp *MetaPeople, a, b *structs.Person) {
	if a.ID == b.ID {
		return
	}

	brel := structs.PersonalRelationSister
	if b.IsMale {
		brel = structs.PersonalRelationBrother
	}

	arel := structs.PersonalRelationSister
	if a.IsMale {
		arel = structs.PersonalRelationBrother
	}

	mp.Relations = append(
		mp.Relations,
		&structs.Tuple{Subject: a.ID, Object: b.ID, Value: int(brel)},
		&structs.Tuple{Subject: b.ID, Object: a.ID, Value: int(arel)},
	)
	mp.Trust = append(
		mp.Trust,
		&structs.Tuple{Subject: a.ID, Object: b.ID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: b.ID, Object: a.ID, Value: dice.RandomTrust()},
	)
}

// AddParentChildRelations adds parent & grandparent relationships
func AddParentChildRelations(dice *dg.Demographic, mp *MetaPeople, child *structs.Person, f *structs.Family) {
	rel := structs.PersonalRelationDaughter
	grel := structs.PersonalRelationGranddaughter
	if child.IsMale {
		rel = structs.PersonalRelationSon
		grel = structs.PersonalRelationGrandson
	}

	// parents
	mp.Relations = append(
		mp.Relations,
		&structs.Tuple{Subject: f.FemaleID, Object: child.ID, Value: int(rel)},
		&structs.Tuple{Subject: f.MaleID, Object: child.ID, Value: int(rel)},
		&structs.Tuple{Subject: child.ID, Object: f.FemaleID, Value: int(structs.PersonalRelationMother)},
		&structs.Tuple{Subject: child.ID, Object: f.MaleID, Value: int(structs.PersonalRelationFather)},
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
			&structs.Tuple{Subject: f.MaGrandmaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandmaID, Value: int(structs.PersonalRelationGrandmother)},
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
			&structs.Tuple{Subject: f.MaGrandpaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandpaID, Value: int(structs.PersonalRelationGrandfather)},
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
			&structs.Tuple{Subject: f.PaGrandmaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandmaID, Value: int(structs.PersonalRelationGrandmother)},
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
			&structs.Tuple{Subject: f.PaGrandpaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandpaID, Value: int(structs.PersonalRelationGrandfather)},
		)
		mp.Trust = append(
			mp.Trust,
			&structs.Tuple{Subject: f.PaGrandpaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandpaID, Value: dice.RandomTrust()},
		)
	}
}
