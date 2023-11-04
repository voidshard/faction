package config

// QueueDriver denotes what kind of queue we'll use
type QueueDriver string

const (
	// QueueInMemory is a queue for processing tasks in memory for local only mode.
	// This is is only intended for testing / debugging / development.
	QueueInMemory QueueDriver = "inmemory"

	// QueueAsynq uses asynq to process tasks.
	QueueAsynq QueueDriver = "asynq"
)

const (
	defaultAsynqQueueURL = "localhost:6379" // redis
)

type Queue struct {
	Driver   QueueDriver
	Location string
}

func DefaultQueue() *Queue {
	if isLocalMode() {
		return &Queue{Driver: QueueInMemory}
	}
	return &Queue{Driver: QueueAsynq, Location: defaultAsynqQueueURL}
}
