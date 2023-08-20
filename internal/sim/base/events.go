package base

import (
	"github.com/voidshard/faction/pkg/structs"
)

func newFactionPromotionEvent(p *structs.Person, tick int, factionID string) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventFactionPromotion,
		Tick:           tick,
		SubjectMetaKey: structs.MetaKeyPerson,
		SubjectMetaVal: p.ID,
		CauseMetaKey:   structs.MetaKeyFaction,
		CauseMetaVal:   factionID,
	}
}

func newFactionChangeEvent(p *structs.Person, tick int, previousPreferredFaction string) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventPersonChangeFaction,
		Tick:           tick,
		SubjectMetaKey: structs.MetaKeyPerson,
		SubjectMetaVal: p.ID,
		CauseMetaKey:   structs.MetaKeyFaction,
		CauseMetaVal:   previousPreferredFaction,
	}
}

func newMarriageEvent(f *structs.Family, tick int) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventFamilyMarriage,
		Tick:           tick,
		SubjectMetaKey: structs.MetaKeyFamily,
		SubjectMetaVal: f.ID,
	}
}

func newDivorceEvent(f *structs.Family, tick int) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventFamilyDivorce,
		Tick:           tick,
		SubjectMetaKey: structs.MetaKeyFamily,
		SubjectMetaVal: f.ID,
	}
}

func newBirthEvent(p *structs.Person, familiyID string) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventPersonBirth,
		Tick:           p.BirthTick,
		SubjectMetaKey: structs.MetaKeyPerson,
		SubjectMetaVal: p.ID,
		CauseMetaKey:   structs.MetaKeyFamily,
		CauseMetaVal:   familiyID,
	}
}

func newDeathEvent(p *structs.Person) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventPersonDeath,
		Tick:           p.DeathTick,
		SubjectMetaKey: structs.MetaKeyPerson,
		SubjectMetaVal: p.ID,
		CauseMetaKey:   p.DeathMetaKey,
		CauseMetaVal:   p.DeathMetaVal,
	}
}

func newDeathEventWithCause(tick int, subject string, causek structs.MetaKey, causev string) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventPersonDeath,
		Tick:           tick,
		SubjectMetaKey: structs.MetaKeyPerson,
		SubjectMetaVal: subject,
		CauseMetaKey:   causek,
		CauseMetaVal:   causev,
	}
}
