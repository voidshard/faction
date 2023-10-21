package base

import (
	"bytes"
	"fmt"

	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/pkg/queue"
	"github.com/voidshard/faction/pkg/structs"

	mapset "github.com/deckarep/golang-set/v2"
)

var (
	keyEventType = "event_type"
)

func eventTask(e structs.EventType) string {
	return fmt.Sprintf("event-%s", e)
}

func (s *Base) registerTasksWithQueue() error {
	errs := []error{
		s.queue.Register(eventTask(structs.EventPersonBirth), s.handlerOnEventPersonBirth),
		s.queue.Register(eventTask(structs.EventPersonDeath), s.handlerOnEventPersonDeath),
		s.queue.Register(eventTask(structs.EventFamilyMarriage), s.handlerOnEventFamilyMarriage),
		s.queue.Register(eventTask(structs.EventFactionCreated), s.handlerOnEventFactionCreated),
	}
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// toEvents returns a slice of events from a slice of jobmeta.
// We also return a slice of unique subjects and causes.
func (s *Base) toEvents(in []*queue.Job) ([]*structs.Event, error) {
	out := make([]*structs.Event, len(in))
	for i, j := range in {
		evt := &structs.Event{}
		if err := evt.UnmarshalJson(bytes.TrimSpace(j.Args)); err != nil {
			log.Debug().Str("payload", string(j.Args)).Err(err).Msg()("failed to unmarshal event")
			return nil, err
		}
		out[i] = evt
	}
	return out, nil
}

func eventSubjects(in []*structs.Event) []string {
	out := mapset.NewSet[string]()
	for _, e := range in {
		if e.SubjectMetaKey != "" {
			out.Add(string(e.SubjectMetaVal))
		}
	}
	return out.ToSlice()
}

func (s *Base) handlerOnEventFamilyMarriage(in ...*queue.Job) error {
	events, err := s.toEvents(in)
	log.Debug().Str(keyEventType, string(structs.EventFamilyMarriage)).Err(err).Int("events", len(in)).Msg()("handling event")
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}

	tick, err := s.dbconn.Tick()
	if err != nil {
		return err
	}
	return s.applyFamilyMarriage(tick, events)
}

func (s *Base) handlerOnEventFactionCreated(in ...*queue.Job) error {
	events, err := s.toEvents(in)
	log.Debug().Str(keyEventType, string(structs.EventFactionCreated)).Err(err).Int("events", len(in)).Msg()("handling event")
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}

	tick, err := s.dbconn.Tick()
	if err != nil {
		return err
	}

	return s.applyFactionCreated(tick, events)
}

func (s *Base) handlerOnEventPersonBirth(in ...*queue.Job) error {
	events, err := s.toEvents(in)
	log.Debug().Str(keyEventType, string(structs.EventPersonBirth)).Err(err).Int("events", len(in)).Msg()("handling event")
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}

	tick, err := s.dbconn.Tick()
	if err != nil {
		return err
	}

	return s.applyBirthSiblingRelations(tick, events)
}

func (s *Base) handlerOnEventPersonDeath(in ...*queue.Job) error {
	events, err := s.toEvents(in)
	log.Debug().Str(keyEventType, string(structs.EventPersonDeath)).Err(err).Int("events", len(in)).Msg()("handling event")
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}

	tick, err := s.dbconn.Tick()
	if err != nil {
		return err
	}

	return s.applyDeathFamilyEffect(tick, events)
}
