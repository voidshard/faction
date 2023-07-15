package base

import (
	dg "github.com/voidshard/faction/internal/random/demographics"
	"github.com/voidshard/faction/pkg/structs"
)

func (s *Base) addSkillAndFaith(demo *dg.Demographic, mp *metaPeople, person *structs.Person) ([]*structs.Tuple, []*structs.Tuple) {
	skills := demo.RandomProfession(person.ID)
	if len(skills) > 0 {
		mp.skills = append(mp.skills, skills...)
		person.PreferredProfession = skills[0].Object
		person.Ethos = *person.Ethos.AddEthos(s.dice.EthosWeightFromProfessions(skills))
	}

	faiths := demo.RandomFaith(person.ID)
	mp.faith = append(mp.faith, faiths...)

	return skills, faiths
}

func siblingRelationship(dice *dg.Demographic, mp *metaPeople, a, b *structs.Person) {
	brel := structs.PersonalRelationSister
	if b.IsMale {
		brel = structs.PersonalRelationBrother
	}

	arel := structs.PersonalRelationSister
	if a.IsMale {
		arel = structs.PersonalRelationBrother
	}

	mp.relations = append(
		mp.relations,
		&structs.Tuple{Subject: a.ID, Object: b.ID, Value: int(brel)},
		&structs.Tuple{Subject: b.ID, Object: a.ID, Value: int(arel)},
	)
	mp.trust = append(
		mp.trust,
		&structs.Tuple{Subject: a.ID, Object: b.ID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: b.ID, Object: a.ID, Value: dice.RandomTrust()},
	)
}

func addChildToFamily(dice *dg.Demographic, child *structs.Person, f *structs.Family) {
	child.Ethos = *structs.EthosAverage(&child.Ethos, &f.Ethos)
	child.BirthFamilyID = f.ID
	child.Race = f.Race
	child.Culture = f.Culture
	f.NumberOfChildren++
	if f.NumberOfChildren >= dice.MaxFamilySize() {
		f.IsChildBearing = false
	}
}

// addParentChildRelations adds parent & grandparent relationships
func addParentChildRelations(dice *dg.Demographic, mp *metaPeople, child *structs.Person, f *structs.Family) {
	rel := structs.PersonalRelationDaughter
	grel := structs.PersonalRelationGranddaughter
	if child.IsMale {
		rel = structs.PersonalRelationSon
		grel = structs.PersonalRelationGrandson
	}

	// parents
	mp.relations = append(
		mp.relations,
		&structs.Tuple{Subject: f.FemaleID, Object: child.ID, Value: int(rel)},
		&structs.Tuple{Subject: f.MaleID, Object: child.ID, Value: int(rel)},
		&structs.Tuple{Subject: child.ID, Object: f.FemaleID, Value: int(structs.PersonalRelationMother)},
		&structs.Tuple{Subject: child.ID, Object: f.MaleID, Value: int(structs.PersonalRelationFather)},
	)
	mp.trust = append(
		mp.trust,
		&structs.Tuple{Subject: f.FemaleID, Object: child.ID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: f.MaleID, Object: child.ID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: child.ID, Object: f.FemaleID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: child.ID, Object: f.MaleID, Value: dice.RandomTrust()},
	)

	// grandparents (first generation families wont have grandparents)
	if f.MaGrandmaID != "" {
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: f.MaGrandmaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandmaID, Value: int(structs.PersonalRelationGrandmother)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: f.MaGrandmaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandmaID, Value: dice.RandomTrust()},
		)
	}
	if f.MaGrandpaID != "" {
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: f.MaGrandpaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandpaID, Value: int(structs.PersonalRelationGrandfather)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: f.MaGrandpaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandpaID, Value: dice.RandomTrust()},
		)
	}
	if f.PaGrandmaID != "" {
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: f.PaGrandmaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandmaID, Value: int(structs.PersonalRelationGrandmother)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: f.PaGrandmaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandmaID, Value: dice.RandomTrust()},
		)
	}
	if f.PaGrandpaID != "" {
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: f.PaGrandpaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandpaID, Value: int(structs.PersonalRelationGrandfather)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: f.PaGrandpaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandpaID, Value: dice.RandomTrust()},
		)
	}
}
