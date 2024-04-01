package simutil

import (
	"github.com/voidshard/faction/pkg/structs"
)

func NewJobPendingEvent(j *structs.Job) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventType_JobPending,
		Tick:           j.TickCreated,
		SubjectMetaKey: structs.Meta_KeyJob,
		SubjectMetaVal: j.ID,
	}
}

func NewFactionCreatedEvent(f *structs.Faction, tick int64) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventType_FactionCreated,
		Tick:           tick,
		SubjectMetaKey: structs.Meta_KeyFaction,
		SubjectMetaVal: f.ID,
	}
}

func NewFactionPromotionEvent(p *structs.Person, tick int64, factionID string) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventType_FactionPromotion,
		Tick:           tick,
		SubjectMetaKey: structs.Meta_KeyPerson,
		SubjectMetaVal: p.ID,
		CauseMetaKey:   structs.Meta_KeyFaction,
		CauseMetaVal:   factionID,
	}
}

func NewFactionChangeEvent(p *structs.Person, tick int64, previousPreferredFaction string) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventType_PersonChangeFaction,
		Tick:           tick,
		SubjectMetaKey: structs.Meta_KeyPerson,
		SubjectMetaVal: p.ID,
		CauseMetaKey:   structs.Meta_KeyFaction,
		CauseMetaVal:   previousPreferredFaction,
	}
}

func NewMarriageEvent(f *structs.Family, tick int64) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventType_FamilyMarriage,
		Tick:           tick,
		SubjectMetaKey: structs.Meta_KeyFamily,
		SubjectMetaVal: f.ID,
	}
}

func NewDivorceEvent(f *structs.Family, tick int64) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventType_FamilyDivorce,
		Tick:           tick,
		SubjectMetaKey: structs.Meta_KeyFamily,
		SubjectMetaVal: f.ID,
	}
}

func NewBirthEvent(p *structs.Person, familiyID string) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventType_PersonBirth,
		Tick:           p.BirthTick,
		SubjectMetaKey: structs.Meta_KeyPerson,
		SubjectMetaVal: p.ID,
		CauseMetaKey:   structs.Meta_KeyFamily,
		CauseMetaVal:   familiyID,
	}
}

func NewDeathEvent(p *structs.Person) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventType_PersonDeath,
		Tick:           p.DeathTick,
		SubjectMetaKey: structs.Meta_KeyPerson,
		SubjectMetaVal: p.ID,
		CauseMetaKey:   p.DeathMetaKey,
		CauseMetaVal:   p.DeathMetaVal,
	}
}

func NewDeathEventWithCause(tick int64, subject string, causek structs.Meta, causev string) *structs.Event {
	return &structs.Event{
		ID:             structs.NewID(),
		Type:           structs.EventType_PersonDeath,
		Tick:           tick,
		SubjectMetaKey: structs.Meta_KeyPerson,
		SubjectMetaVal: subject,
		CauseMetaKey:   causek,
		CauseMetaVal:   causev,
	}
}
