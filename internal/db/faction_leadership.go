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

	People map[string]*structs.Person
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
		People:      map[string]*structs.Person{},
	}
}

func (l *FactionLeadership) Get(index int) string {
	// we only intend to return a small number of entries here, so in theory
	// this shouldn't be too bad.
	// Obviously if the faction has thousands of people and an int is requested for
	// some random appentice that's .. not great.
	// But this internal struct is simply for looking at faction leaders .. which
	// should be a small number of people.
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

func (l *FactionLeadership) Add(rank structs.FactionRank, p *structs.Person) {
	l.People[p.ID] = p
	l.Total++
	switch rank {
	case structs.FactionRank_Ruler:
		l.Ruler = append(l.Ruler, p.ID)
	case structs.FactionRank_Elder:
		l.Elder = append(l.Elder, p.ID)
	case structs.FactionRank_GrandMaster:
		l.GrandMaster = append(l.GrandMaster, p.ID)
	case structs.FactionRank_Master:
		l.Master = append(l.Master, p.ID)
	case structs.FactionRank_Expert:
		l.Expert = append(l.Expert, p.ID)
	case structs.FactionRank_Adept:
		l.Adept = append(l.Adept, p.ID)
	case structs.FactionRank_Journeyman:
		l.Journeyman = append(l.Journeyman, p.ID)
	case structs.FactionRank_Novice:
		l.Novice = append(l.Novice, p.ID)
	case structs.FactionRank_Apprentice:
		l.Apprentice = append(l.Apprentice, p.ID)
	case structs.FactionRank_Associate:
		l.Apprentice = append(l.Apprentice, p.ID)
	}
}
