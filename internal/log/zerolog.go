/*
We mimic zerologs no-allocation approach to logging by using a sync.Pool here.
https://github.com/rs/zerolog/blob/master/event.go#L13

We're not adding any functionality at all (actually, we hiding a lot), the aim is only to
remove library specific stuff from our logging interface.

Technically we add an allocation here, so it's now onelog not zerolog .. but .. eh.
*/
package log

import (
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var linePool = &sync.Pool{
	New: func() interface{} { return &logLine{} },
}

func putLine(l *logLine) {
	linePool.Put(l)
}

func getLine(e *zerolog.Event) *logLine {
	l := linePool.Get().(*logLine)
	l.e = e
	l.msg = e.Msg
	return l
}

type logLine struct {
	e   *zerolog.Event
	msg func(string)
}

func (l *logLine) Str(key string, val string) LogLine {
	l.e.Str(key, val)
	return l
}

func (l *logLine) Err(err error) LogLine {
	l.e.Err(err)
	return l
}

func (l *logLine) Int(key string, val int) LogLine {
	l.e.Int(key, val)
	return l
}

func (l *logLine) Float64(key string, val float64) LogLine {
	l.e.Float64(key, val)
	return l
}

func (l *logLine) Msg() func(msg string) {
	// Ensures we call "Msg" where the log is written (for tracing)
	defer putLine(l)
	return l.msg
}

type zLogger struct{}

func (z *zLogger) Info() LogLine {
	return getLine(log.Info())
}

func (z *zLogger) Warn() LogLine {
	return getLine(log.Warn())
}

func (z *zLogger) Error() LogLine {
	return getLine(log.Error())
}

func (z *zLogger) Fatal() LogLine {
	return getLine(log.Fatal())
}

func (z *zLogger) Debug() LogLine {
	return getLine(log.Debug())
}

func newZeroLog() *zLogger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger() // enable logs to include line numbers & file names

	debug := strings.ToLower(strings.TrimSpace(os.Getenv(EnvDebug)))
	if debug == "true" || debug == "1" || debug == "on" || debug == "yes" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	return &zLogger{}
}
