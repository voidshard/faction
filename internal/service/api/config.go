package api

import "time"

const (
	defaultAPIRoutines = 5
	defaultMaxAge      = 60 * time.Minute // implies something is horribly wrong
	defaultLimit       = 100
	defaultMaxLimit    = 1000
)

type Config struct {
	// WorkersAPI is the number of worker routines to use for processing API requests.
	WorkersAPI int

	// WorkersDefer is the number of worker routines to use for processing deferred API requests.
	WorkersDefer int

	// MaxMessageAge is the maximum age of a message before it is considered stale.
	MaxMessageAge time.Duration

	// FlushSearch waits for writes to searchbase before returning API writes (slow).
	FlushSearch bool
}

func (c *Config) setDefaults() {
	if c.WorkersAPI < 1 {
		c.WorkersAPI = defaultAPIRoutines
	}
	if c.MaxMessageAge == 0 {
		c.MaxMessageAge = defaultMaxAge
	}
}
