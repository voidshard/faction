package config

// QueueDriver denotes what kind of queue we'll use
type QueueDriver string

const (
	// QueueInMemory is a queue for processing tasks in memory for local only mode.
	// This is is only intended for testing / debugging / development.
	QueueInMemory QueueDriver = "inmemory"

	// QueueIgor uses Igor to process tasks.
	QueueIgor QueueDriver = "igor"
)

const (
	defaultIgorQueueDB  = "postgres://postgres:test@localhost:5432/igor?sslmode=disable"
	defaultIgorQueueURL = "localhost:6379"
)

type Queue struct {
	Driver QueueDriver

	// Location of queue
	// [inmemory]: ignored
	// [igor]: URL of igor redis
	Queue string

	// Location of queue database
	// [inmemory]: ignored
	// [igor]: igor postgres connection string
	Database string
}

func DefaultQueue() *Queue {
	if isLocalMode() {
		return &Queue{Driver: QueueInMemory}
	}
	return &Queue{
		Driver:   QueueIgor,
		Queue:    defaultIgorQueueURL,
		Database: defaultIgorQueueDB,
	}
}
