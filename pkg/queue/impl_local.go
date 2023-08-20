package queue

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/groupcache/lru"

	"github.com/voidshard/faction/internal/dbutils"
)

// Implementation of the queue interface with local routines.
// This is super basic for local-only mode, not intended for use in large scales.

const (
	maxOldTasks = 1024
)

type localQueue struct {
	handlerslock *sync.Mutex
	handlers     map[string]*aggregatingQueue

	jobsLock *sync.Mutex
	jobs     map[string]*JobMeta
	jobsOld  *lru.Cache

	queue chan []*JobMeta
}

func NewLocalQueue(routines int) Queue {
	if routines < 1 {
		routines = 1
	}
	me := &localQueue{
		handlerslock: &sync.Mutex{},
		handlers:     make(map[string]*aggregatingQueue),
		jobsLock:     &sync.Mutex{},
		jobs:         make(map[string]*JobMeta),
		jobsOld:      lru.New(maxOldTasks),
		queue:        make(chan []*JobMeta),
	}

	for i := 0; i < routines; i++ {
		go me.run()
	}

	return me
}

func (q *localQueue) Close() error {
	q.handlerslock.Lock()
	defer q.handlerslock.Unlock()

	for _, h := range q.handlers {
		h.Close()
	}

	for _, h := range q.handlers {
		h.Flush()
	}

	close(q.queue)
	return nil
}

func (q *localQueue) Await(ids ...string) (State, error) {
	if len(ids) == 0 {
		// If no ids are given, await all tasks. Strictly a local-mode thing.
		// This doesn't prevent adding new tasks.
		for {
			q.jobsLock.Lock()
			num := len(q.jobs)
			q.jobsLock.Unlock()

			if num <= 0 {
				return StateComplete, nil
			}
			staggeredSleep(num)
		}
	}
	return q.awaitTasks(ids...)
}

func (q *localQueue) awaitTasks(in ...string) (State, error) {
	if len(in) == 0 {
		return StateComplete, nil
	}

	seenFailure := false
	waitingOn := in
	for {
		notDone := []string{}
		for _, id := range waitingOn {
			j, err := q.Status(id)
			if err != nil {
				return StateFailed, err
			}

			if j.State == StateFailed {
				seenFailure = true
			} else if j.State == StateComplete {
				// no op
			} else {
				notDone = append(notDone, id)
			}
		}
		if len(notDone) == 0 {
			if seenFailure {
				return StateFailed, nil
			}
			return StateComplete, nil
		}
		waitingOn = notDone
		staggeredSleep(len(waitingOn))
	}
}

func staggeredSleep(count int) {
	if count < 10 {
		time.Sleep(1 * time.Second)
	} else if count < 100 {
		time.Sleep(5 * time.Second)
	} else {
		time.Sleep(10 * time.Second)
	}
}

func (q *localQueue) run() {
	for jobs := range q.queue {
		if len(jobs) == 0 {
			continue // ??
		}

		fn := q.handler(jobs[0].Job.Task)
		if fn == nil {
			// ??
			continue
		}

		for _, j := range jobs {
			q.setStatus(j.Job.ID, StateRunning)
		}

		err := fn.Handler(jobs...)
		if err != nil {
			// implies something beyond the task is broken, in a non-local case we should retry
			for _, j := range jobs {
				q.setStatus(j.Job.ID, StateFailed)
			}
			continue
		}

		for _, j := range jobs {
			if j.State == StateComplete {
				q.setStatus(j.Job.ID, StateComplete)
			} else {
				q.setStatus(j.Job.ID, StateFailed)
			}
		}
	}
}

func (q *localQueue) setStatus(id string, state State) {
	q.jobsLock.Lock()
	defer q.jobsLock.Unlock()

	j, ok := q.jobs[id]
	if !ok {
		return
	}

	j.State = state

	if state == StateComplete || state == StateFailed {
		q.jobsOld.Add(id, j)
		delete(q.jobs, id)
	}

	//log.Printf("%s %s -> %s (%s)\n", j.Job.ID, j.Job.Task, j.State, j.Msg)
}

func (q *localQueue) handler(task string) *aggregatingQueue {
	q.handlerslock.Lock()
	defer q.handlerslock.Unlock()
	h, _ := q.handlers[task]
	return h
}

func (q *localQueue) Register(task string, handler Handler) error {
	q.handlerslock.Lock()
	defer q.handlerslock.Unlock()

	q.handlers[task] = newAggregatingQueue(handler, q.queue)
	return nil
}

func (q *localQueue) Enqueue(task string, args []byte) (string, error) {
	id := dbutils.RandomID()
	j := &JobMeta{Job: &Job{ID: id, Task: task, Args: args}}

	q.jobsLock.Lock()
	q.jobs[j.Job.ID] = j
	q.jobsLock.Unlock()

	q.handlerslock.Lock()
	hnd, ok := q.handlers[task]
	defer q.handlerslock.Unlock()
	if !ok {
		return "", fmt.Errorf("%w %s", ErrNoHandler, task)
	}

	q.setStatus(j.Job.ID, StateQueued)
	hnd.Enqueue(j)

	return j.Job.ID, nil
}

func (q *localQueue) Status(id string) (*JobMeta, error) {
	q.jobsLock.Lock()
	defer q.jobsLock.Unlock()

	if j, ok := q.jobs[id]; ok {
		return j, nil
	}
	if j, ok := q.jobsOld.Get(id); ok {
		return j.(*JobMeta), nil
	}

	return nil, fmt.Errorf("unknown task id %s", id)
}
