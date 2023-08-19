package queue

import (
	"sync"
	"time"
)

const (
	gracePeriod = 5 * time.Second
	maxBuffer   = 1000
)

type aggregatingQueue struct {
	Handler  Handler
	lock     *sync.Mutex
	buffer   []*JobMeta
	incoming chan *JobMeta
	outgoing chan<- []*JobMeta
	quit     chan bool
}

func newAggregatingQueue(hnd Handler, out chan<- []*JobMeta) *aggregatingQueue {
	q := &aggregatingQueue{
		Handler:  hnd,
		lock:     &sync.Mutex{},
		buffer:   []*JobMeta{},
		incoming: make(chan *JobMeta),
		outgoing: out,
		quit:     make(chan bool),
	}
	go q.enqueue()
	go q.run()
	return q
}

func (q *aggregatingQueue) Close() {
	close(q.incoming)
	q.quit <- true
	close(q.quit)
}

func (q *aggregatingQueue) Enqueue(j *JobMeta) {
	q.incoming <- j
}

func (q *aggregatingQueue) enqueue() {
	for j := range q.incoming {
		q.lock.Lock()
		q.buffer = append(q.buffer, j)
		q.lock.Unlock()
	}
}

func (q *aggregatingQueue) run() {
	exit := false
	ticker := time.NewTicker(gracePeriod)
	for {
		select {
		case <-q.quit:
			exit = true
		case <-ticker.C:
			q.lock.Lock()
			for i := 0; i < len(q.buffer); i += maxBuffer {
				j := i + maxBuffer
				if j > len(q.buffer) {
					q.outgoing <- q.buffer[i:]
				} else {
					q.outgoing <- q.buffer[i:j]
				}
			}
			q.buffer = []*JobMeta{}
			q.lock.Unlock()
			if exit {
				return
			}
		}
	}
}
