package log

import (
	"fmt"
	"os"
	"reflect"
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

func Sublogger(name string, attrs ...map[string]interface{}) zerolog.Logger {
	logger := logger.With()
	// first set any attrs
	for _, attr := range attrs {
		for k, v := range attr {
			var val reflect.Value
			if reflect.TypeOf(v).Kind() == reflect.Ptr {
				// if the value is a pointer, dereference it (otherwise we get the pointer address)
				val = reflect.ValueOf(v).Elem()
			} else {
				val = reflect.ValueOf(v)
			}

			switch val.Kind() {
			case reflect.String:
				logger = logger.Str(k, val.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				logger = logger.Int64(k, val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				logger = logger.Uint64(k, val.Uint())
			case reflect.Float32, reflect.Float64:
				logger = logger.Float64(k, val.Float())
			case reflect.Bool:
				logger = logger.Bool(k, val.Bool())
			default:
				Warn().Str("key", k).Str("value", fmt.Sprintf("%v", v)).Msg("Unknown type for logger attribute")
			}
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
