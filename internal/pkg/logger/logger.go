package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"runtime"
	"strings"
)

type Logger interface {
	WithRequestID(requestID string) Logger
	Print(message string)
	Printf(format string, args ...interface{})
	Debug(message string)
	Debugf(format string, args ...interface{})
	Info(message string)
	Infof(format string, args ...interface{})
	Warn(message string)
	Warnf(format string, args ...interface{})
	Error(message string)
	Errorf(format string, args ...interface{})
	Output() io.Writer
}

type log struct {
	logger    *slog.Logger
	prefix    string
	requestID string
}

func NewLogger(prefix ...string) Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	instance := &log{logger: logger}
	if len(prefix) > 0 {
		instance.prefix = prefix[0]
	}

	return instance
}

func (q *log) WithRequestID(requestID string) Logger {
	return &log{
		logger:    q.logger,
		prefix:    q.prefix,
		requestID: requestID,
	}
}

func (q *log) decorateLog() []any {
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

func (q *log) Output() io.Writer {
	return os.Stdout
}

func (q *log) Print(message string) {
	q.logger.Info(message, q.decorateLog()...)
}

func (q *log) Printf(format string, args ...interface{}) {
	q.logger.Info(fmt.Sprintf(format, args...), q.decorateLog()...)
}

func (q *log) Debug(message string) {
	q.logger.Debug(message, q.decorateLog()...)
}

func (q *log) Debugf(format string, args ...interface{}) {
	q.logger.Debug(fmt.Sprintf(format, args...), q.decorateLog()...)
}

func (q *log) Info(message string) {
	q.logger.Info(message, q.decorateLog()...)
}

func (q *log) Infof(format string, args ...interface{}) {
	q.logger.Info(fmt.Sprintf(format, args...), q.decorateLog()...)
}

func (q *log) Warn(message string) {
	q.logger.Warn(message, q.decorateLog()...)
}

func (q *log) Warnf(format string, args ...interface{}) {
	q.logger.Warn(fmt.Sprintf(format, args...), q.decorateLog()...)
}

func (q *log) Error(message string) {
	q.logger.Error(message, q.decorateLog()...)
}

func (q *log) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	q.logger.Error(msg, q.decorateLog()...)
}
