package config

const (
	defaultQueueURL = "localhost:6379" // redis
)

type Queue struct {
	Location string
}

func DefaultQueue() *Queue {
	return &Queue{Location: defaultQueueURL}
}
