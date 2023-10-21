package log

const (
	// Env var to control debug
	EnvDebug = "ENABLE_DEBUG"
)

var activeLogger logger

// logger is the interface that wraps the basic logging methods.
// We do this so it's easy to switch loggers if we want.
type logger interface {
	Info() LogLine
	Warn() LogLine
	Error() LogLine
	Fatal() LogLine
	Debug() LogLine
}

// LogLine is the interface that wraps the basic logging methods for
// setting key/value pairs and logging a message.
type LogLine interface {
	// Log various key/value pairs
	Str(key string, val string) LogLine
	Int(key string, val int) LogLine
	Float64(key string, val float64) LogLine
	Err(err error) LogLine

	// Msg logs a message with the key/value pairs previously set. It must be called.
	Msg() func(msg string)
}

func Info() LogLine {
	return activeLogger.Info()
}

func Warn() LogLine {
	return activeLogger.Warn()
}

func Error() LogLine {
	return activeLogger.Error()
}

func Fatal() LogLine {
	return activeLogger.Fatal()
}

func Debug() LogLine {
	return activeLogger.Debug()
}

func init() {
	activeLogger = newZeroLog()
}
