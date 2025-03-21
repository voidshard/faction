package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/voidshard/faction/pkg/structs/api"
	v1 "github.com/voidshard/faction/pkg/structs/v1"
	"github.com/voidshard/faction/pkg/util/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	errmsg  = []byte("ERROR")
)

type EventStream struct {
	url *url.URL

	mu     sync.RWMutex
	conn   *websocket.Conn
	killed bool

	events chan *v1.Event
	acks   chan string
}

func newEventStream(u *url.URL) (*EventStream, error) {
	me := &EventStream{url: u, events: make(chan *v1.Event), acks: make(chan string), mu: sync.RWMutex{}}
	conn, err := connect(u)
	if err != nil {
		return me, err
	}
	me.conn = conn
	go me.readPump()
	go me.writePump()
	return me, nil
}

func (e *EventStream) connectForever(u *url.URL) *websocket.Conn {
	for {
		if e.killed {
			return nil
		}
		conn, err := connect(u)
		if err == nil {
			return conn
		}
		log.Warn().Err(err).Msg("Failed to connect")
		time.Sleep(1 * time.Second)
	}
}

func (e *EventStream) readPump() {
	defer close(e.events)
	for {
		if e.killed {
			return
		}
		e.mu.RLock()
		_, message, err := e.conn.ReadMessage()
		e.mu.RUnlock()
		if err != nil {
			if !e.killed { // if we're not shutting down, reconnect
				log.Warn().Err(err).Msg("Unexpected websocket error, reconnecting")
				e.mu.Lock()
				// We hold the lock until we reconnect; anything waiting for the
				// connection has to wait anyway (since .. we can't talk to the server)
				// this simply forces our Write routine to block until we're online.
				e.conn = e.connectForever(e.url)
				e.mu.Unlock()
				continue
			}
			log.Debug().Err(err).Msg("Websocket closed")
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		if bytes.HasPrefix(message, errmsg) {
			message = bytes.TrimPrefix(message, errmsg)
			errresp := &api.ErrorResponse{}
			err = json.Unmarshal(message, errresp)
			if err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal error from server")
			} else {
				log.Error().Int("code", errresp.Code).Str("message", errresp.Message).Msg("Received error")
			}
			continue
		}

		evt := &v1.Event{}
		err = json.Unmarshal(message, evt)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal event")
			continue
		}
		e.events <- evt
	}
}

func (e *EventStream) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case id := <-e.acks:
			e.mu.RLock()
			e.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := e.conn.NextWriter(websocket.TextMessage)
			e.mu.RUnlock()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get writer")
				return
			}
			w.Write([]byte(fmt.Sprintf("%s\n", id)))
			log.Debug().Err(w.Close()).Msg("Acked event")
		case <-ticker.C:
			e.mu.RLock()
			e.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := e.conn.WriteMessage(websocket.PingMessage, nil)
			e.mu.RUnlock()
			log.Debug().Err(err).Msg("Ping")
		}
	}
}

func (e *EventStream) Close() {
	e.killed = true
	e.conn.Close()
}

func (e *EventStream) Events() <-chan *v1.Event {
	return e.events
}

func (e *EventStream) Ack(id string) error {
	e.acks <- id
	return nil
}

func connect(u *url.URL) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Warn().Err(err).Str("url", u.String()).Msg("Failed to connect websocket")
		return nil, err
	}
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		log.Debug().Msg("Pong received")
		return nil
	})
	return conn, nil
}
