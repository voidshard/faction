package api

import (
	"time"

	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"
)

var (
	kindSetWorld   = (&structs.SetWorldRequest{}).Kind()
	kindSetFaction = (&structs.SetFactionsRequest{}).Kind()
	kindSetActor   = (&structs.SetActorsRequest{}).Kind()
	kindWorld      = (&structs.World{}).Kind()
	kindFaction    = (&structs.Faction{}).Kind()
	kindActor      = (&structs.Actor{}).Kind()
)

type worker struct {
	kill chan bool
	sub  queue.Subscription

	svr *Server
	log log.Logger
}

func newWorker(name string, svr *Server, sub queue.Subscription) *worker {
	return &worker{
		kill: make(chan bool),
		sub:  sub,
		svr:  svr,
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
		if msg.Timestamp().Add(w.svr.cfg.MaxMessageAge).Before(time.Now()) {
			w.log.Debug().Str("MessageId", msg.Id()).Err(msg.Ack()).Msg("Message too old, acking and dropping")
			return
		}

		w.log.Debug().Str("MessageId", msg.Id()).Msg("Processing message")
		defer w.log.Debug().Str("MessageId", msg.Id()).Msg("Finished processing message")

		pan := log.NewSpan(msg.Context(), "api.asyncAPIRequest", map[string]interface{}{"mid": msg.Id()})
		defer pan.End()

		err := w.svr.asyncAPIRequest(msg.Context(), msg)
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
