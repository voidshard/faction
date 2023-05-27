package work

type Handler func(args []byte) ([]byte, error)

type Queue interface {
	Register(task string, handler Handler) error
	Subscribe() <-chan *Message

	EnqueueSequential(id, name string, jobs ...*Job) error
	EnqueueParallel(id, name string, jobs ...*Job) error
}

type Message struct {
	GroupID string
	JobID   string

	Data  []byte
	Error bool
}

type Job struct {
	ID   string
	Name string

	Task string

	Queue      string
	MaxRetries uint
	Args       []byte
}
