package logger

import (
	stdlog "log"
	"os"

	kitlog "github.com/go-kit/kit/log"
)

// Log for logging in all components
type Log struct {
	logger kitlog.Logger
}

// New JSON logger based off kitlog writing to stdout
func New(out *os.File) *Log {
	l := &Log{}
	l.logger = kitlog.NewJSONLogger(kitlog.NewSyncWriter(out))
	l.logger = kitlog.With(l.logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.Caller(4))
	stdlog.SetOutput(kitlog.NewStdlibAdapter(l.logger))
	return l
}

// Info level logs
func (l *Log) Info(fields ...interface{}) {
	kvs := append([]interface{}{"info"}, fields...)
	l.logger.Log(kvs...)
}

// Error level logs
func (l *Log) Error(fields ...interface{}) {
	kvs := append([]interface{}{"error"}, fields...)
	l.logger.Log(kvs...)
}
