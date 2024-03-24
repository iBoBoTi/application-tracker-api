package log

type Fields map[string]interface{}

type Logger interface {
	Info(message string, properties Fields)
	Error(err error, properties Fields)
	Fatal(err error, properties Fields)
	Debug(message string, properties Fields)
	SetLevel(level Level)
}

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
	LevelDebug
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelDebug:
		return "DEBUG"
	default:
		return ""
	}
}
