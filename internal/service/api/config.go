package api

import "time"

const (
	defaultPublishRoutines = 5
	defaultMaxAge          = 60 * time.Minute // implies something is horribly wrong
	defaultLimit           = 100
	defaultMaxLimit        = 1000
)

type Config struct {
	// MaxMessageAge is the maximum age of a message before it is considered stale.
	MaxMessageAge time.Duration

	// FlushSearch waits for writes to searchbase before returning API writes (slow).
	FlushSearch bool

	// Timeout is the maximum time to wait for a read or write operation.
	TimeoutRead  time.Duration
	TimeoutWrite time.Duration

	PublishRoutines int
}

func (c *Config) setDefaults() {
	if c.MaxMessageAge == 0 {
		c.MaxMessageAge = defaultMaxAge
	}
	if c.TimeoutRead == 0 {
		c.TimeoutRead = 60 * time.Second
	}
	if c.TimeoutWrite == 0 {
		c.TimeoutWrite = 60 * time.Second
	}
	if c.PublishRoutines <= 0 {
		c.PublishRoutines = defaultPublishRoutines
	}
}
