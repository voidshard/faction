package api

import (
	"context"
	"fmt"
	"time"

	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/queue"
)

type worker struct {
	kill chan bool
	sub  queue.Subscription

	svc *Service
	log log.Logger
}

func newWorker(name string, svc *Service, sub queue.Subscription) *worker {
	return &worker{
		kill: make(chan bool),
		sub:  sub,
		svc:  svc,
		log:  log.Sublogger(name),
	}
}

func (w *worker) Kill() {
	w.log.Debug().Msg("Killing worker")
	defer w.log.Debug().Msg("Worker killed")

	w.kill <- true
	w.sub.Close()

	close(w.kill)
}

func (w *worker) Run() {
	w.log.Debug().Msg("Worker started")
	defer w.log.Debug().Msg("Worker stopped")

	processMessage := func(msg queue.Message) {
		if msg.Timestamp().Add(w.svc.cfg.MaxMessageAge).Before(time.Now()) {
			w.log.Debug().Str("MessageId", msg.Id()).Err(msg.Ack()).Msg("Message too old, acking and dropping")
			return
		}

		w.log.Debug().Str("MessageId", msg.Id()).Msg("Processing message")
		defer w.log.Debug().Str("MessageId", msg.Id()).Msg("Finished processing message")

		pan := log.NewSpan(msg.Context(), "api.asyncAPIRequest", map[string]interface{}{"mid": msg.Id()})
		defer pan.End()

		err := w.asyncAPIRequest(msg.Context(), msg)
		if err != nil {
			w.log.Error().Err(err).Msg("Failed to process message")
			pan.Err(err)
		}
	}

	for {
		select {
		case <-w.kill:
			return
		case msg := <-w.sub.Channel():
			processMessage(msg)
		}
	}
}

func (w *worker) asyncAPIRequest(ctx context.Context, msg queue.Message) error {
	method, data, err := decodeRequest(msg.Data())
	if err != nil {
		return err
	}

	switch method {
	case "SetWorld":
		return w.svc.setWorld(ctx, msg, data)
	case "DeleteWorld":
		return w.svc.deleteWorld(ctx, msg, data)
	case "SetFactions":
		return w.svc.setFactions(ctx, msg, data)
	case "DeleteFaction":
		return w.svc.deleteFaction(ctx, msg, data)
	case "SetActors":
		return w.svc.setActors(ctx, msg, data)
	case "DeleteActor":
		return w.svc.deleteActor(ctx, msg, data)
	default:
		return fmt.Errorf("unknown method")
	}
}
