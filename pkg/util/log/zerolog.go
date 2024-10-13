package log

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

const (
	// Env var to control debug
	EnvDebug = "DEBUG"
)

var (
	// clogger "console logger" we know we can always write to
	clogger zerolog.Logger

	// logger is the main logger, which can be multi level or just the console logger
	logger zerolog.Logger
)

type Logger = zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	SetGlobalLevel()

	console := zerolog.ConsoleWriter{Out: os.Stdout}
	clogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()

	logger = zerolog.New(zerolog.MultiLevelWriter(console)).With().Caller().Logger()
}

func SetGlobalLevel() {
	debug := strings.ToLower(strings.TrimSpace(os.Getenv(EnvDebug)))
	if debug == "true" || debug == "1" || debug == "on" || debug == "yes" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func Sublogger(name string, attrs ...map[string]string) zerolog.Logger {
	logger := logger.With()
	// first set any attrs
	for _, attr := range attrs {
		for k, v := range attr {
			logger = logger.Str(k, v)
		}
	}
	// make sure we set name and return the logger
	return logger.Str("logger", name).Logger()
}

func Info() *zerolog.Event {
	return logger.Info()
}

func Warn() *zerolog.Event {
	return logger.Warn()
}

func Debug() *zerolog.Event {
	return logger.Debug()
}

func Error() *zerolog.Event {
	return logger.Error()
}

func Fatal() *zerolog.Event {
	return logger.Fatal()
}
