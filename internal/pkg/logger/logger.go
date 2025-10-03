package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
)

type Logger struct {
	logger    *slog.Logger
	prefix    string
	requestID string
}

var (
	instance *Logger
	once     sync.Once
)

func NewLogger(prefix ...string) *Logger {
	once.Do(func() {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		instance = &Logger{logger: logger}
		if len(prefix) > 0 {
			instance.prefix = prefix[0]
		}
	})

	if len(prefix) > 0 && instance.prefix != prefix[0] {
		instance.prefix = prefix[0]
	}
	return instance
}

func (q *Logger) decorateLog() []any {
	var source string
	if pc, file, line, ok := runtime.Caller(2); ok {
		funcName := runtime.FuncForPC(pc).Name()
		funcName = path.Base(funcName[strings.LastIndex(funcName, ".")+1:])
		source = fmt.Sprintf("%s:%d:%s()", path.Base(file), line, funcName)
	}

	attrs := []any{slog.String("source", source)}

	if q.prefix != "" {
		attrs = append(attrs, slog.String("prefix", q.prefix))
	}

	if q.requestID != "" {
		attrs = append(attrs, slog.String("requestID", q.requestID))
	}

	return attrs
}

func (q *Logger) Output() io.Writer {
	return os.Stdout
}

func (q *Logger) Print(message string) {
	q.logger.Info(message, q.decorateLog()...)
}

func (q *Logger) Printf(format string, args ...interface{}) {
	q.logger.Info(fmt.Sprintf(format, args...), q.decorateLog()...)
}

func (q *Logger) Debug(message string) {
	q.logger.Debug(message, q.decorateLog()...)
}

func (q *Logger) Debugf(format string, args ...interface{}) {
	q.logger.Debug(fmt.Sprintf(format, args...), q.decorateLog()...)
}

func (q *Logger) Info(message string) {
	q.logger.Info(message, q.decorateLog()...)
}

func (q *Logger) Infof(format string, args ...interface{}) {
	q.logger.Info(fmt.Sprintf(format, args...), q.decorateLog()...)
}

func (q *Logger) Warn(message string) {
	q.logger.Warn(message, q.decorateLog()...)
}

func (q *Logger) Warnf(format string, args ...interface{}) {
	q.logger.Warn(fmt.Sprintf(format, args...), q.decorateLog()...)
}

func (q *Logger) Error(message string) {
	q.logger.Error(message, q.decorateLog()...)
}

func (q *Logger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	q.logger.Error(msg, q.decorateLog()...)
}
