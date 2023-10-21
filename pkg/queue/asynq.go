package queue

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hibiken/asynq"

	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/pkg/config"
)

const (
	asyncWorkQueue   = "work"
	asyncPollTime    = 1 * time.Second
	asyncAggMaxSize  = 100
	asyncAggMaxDelay = 2 * time.Second
	asyncAggRune     = "¬"
)

type Asynq struct {
	cfg *config.Queue
	ins *asynq.Inspector

	lock *sync.Mutex
	srv  *asynq.Server
	mux  *asynq.ServeMux

	clt *asynq.Client
}

func NewAsynqQueue(cfg *config.Queue) (*Asynq, error) {
	ins := asynq.NewInspector(asynq.RedisClientOpt{Addr: cfg.Location})
	return &Asynq{
		cfg:  cfg,
		ins:  ins,
		lock: &sync.Mutex{},
	}, nil
}

func (a *Asynq) buildServer() {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.mux != nil {
		// someone locked and set this first
		return
	}
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: a.cfg.Location},
		asynq.Config{
			Queues:          map[string]int{asyncWorkQueue: 1},
			GroupAggregator: asynq.GroupAggregatorFunc(aggregate),
			GroupMaxSize:    asyncAggMaxSize,
			GroupMaxDelay:   asyncAggMaxDelay,
		},
	)
	mux := asynq.NewServeMux()
	a.srv = srv
	a.mux = mux
}

func (a *Asynq) Register(task string, handler Handler) error {
	if a.mux == nil {
		a.buildServer()
	}
	log.Debug().Str("task", task).Msg()("registering task")
	a.mux.HandleFunc(aggregatedTask(task), func(ctx context.Context, t *asynq.Task) error {
		jobs := []*Job{}
		for _, load := range bytes.Split(t.Payload(), []byte(asyncAggRune)) {
			load = bytes.TrimSpace(load)
			if len(load) == 0 {
				continue
			}
			jobs = append(jobs, &Job{Task: task, Args: load})
		}
		err := handler(jobs...)
		log.Info().Err(err).Str("task", task).Int("jobs", len(jobs)).Msg()("handling jobs")
		return err
	})
	return nil
}

func (a *Asynq) Start() error {
	log.Info().Msg()("starting queue processing")
	if a.srv == nil {
		a.buildServer()
	}
	return a.srv.Run(a.mux)
}

func (a *Asynq) Enqueue(task string, args []byte) (string, error) {
	if a.clt == nil {
		a.clt = asynq.NewClient(asynq.RedisClientOpt{Addr: a.cfg.Location})
	}
	qtask := asynq.NewTask(task, args)
	info, err := a.clt.Enqueue(qtask, asynq.Queue(asyncWorkQueue), asynq.Group(aggregatedTask(task)))
	log.Debug().Err(err).Str("task", task).Str("id", info.ID).Msg()("enqueued task")
	return info.ID, err
}

func (a *Asynq) Status(id string) (*JobMeta, error) {
	info, err := a.ins.GetTaskInfo(asyncWorkQueue, id)
	if err != nil {
		return nil, err
	}
	return toMeta(info), nil
}

func (a *Asynq) await() error {
	// wait for everything
	for {
		info, err := a.ins.GetQueueInfo(asyncWorkQueue)
		log.Debug().Err(err).Int("pending", info.Pending).Int("active", info.Active).Int("retry", info.Retry).Msg()("awaiting empty queue")
		if err != nil {
			log.Debug().Err(err).Msg()("awaiting empty queue")
			return err
		}
		log.Debug().Int("pending", info.Pending).Int("active", info.Active).Int("retry", info.Retry).Msg()("awaiting empty queue")
		if info.Pending+info.Active+info.Retry == 0 {
			return nil
		}
		time.Sleep(asyncPollTime)
	}
}

func (a *Asynq) Stop() error {
	log.Debug().Msg()("stopping queue processing")
	if a.srv == nil {
		return nil
	}

	// stop accepting new tasks
	a.srv.Stop()

	// shutdown
	a.srv.Shutdown()
	return nil
}

func toMeta(info *asynq.TaskInfo) *JobMeta {
	return &JobMeta{
		ID:     info.ID,
		State:  toState(info.LastErr, info.State),
		Result: info.Result,
		Msg:    info.LastErr,
		Job: &Job{
			Task: info.Type,
			Args: info.Payload,
		},
	}
}

func toState(lastErr string, s asynq.TaskState) State {
	switch s {
	case asynq.TaskStateActive:
		return StateRunning
	case asynq.TaskStatePending, asynq.TaskStateScheduled, asynq.TaskStateRetry, asynq.TaskStateAggregating:
		return StateQueued
	default:
		if lastErr != "" {
			return StateComplete
		} else {
			return StateFailed
		}
	}
}

func aggregate(group string, tasks []*asynq.Task) *asynq.Task {
	log.Debug().Str("group", group).Int("tasks", len(tasks)).Msg()("aggregating tasks")
	var b strings.Builder
	for _, t := range tasks {
		if t == nil || t.Payload() == nil {
			continue
		}
		b.Write(t.Payload())
		b.WriteString(asyncAggRune)
	}
	log.Debug().Str("payload", b.String()).Msg()("aggregated tasks payload")
	return asynq.NewTask(group, []byte(b.String()))
}

func aggregatedTask(group string) string {
	return fmt.Sprintf("aggregated%s", group)
}
