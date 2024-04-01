package queue

import (
	"fmt"
	"time"

	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/igor/pkg/api"
	"github.com/voidshard/igor/pkg/database"
	que "github.com/voidshard/igor/pkg/queue"
	"github.com/voidshard/igor/pkg/structs"
)

const (
	igorTaskBatchSize = 500
)

type IgorQueue struct {
	svc *api.Service
}

func NewIgorQueue(cfg *config.Queue) (*IgorQueue, error) {
	svc, err := api.New(
		&database.Options{URL: cfg.Database},
		&que.Options{URL: cfg.Queue},
		&api.Options{
			EventRoutines:       2,
			TidyRoutines:        2,
			TidyLayerFrequency:  2 * time.Minute,
			TidyTaskFrequency:   2 * time.Minute,
			TidyUpdateThreshold: 2 * time.Minute,
			MaxTaskRuntime:      10 * time.Minute,
		},
	)
	return &IgorQueue{svc: svc}, err
}

type IgorWorkflow struct {
	svc        *api.Service
	job        *structs.CreateJobResponse
	layerIDs   map[string]string
	layerETags map[string]string
	tasks      map[string][]*structs.CreateTaskRequest
}

func (q *IgorQueue) NewWorkflow(name string, layers ...[]string) (Workflow, error) {
	cjr := &structs.CreateJobRequest{
		JobSpec: structs.JobSpec{Name: name},
		Layers:  []structs.JobLayerRequest{},
	}

	seen := map[string]bool{}
	now := time.Now().Unix()

	for i, layerSet := range layers {
		for _, layer := range layerSet {
			if _, ok := seen[layer]; ok {
				// strictly speaking, Igor allows this, but it makes it easy for us
				// to differentiate between layers without the user needing to do translation
				return nil, fmt.Errorf("duplicate layer: %s", layer)
			}
			seen[layer] = true
			cjr.Layers = append(cjr.Layers, structs.JobLayerRequest{
				LayerSpec: structs.LayerSpec{
					Name:     layer,
					PausedAt: now, // nb. launching paused
					Priority: int64(i),
				},
			})
		}
	}

	resp, err := q.svc.CreateJob(cjr)
	if err != nil {
		return nil, err
	}

	wf := &IgorWorkflow{
		svc:        q.svc,
		job:        resp,
		layerIDs:   map[string]string{},
		layerETags: map[string]string{},
		tasks:      map[string][]*structs.CreateTaskRequest{},
	}
	for _, l := range resp.Layers {
		wf.layerIDs[l.Name] = l.ID
		wf.layerETags[l.Name] = l.ETag
	}

	return wf, nil
}

func (q *IgorQueue) Register(task string, handler Handler) error {
	return q.svc.Register(task, func(work []*que.Meta) error {
		jobs := make([]*Job, len(work))
		for i, w := range work {
			jobs[i] = &Job{Task: task, Args: w.Task.Args}
		}
		return handler(jobs...)
	})
}

func (q *IgorQueue) Start() error {
	go func() {
		err := q.svc.Run()
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

func (q *IgorQueue) Stop() error {
	return q.svc.Close()
}

func (w *IgorWorkflow) Task(layer, task string, args []byte) error {
	id, ok := w.layerIDs[layer]
	if !ok {
		return fmt.Errorf("layer not found: %s", layer)
	}

	tasks, ok := w.tasks[layer]
	if !ok {
		tasks = []*structs.CreateTaskRequest{}
	}

	tasks = append(tasks, &structs.CreateTaskRequest{
		LayerID: id,
		TaskSpec: structs.TaskSpec{
			Name: task,
			Args: args,
		},
	})
	if len(tasks) >= igorTaskBatchSize {
		_, err := w.svc.CreateTasks(tasks)
		if err != nil {
			return err
		}
		w.tasks[layer] = []*structs.CreateTaskRequest{}
	} else {
		w.tasks[layer] = tasks
	}

	return nil
}

func (w *IgorWorkflow) Unpause(layers ...string) error {
	// if no layers are given, unpause all layers
	if layers == nil || len(layers) == 0 {
		layers = []string{}
		for l := range w.layerIDs {
			layers = append(layers, l)
		}
	}

	// prepare all tasks
	refs := []*structs.ObjectRef{}
	for _, l := range layers {
		// flush any remaining tasks
		tasks, ok := w.tasks[l]
		if ok && len(tasks) > 0 {
			_, err := w.svc.CreateTasks(tasks)
			if err != nil {
				return err
			}
			w.tasks[l] = []*structs.CreateTaskRequest{}
		}

		// work out the object reference
		id, ok := w.layerIDs[l]
		if !ok {
			continue
		}
		etag, ok := w.layerETags[l]
		if !ok {
			continue
		}

		refs = append(refs, &structs.ObjectRef{
			ID:   id,
			ETag: etag,
			Kind: structs.KindLayer,
		})
	}

	count, err := w.svc.Unpause(refs)
	if err != nil {
		return err
	}
	if count != int64(len(refs)) {
		// TODO: we could fetch new ETags and Retry but the current impl launches all layers paused
		// (so they cant change) and we only unpause, after which we don't touch them.
		// So in theory, this should never happen until / unless we add more features.
		// Note that adding tasks to a layer will not change the ETag.
		return fmt.Errorf("failed to unpause all layers")
	}

	return nil
}
