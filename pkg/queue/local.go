package queue

import (
	"fmt"
)

// InMemoryQueue is a queue for processing tasks in memory for local only mode.
//
// This is is only intended for testing / debugging / development.
type InMemoryQueue struct {
	handlers map[string]Handler
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{handlers: map[string]Handler{}}
}

func (q *InMemoryQueue) Register(task string, handler Handler) error {
	q.handlers[task] = handler
	return nil
}

func (q *InMemoryQueue) Start() error {
	return nil
}

func (q *InMemoryQueue) Stop() error {
	return nil
}

func (q *InMemoryQueue) Enqueue(task string, args []byte) (string, error) {
	// We don't queue anything, we simply do the task immediately
	hnd, _ := q.handlers[task]
	if hnd == nil {
		hnd = noOpHandler
	}
	return "", hnd(&Job{Task: task, Args: args})
}

func (q *InMemoryQueue) Status(id string) (*JobMeta, error) {
	return nil, fmt.Errorf("not implemented: in memory queue does not support status")
}
