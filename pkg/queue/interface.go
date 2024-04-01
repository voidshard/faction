package queue

import (
	"fmt"

	"github.com/voidshard/faction/pkg/config"
)

var (
	ErrNoHandler = fmt.Errorf("no handler for task")
)

type Handler func(in ...*Job) error

type Queue interface {
	// Register a task to be handled by the queue
	Register(task string, handler Handler) error

	// Start starts tasks processing
	Start() error

	// Stop stops tasks processing
	Stop() error

	// NewWorkflow creates a new workflow with the given name and layers.
	// Layers are a list of tasks to be run in parallel, each set of layers is run in sequence.
	//
	// ie.
	// > NewWorkflow("test", []string{"a", "b"}, []string{"c", "d"})
	// would run a and b in parallel, then c and d in parallel.
	//
	// A workflow is created paused, and Unpause() is called to actually begin processing
	// queued tasks.
	NewWorkflow(name string, layers ...[]string) (Workflow, error)
}

type Workflow interface {
	// Task queues task to be run in the given layer of the workflow
	Task(layer, task string, args []byte) error

	// Unpause allows the given layer(s) to run. If no layers are given, all layers are permitted to unpause.
	// (Layer order is still respected).
	Unpause(layers ...string) error
}

type Job struct {
	Task string `json:"task,omitempty"`
	Args []byte `json:"args,omitempty"`
}

func New(cfg *config.Queue) (Queue, error) {
	switch cfg.Driver {
	case config.QueueInMemory:
		return NewInMemoryQueue(), nil
	case config.QueueIgor:
		return NewIgorQueue(cfg)
	default:
		return nil, fmt.Errorf("unknown queue driver: %s", cfg.Driver)
	}
}
