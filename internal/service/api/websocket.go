package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/voidshard/faction/pkg/structs/api"
	v1 "github.com/voidshard/faction/pkg/structs/v1"
	"github.com/voidshard/faction/pkg/util/log"
)

// WebSocket code here heavily inspired by the Gorilla Websocket Chat example
// https://github.com/gorilla/websocket/blob/main/examples/chat/client.go
// https://github.com/gorilla/websocket/blob/main/examples/chat/hub.go

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
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type WebSocket struct {
	// The service that this websocket is connected to.
	svc *Service

	// The websocket connection.
	conn *websocket.Conn

	// channel of events from service
	events <-chan *v1.Event

	// channel to kill events subscription
	killEvents chan<- bool
}

func newWebSocket(svc *Service, conn *websocket.Conn) *WebSocket {
	return &WebSocket{svc: svc, conn: conn}
}

// Close the connection, attempt to send an error message to the client
// explaining why the connection is being closed if provided.
func (c *WebSocket) Close(errMsg *api.ErrorResponse) {
	c.killEvents <- true
	c.conn.Close()
	c.conn = nil

	if errMsg == nil {
		return
	}

	data, err := json.Marshal(errMsg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal error message")
		return
	}

	c.conn.SetWriteDeadline(time.Now().Add(writeWait))

	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get writer")
		return
	}

	w.Write([]byte("ERROR"))
	w.Write(data)
	w.Write(newline)

	err = w.Close()
	log.Debug().Err(err).Msg("Closed connection")
}

func (c *WebSocket) Pump(events <-chan *v1.Event, killEvents chan<- bool) {
	c.events = events
	c.killEvents = killEvents

	go c.writePump()
	go c.readPump()
}

// readPump reads messages from the websocket connection and enacts them.
func (c *WebSocket) readPump() {
	// Messages we expect from the client:
	// - AckId (str): acknowledge that the event has been processed
	defer func() {
		c.killEvents <- true
		c.conn.Close()
		c.conn = nil
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		pan := log.NewSpan(context.Background(), "websocket.Ack")

		if err != nil {
			log.Error().Err(err).Msg("Failed to read message from websocket")
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				pan.Err(err)
				pan.End()
			}
			pan.End()
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if len(message) == 0 {
			pan.Err(fmt.Errorf("empty message"))
			pan.End()
			continue
		}
		err = c.svc.ackEvent(string(message))
		if err != nil {
			pan.Err(err)
		}
		pan.SetAttributes(map[string]interface{}{"AckId": string(message)})
		log.Debug().Err(err).Str("AckId", string(message)).Msg("Acked event")
		pan.End()
	}
}

// writePump sends messages from the Events channel to the websocket connection.
func (c *WebSocket) writePump() {
	// Messages we send to the client:
	// - Event (json): event to send to client
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case event, ok := <-c.events:
			message, err := json.Marshal(event)
			pan := log.NewSpan(context.Background(), "websocket.Event")
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal event")
				pan.Err(err)
				pan.End()
				continue
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				pan.End()
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get writer")
				pan.Err(err)
				pan.End()
				return
			}
			w.Write(message)
			w.Write(newline)

			err = w.Close()
			if err != nil {
				log.Error().Err(err).Msg("Failed to close writer")
				pan.Err(err)
			}
			pan.End()
		case <-ticker.C:
			if c.conn == nil {
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			log.Debug().Err(err).Msg("Ping")
		}
	}
}
