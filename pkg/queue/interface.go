package queue

type State string

const (
	StateQueued   State = "queued"
	StateRunning  State = "running"
	StateComplete State = "complete"
	StateFailed   State = "failed"
)

type Handler func(in ...*JobMeta) error

type Queue interface {
	Register(task string, handler Handler) error
	Enqueue(task string, args []byte) (string, error)
	Status(id string) (*JobMeta, error)
	Await(ids ...string) (State, error)
	Close() error
}

type Job struct {
	ID   string
	Task string
	Args []byte
}

type JobMeta struct {
	State  State
	Result []byte
	Msg    string
	Job    *Job
}
