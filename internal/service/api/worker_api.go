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

// API Worker handles write operations to the API, reading from a queue and acknowledging
// a request as done only when all writes have been completed or cannot occur.
//
// Ie. an update to the database with a given etag followed by submitting a message to a queue
// can be retried if the database resource either hasn't been updated (etag is the same as
// what our request is updating) or the database update would do nothing (etag is what our
// request would set). In this way we can retry an operation across the database,
// queue and / or search, up until something else has done an update (in which case our etag will
// mismatch and we know the operation cannot be done).
type apiWorker struct {
	kill chan bool
	sub  queue.Subscription

	svr *Server
	log log.Logger
}

func newApiWorker(name string, svr *Server, sub queue.Subscription) *apiWorker {
	return &apiWorker{
		kill: make(chan bool),
		sub:  sub,
		svr:  svr,
		log:  log.Sublogger(name),
	}
}

func (w *apiWorker) Kill() {
	w.log.Debug().Msg("Killing worker")
	defer w.log.Debug().Msg("Worker killed")

	w.kill <- true

	close(w.kill)
}

func (w *apiWorker) Run() {
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
