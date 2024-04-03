package queue

import (
	"fmt"
	"sync"
)

const (
	inmemRoutines = 10
)

// InMemoryQueue is a queue for processing tasks in memory for local only mode.
//
// This is is only intended for testing / debugging / development.
//
// This is not efficient at all, and is not intended to be.
type InMemoryQueue struct {
	handlers map[string]Handler
	hlock    sync.Mutex
}

type InMemoryWorkflow struct {
	que    *InMemoryQueue
	layers [][]string                     // input layers
	tasks  map[string]map[string][][]byte // layer -> task -> []task args
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{handlers: map[string]Handler{}, hlock: sync.Mutex{}}
}

func (q *InMemoryQueue) Register(task string, handler Handler) error {
	q.hlock.Lock()
	defer q.hlock.Unlock()
	q.handlers[task] = handler
	return nil
}

func (q *InMemoryQueue) Start() error {
	return nil
}

func (q *InMemoryQueue) Stop() error {
	return nil
}

func (q *InMemoryQueue) NewWorkflow(name string, layers ...[]string) (Workflow, error) {
	seen := map[string]bool{}
	tasks := map[string]map[string][][]byte{}
	for _, setOfLayers := range layers {
		for _, l := range setOfLayers {
			_, ok := seen[l]
			if ok {
				return nil, fmt.Errorf("duplicate layer: %s", l)
			}
			seen[l] = true
			tasks[l] = map[string][][]byte{}
		}
	}
	return &InMemoryWorkflow{que: q, layers: layers, tasks: tasks}, nil
}

func (w *InMemoryWorkflow) start() error {
	for _, layerSet := range w.layers { // all tasks of each layer at the same priority can run concurrently
		errs := make(chan error)
		errWg := sync.WaitGroup{}
		errWg.Add(1)

		var ferr error
		go func() {
			defer errWg.Done()
			for err := range errs {
				if ferr == nil {
					ferr = err
				} else {
					ferr = fmt.Errorf("%v %w", ferr, err)
				}
			}
		}()

		wg := sync.WaitGroup{}
		work := make(chan *Job)

		for i := 0; i < inmemRoutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for task := range work {
					w.que.hlock.Lock()
					handler, ok := w.que.handlers[task.Task]
					w.que.hlock.Unlock()

					if !ok {
						continue
					}

					err := handler(task)
					if err != nil {
						errs <- err
					}
				}
			}()
		}

		for _, layerName := range layerSet {
			for task, allTaskArgs := range w.tasks[layerName] {
				for _, args := range allTaskArgs {
					work <- &Job{Task: task, Args: args}
				}
			}
		}

		close(work)
		wg.Wait()
		close(errs)
		errWg.Wait()

		if ferr != nil {
			return ferr
		}
	}

	return nil
}

func (w *InMemoryWorkflow) Task(layer, task string, args []byte) error {
	if _, ok := w.tasks[layer]; !ok {
		return fmt.Errorf("unknown layer: %s", layer)
	}
	regTasks, ok := w.tasks[layer][task]
	if !ok {
		w.tasks[layer][task] = [][]byte{args}
	} else {
		w.tasks[layer][task] = append(regTasks, args)
	}
	return nil
}

func (w *InMemoryWorkflow) Unpause(layers ...string) error {
	if layers != nil {
		// Nb. we unpause everything; otherwise we'd have to implement the logic to figure out
		// what can run and when .. and this is just a local test queue
		return fmt.Errorf("unpausing specific layers not supported in local mode")
	}
	return w.start()
}
