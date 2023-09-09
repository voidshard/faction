package simutil

import (
	"github.com/voidshard/faction/pkg/structs"
)

func NewJobPending(j *structs.Job) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Tick:           j.TickCreated,
		SubjectMetaKey: structs.MetaKeyJob,
		SubjectMetaVal: j.ID,
	}
}

func NewFactionPromotionEvent(p *structs.Person, tick int, factionID string) *structs.Event {
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

func NewFactionChangeEvent(p *structs.Person, tick int, previousPreferredFaction string) *structs.Event {
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

func NewMarriageEvent(f *structs.Family, tick int) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventFamilyMarriage,
		Tick:           tick,
		SubjectMetaKey: structs.MetaKeyFamily,
		SubjectMetaVal: f.ID,
	}
}

func NewDivorceEvent(f *structs.Family, tick int) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventFamilyDivorce,
		Tick:           tick,
		SubjectMetaKey: structs.MetaKeyFamily,
		SubjectMetaVal: f.ID,
	}
}

func NewBirthEvent(p *structs.Person, familiyID string) *structs.Event {
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

func NewDeathEvent(p *structs.Person) *structs.Event {
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

func NewDeathEventWithCause(tick int, subject string, causek structs.MetaKey, causev string) *structs.Event {
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
