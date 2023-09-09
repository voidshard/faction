package db

import (
	"github.com/voidshard/faction/pkg/structs"
)

type FactionLeadership struct {
	Ruler       []string
	Elder       []string
	GrandMaster []string
	Master      []string
	Expert      []string
	Adept       []string
	Journeyman  []string
	Novice      []string
	Apprentice  []string

	Total int
}

func NewFactionLeadership() *FactionLeadership {
	return &FactionLeadership{
		Ruler:       []string{},
		Elder:       []string{},
		GrandMaster: []string{},
		Master:      []string{},
		Expert:      []string{},
		Adept:       []string{},
		Journeyman:  []string{},
		Novice:      []string{},
		Apprentice:  []string{},
	}
}

func (l *FactionLeadership) Get(index int) string {
	if index < 0 || index >= l.Total {
		return ""
	}
	if index < len(l.Ruler) {
		return l.Ruler[index]
	}
	index -= len(l.Ruler)
	if index < len(l.Elder) {
		return l.Elder[index]
	}
	index -= len(l.Elder)
	if index < len(l.GrandMaster) {
		return l.GrandMaster[index]
	}
	index -= len(l.GrandMaster)
	if index < len(l.Master) {
		return l.Master[index]
	}
	index -= len(l.Master)
	if index < len(l.Expert) {
		return l.Expert[index]
	}
	index -= len(l.Expert)
	if index < len(l.Adept) {
		return l.Adept[index]
	}
	index -= len(l.Adept)
	if index < len(l.Journeyman) {
		return l.Journeyman[index]
	}
	index -= len(l.Journeyman)
	if index < len(l.Novice) {
		return l.Novice[index]
	}
	index -= len(l.Novice)
	if index < len(l.Apprentice) {
		return l.Apprentice[index]
	}
	return ""
}

func (l *FactionLeadership) Add(rank structs.FactionRank, personID string) {
	l.Total++
	switch rank {
	case structs.FactionRankRuler:
		l.Ruler = append(l.Ruler, personID)
	case structs.FactionRankElder:
		l.Elder = append(l.Elder, personID)
	case structs.FactionRankGrandMaster:
		l.GrandMaster = append(l.GrandMaster, personID)
	case structs.FactionRankMaster:
		l.Master = append(l.Master, personID)
	case structs.FactionRankExpert:
		l.Expert = append(l.Expert, personID)
	case structs.FactionRankAdept:
		l.Adept = append(l.Adept, personID)
	case structs.FactionRankJourneyman:
		l.Journeyman = append(l.Journeyman, personID)
	case structs.FactionRankNovice:
		l.Novice = append(l.Novice, personID)
	case structs.FactionRankApprentice:
		l.Apprentice = append(l.Apprentice, personID)
	case structs.FactionRankAssociate:
		l.Apprentice = append(l.Apprentice, personID)
	}
}
