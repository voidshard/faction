package api

/*
import (
	"fmt"

	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/pkg/structs"
)

type OnChangeSubscription struct {
	qu  queue.Queue
	sub queue.Subscription
	log log.Logger
	out chan *structs.Change
}

func newOnChangeSubscription(qu queue.Queue, pattern *structs.Change) (*OnChangeSubscription, error) {
	subject, err := qu.ToSubject(pattern)
	if err != nil {
		return nil, err
	}

	sub, err := qu.Subscribe(subject)
	if err != nil {
		return nil, err
	}

	me := &OnChangeSubscription{
		qu:  qu,
		sub: sub,
		log: log.Sublogger(fmt.Sprintf("api-onchange-subscription:%s", subject)),
		out: make(chan *structs.Change),
	}

	go me.work()
	return me, nil
}

func (o *OnChangeSubscription) work() {
	for message := range o.sub.Channel() {
		change, err := o.qu.FromSubject(message.Subject)
		if err != nil {
			o.log.Error().Err(err).Msg("failed to parse message")
			continue
		}
		o.out <- change
	}
}

func (o *OnChangeSubscription) Channel() <-chan *structs.Change {
	return o.out
}

func (o *OnChangeSubscription) Close() error {
	return o.sub.Close()
}
*/
