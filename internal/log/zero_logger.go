package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"runtime"
)

type ZeroLogger struct {
	writer io.Writer
	level  Level
}

type CallerHook struct{}

func (h CallerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if _, file, line, ok := runtime.Caller(5); ok {
		e.Str("file", file)
		e.Int("line", line)
	}
}

// NewZeroLogger return a configured instance of NewZeroLogger
func NewZeroLogger(writer io.Writer, level Level) *ZeroLogger {
	zeroLogger := ZeroLogger{writer, level}
	zeroLogger.configureLogger()
	return &zeroLogger
}

func (l *ZeroLogger) configureLogger() {
	var zLevel zerolog.Level
	switch l.level {
	case LevelInfo:
		zLevel = zerolog.InfoLevel
	case LevelError:
		zLevel = zerolog.ErrorLevel
	case LevelFatal:
		zLevel = zerolog.FatalLevel
	case LevelOff:
		zLevel = zerolog.Disabled
	default:
		zLevel = zerolog.InfoLevel
	}
	log.Logger = zerolog.New(l.writer).With().Timestamp().Logger().Level(zLevel)
	log.Logger = log.Hook(CallerHook{})
}

// Info only logs information
func (l *ZeroLogger) Info(message string, properties Fields) {
	log.Info().Fields(getFields(properties)).Msg(message)
}

func getFields(properties Fields) map[string]interface{} {
	props := make(map[string]interface{})
	for k, v := range properties {
		props[k] = v
	}
	return props
}

// Error reports all error at error level
func (l *ZeroLogger) Error(err error, properties Fields) {
	log.Error().Fields(getFields(properties)).Err(err).Msg(err.Error())
}

// Fatal write the log to output and stop the process
func (l *ZeroLogger) Fatal(err error, properties Fields) {
	log.Fatal().Fields(getFields(properties)).Err(err).Msg(err.Error())
}

// Debug this is for debugging and we use it to store some information in the log
func (l *ZeroLogger) Debug(message string, properties Fields) {
	log.Debug().Fields(getFields(properties)).Msg(message)
}

func (l *ZeroLogger) SetLevel(level Level) {
	l.level = level
	l.configureLogger()
}
