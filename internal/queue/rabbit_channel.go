package queue

import (
	"errors"
	"fmt"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/util/uuid"
)

// rabbitChannel wraps a connection and channel to rabbitmq
// and provides a way to reconnect if the connection is lost
type rabbitChannel struct {
	id string

	cfg      *RabbitConfig
	log      log.Logger
	prefetch int

	closed  bool
	conn    *amqp.Connection
	channel *amqp.Channel

	onReconnect chan *amqp.Channel
}

func newRabbitChannel(cfg *RabbitConfig, name string, prefetch int) (*rabbitChannel, error) {
	id := uuid.New()
	rc := &rabbitChannel{
		id:       id,
		cfg:      cfg,
		log:      log.Sublogger("rabbit-channel", map[string]interface{}{"name": name, "id": id}),
		prefetch: prefetch,
	}
	return rc, rc.connectAndMonitor()
}

func (rc *rabbitChannel) Close() {
	rc.log.Info().Msg("Close called by user, killing rabbitmq channel")
	rc.closed = true
	rc.close()
}

func (rc *rabbitChannel) close() {
	rc.log.Debug().Msg("Internal close called, closing rabbitmq channel")
	if rc.channel != nil {
		rc.channel.Close()
	}
	if rc.conn != nil {
		rc.conn.Close()
	}
}

// Channel returns the current channel.
// Nb. this channel may be in the process of connecting / reconnecting
func (rc *rabbitChannel) Channel() *amqp.Channel {
	return rc.channel
}

// ChannelStream returns a channel that will emit the current (rabbit) channel,
// as it changes (eg. reconnects occur).
func (rc *rabbitChannel) ChannelStream(notify chan *amqp.Channel) <-chan *amqp.Channel {
	rc.onReconnect = notify
	return rc.onReconnect
}

func isConnectError(err error) bool {
	if err == nil {
		return false
	} else if errors.Is(err, amqp.ErrClosed) {
		return true
	} else if strings.Contains(err.Error(), "connection reset") {
		return true
	} else if strings.Contains(err.Error(), "broken pipe") {
		return true
	} else if strings.Contains(err.Error(), "connection closure") {
		return true
	} else if strings.Contains(err.Error(), "connection refused") {
		return true
	}
	return false
}

func (rc *rabbitChannel) connectAndMonitor() error {
	conn, ch, err := rc.connect()
	if err != nil {
		return err
	}

	// kill lingering conn / channel (if any) as we've made new ones
	rc.close()

	rc.conn = conn
	rc.channel = ch

	go func() {
		// listen for anything disconnecting and trigger a reconnect
		// unless we've been told to close
		for {
			select {
			case err := <-rc.conn.NotifyClose(make(chan *amqp.Error)):
				if err == nil {
					continue
				}
				if !rc.closed {
					rc.log.Warn().Err(err).Msg("detected rabbitmq connection/channel error")
				}
				if !isConnectError(*err) {
					continue
				}
				rc.close()
				rc.connectAndMonitor()
				if rc.onReconnect != nil {
					rc.onReconnect <- rc.channel
				}
				return
			case err := <-rc.channel.NotifyClose(make(chan *amqp.Error)):
				if err == nil {
					continue
				}
				if !rc.closed {
					rc.log.Warn().Err(err).Msg("detected rabbitmq connection/channel error")
				}
				if !isConnectError(*err) {
					continue
				}
				rc.close()
				rc.connectAndMonitor()
				if rc.onReconnect != nil {
					rc.onReconnect <- rc.channel
				}
				return
			}
		}
	}()

	return nil
}

// connect to rabbitmq, retrying forever until we get a connection
// and setup our desired channel
func (rc *rabbitChannel) connect() (*amqp.Connection, *amqp.Channel, error) {
	rc.log.Info().Str("username", rc.cfg.Username).Str("host", rc.cfg.Host).Int("port", rc.cfg.Port).Msg("connecting to rabbitmq")
	i := 0
	for {
		i += 1
		if i > 1 {
			time.Sleep(time.Duration(2) * time.Second)
			rc.log.Debug().Int("attempt", i).Msg("retrying connection to rabbitmq")
		}

		url := fmt.Sprintf("amqp://%s:%s@%s:%d/", rc.cfg.Username, rc.cfg.Password, rc.cfg.Host, rc.cfg.Port)
		conn, err := amqp.Dial(url)
		if err != nil {
			rc.log.Warn().Err(err).Msg("failed to connect to rabbitmq")
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			rc.log.Warn().Err(err).Msg("failed to open channel")
			conn.Close()
			continue
		}

		err = ch.Qos(rc.prefetch, 0, false)
		if err != nil {
			rc.log.Warn().Err(err).Msg("failed to set qos")
			ch.Close()
			conn.Close()
			continue
		}

		rc.log.Info().Msg("connected to rabbitmq & established channel")
		return conn, ch, nil
	}
}
