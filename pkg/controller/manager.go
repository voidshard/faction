package controller

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/voidshard/faction/pkg/structs/v1"
	"github.com/voidshard/faction/pkg/util/lock"
	"github.com/voidshard/faction/pkg/util/log"

	"google.golang.org/grpc"
)

var (
	defaultLockTTL = 1 * time.Minute
)

type Manager struct {
	v1.APIClient

	name       string
	controller Controller
	locker     lock.Locker

	l log.Logger

	streamlock sync.Mutex
	streams    map[string]*stream
	events     chan *v1.OnChangeResponse
	ack        grpc.ClientStreamingClient[v1.AckRequest, v1.AckResponse]
}

func NewManager(name string, client v1.APIClient, locker lock.Locker) *Manager {
	if locker == nil {
		locker = &nullLocker{}
	}
	return &Manager{
		APIClient:  client,
		name:       name,
		locker:     locker,
		l:          log.Sublogger(name),
		streamlock: sync.Mutex{},
		streams:    make(map[string]*stream),
		events:     make(chan *v1.OnChangeResponse),
	}
}

func (m *Manager) close() error {
	m.streamlock.Lock()
	defer m.streamlock.Unlock()
	for _, st := range m.streams {
		st.Close()
	}
	m.ack.CloseSend()
	close(m.events)
	return nil
}

func (m *Manager) Run(ctrl Controller) error {
	m.controller = ctrl

	err := ctrl.Init(m)
	if err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	defer m.close()

	for {
		select {
		case sig <- sigs:
			m.l.Info().Str("signal", sig.String()).Msg("received signal")
			return nil
		case req, ok := <-m.events:
			if !ok {
				continue
			}
			ctx := context.Background()
			pan := log.NewSpan(ctx, "controller.Event", map[string]interface{}{
				"Controller": m.name,
				"World":      req.Data.World,
				"Area":       req.Data.Area,
				"Key":        req.Data.Key,
				"Id":         req.Data.Id,
				"Ack":        req.Ack,
			})

			key := changeKey(req.Data)
			err := m.locker.Lock(key, defaultLockTTL)
			if err != nil {
				m.l.Warn().Err(err).Msg("failed to lock event")
				pan.Err(err)
				pan.End()
				continue
			}

			ok, err := ctrl.Handle(ctx, ch)

			uerr := m.locker.Unlock(key)
			if uerr != nil {
				m.l.Warn().Err(uerr).Msg("failed to unlock event")
				// we will try to keep going
			}

			if err != nil {
				m.l.Warn().Err(err).Msg("error handling event")
				pan.Err(err)
				pan.End()
				continue
			}
			if !ok {
				pan.End()
				continue
			}

			err = m.ack.Send(&v1.AckRequest{Ack: []string{req.Ack}})
			if err != nil {
				m.l.Warn().Err(err).Msg("error sending ack")
				pan.Err(err)
			}
			pan.End()
		}
	}
}

func (m *Manager) Deregister(ch *v1.Change) error {
	key := changeKey(ch)

	m.streamlock.Lock()
	defer m.streamlock.Unlock()
	st, ok := m.streams[key]
	if !ok {
		return nil // not registered
	}

	delete(m.streams, key)
	return st.Close()
}

func (m *Manager) Register(ch *v1.Change) error {
	key := changeKey(ch)

	m.streamlock.Lock()
	defer m.streamlock.Unlock()
	_, ok := m.streams[key]
	if ok {
		return nil // already registered
	}

	if m.ack == nil {
		out, err := client.AckStream(context.Background())
		if err != nil {
			return err
		}
		m.ack = out
	}

	m.streams[key] = newStream(m.name, m, ch)
	return nil
}

func changeKey(ch *v1.Change) string {
	return fmt.Sprintf("%s_%s_%d_%s", ch.World, ch.Area, ch.Key, ch.Id)
}
