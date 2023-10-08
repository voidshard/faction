package base

import (
	"fmt"

	"github.com/voidshard/faction/pkg/queue"
	"github.com/voidshard/faction/pkg/structs"

	mapset "github.com/deckarep/golang-set/v2"
)

func eventTask(e structs.EventType) string {
	return fmt.Sprintf("event:%s", e)
}

func (s *Base) registerTasksWithQueue() error {
	s.queue.Register(eventTask(structs.EventPersonBirth), s.handlerOnEventPersonBirth)
	s.queue.Register(eventTask(structs.EventPersonDeath), s.handlerOnEventPersonDeath)
	s.queue.Register(eventTask(structs.EventFamilyMarriage), s.handlerOnEventFamilyMarriage)
	s.queue.Register(eventTask(structs.EventFactionCreated), s.handlerOnEventFactionCreated)
	return nil
}

// toEvents returns a slice of events from a slice of jobmeta.
// We also return a slice of unique subjects and causes.
func (s *Base) toEvents(in []*queue.JobMeta) ([]*structs.Event, error) {
	out := make([]*structs.Event, len(in))
	for i, j := range in {
		evt := &structs.Event{}
		if err := evt.UnmarshalJson(j.Job.Args); err != nil {
			return nil, err
		}
		out[i] = evt
	}
	return out, nil
}

func eventSubjects(in []*structs.Event) []string {
	out := mapset.NewSet[string]()
	for _, e := range in {
		out.Add(string(e.SubjectMetaVal))
	}
	return out.ToSlice()
}

func (s *Base) handlerOnEventFamilyMarriage(in ...*queue.JobMeta) error {
	events, err := s.toEvents(in)
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

func (s *Base) handlerOnEventFactionCreated(in ...*queue.JobMeta) error {
	events, err := s.toEvents(in)
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

func (s *Base) handlerOnEventPersonBirth(in ...*queue.JobMeta) error {
	events, err := s.toEvents(in)
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

func (s *Base) handlerOnEventPersonDeath(in ...*queue.JobMeta) error {
	events, err := s.toEvents(in)
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
