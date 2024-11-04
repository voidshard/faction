package controller

import (
	"context"
	"fmt"
	"sync"

	"github.com/voidshard/faction/pkg/structs"
)

type Manager struct {
	structs.APIClient

	name string

	sublock sync.Mutex
	subs    map[string]*subscriber
}

func NewManager(name string, client structs.APIClient) *Manager {
	return &Manager{
		APIClient: client,
		name:      name,
		sublock:   sync.Mutex{},
		subs:      make(map[string]*subscriber),
	}
}

func (m *Manager) Deregsiter(ctx context.Context, ch *structs.Change) error {
	key := changeKey(ch)

	m.sublock.Lock()
	defer m.sublock.Unlock()
	sub, ok := m.subs[key]
	if !ok {
		return nil // not registered
	}

	delete(m.subs, key)
	return sub.stream.CloseSend()
}

func (m *Manager) Register(ctx context.Context, ch *structs.Change) error {
	key := changeKey(ch)

	m.sublock.Lock()
	defer m.sublock.Unlock()
	_, ok := m.subs[key]
	if ok {
		return nil // already registered
	}

	stream, err := m.OnChange(ctx, &structs.OnChangeRequest{Data: ch, Queue: m.name})
	if err != nil {
		return err
	}

	sub := newSubscriber(stream)
	m.subs[key] = sub
	return nil
}

func changeKey(ch *structs.Change) string {
	return fmt.Sprintf("%s:%s:%d:%s", ch.World, ch.Area, ch.Key, ch.Id)
}
