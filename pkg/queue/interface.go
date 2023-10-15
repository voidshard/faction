package queue

import (
	"fmt"
)

var (
	ErrNoHandler = fmt.Errorf("no handler for task")
)

type State string

const (
	StateQueued   State = "queued"
	StateRunning  State = "running"
	StateComplete State = "complete"
	StateFailed   State = "failed"
)

type Handler func(in ...*Job) error

type Queue interface {
	// -- server functions
	// Register a task to be handled by the queue
	Register(task string, handler Handler) error

	// Start starts tasks processing
	Start() error

	// Stop stops tasks processing
	Stop() error

	// -- client functions
	// Enqueue adds a task to the queue
	Enqueue(task string, args []byte) (string, error)

	// -- misc functions
	// Status returns the status of a task
	Status(id string) (*JobMeta, error)
}

type Job struct {
	Task string `json:"task,omitempty"`
	Args []byte `json:"args,omitempty"`
}

type JobMeta struct {
	ID     string `json:"id"`
	State  State  `json:"state"`
	Result []byte `json:"result,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Job    *Job
}
