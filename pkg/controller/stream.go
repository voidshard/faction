package controller

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/voidshard/faction/pkg/structs/v1"
	"github.com/voidshard/faction/util/log"
)

type stream struct {
	name    string
	l       log.Logger
	mgr     *Manager
	kill    chan bool
	watch   *v1.Change
	fromAPI v1.API_OnChangeClient
	closed  bool
}

func newStream(name string, mgr *Manager, watch *v1.Change) *stream {
	me := &stream{
		name:  name,
		l:     log.Sublogger(name),
		kill:  make(chan bool),
		watch: watch,
	}
	go me.readPump()
	return me
}

func (s *stream) connect() {
	queueName := fmt.Sprintf("%s_%s", s.mgr.name, s.name)
	attempt := -1
	for {
		attempt++
		in, err := m.OnChange(
			context.Background(),
			&v1.OnChangeRequest{Data: s.watch, Queue: queueName},
		)
		if err != nil {
			s.l.Warn().Int("Attempt", attempt).Err(err).Msg("error connecting to watch stream")
			if attempt > 0 {
				time.Sleep(2 * time.Second)
			}
			continue
		}
		s.fromAPI = in
	}
}

func (s *stream) readPump() {
	defer close(s.kill)

	s.connect() // no error, because we retry forever

	for {
		select {
		case <-s.kill:
			l.Info().Msg("watch stream closed")
			return
		default:
			resp, err := s.fromAPI.Recv()
			if err == io.EOF {
				if s.closed {
					return // we've been asked to close
				}
				l.Warn().Err(err).Msg("watch stream closed, reconnecting")
				s.connect() // reconnect
				continue
			} else if err != nil {
				l.Warn().Err(err).Msg("error reading watch stream response")
				continue
			} else if resp == nil {
				l.Warn().Msg("watch stream got nil response")
				continue // ???
			} else if resp.Error != nil {
				l.Warn().Str("error", resp.Error.Message).Msg("server sent error on watch stream")
				continue // error sent from server
			}
			s.events <- resp
		}
	}
}

func (s *stream) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	s.kill <- true
	return s.fromAPI.CloseSend()
}
