package logger

type Level uint8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var LevelMap = map[string]Level{"debug": DEBUG, "info": INFO, "warn": WARN, "error": ERROR, "fatal": FATAL}

//附加信息
type Extra struct {
	Key   string
	Value interface{}
}

type Logger interface {
	SetLevel(level Level)
	Debugf(msg string, args []Extra)
	Infof(msg string, args []Extra)
	Warnf(msg string, args []Extra)
	Errorf(msg string, args []Extra)
	Fatalf(msg string, args []Extra)
}
