package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

// structuredLogger holds our application's instance of our logger
type structuredLogger struct {
	logger *logrus.Logger
}

// newLogEntry will return a new log entry scoped to the http.Request
func (l *structuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	le := logrus.NewEntry(l.logger)
	le.WithContext(r.Context())
	entry := &structuredLoggerEntry{logger: le}

	logFields := logrus.Fields{}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	logFields["http_method"] = r.Method

	logFields["uri"] = r.RequestURI

	entry.logger = entry.logger.WithFields(logFields)

	return entry
}

// structuredLoggerEntry holds our FieldLogger entry
type structuredLoggerEntry struct {
	logger logrus.FieldLogger
}

// Write will write to logger entry once the http.Request is complete
func (l *structuredLoggerEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, extra interface{}) {
	l.logger.WithFields(logrus.Fields{
		"status":        status,
		"bytes_written": bytes,
		"elapsed_ms":    float64(elapsed.Nanoseconds()) / 1000000.0,
	}).Infoln()
}

// Panic attaches the panic stack and text to the log entry
func (l *structuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	}).Errorln("request panic'd")
}

// LogEntrySetField will set a new field on a log entry
func LogEntrySetField(ctx context.Context, key string, value interface{}) {
	if entry, ok := ctx.Value(middleware.LogEntryCtxKey).(*structuredLoggerEntry); ok {
		entry.logger = entry.logger.WithField(key, value)
	}
}

// LogEntrySetFields will set a map of key/value pairs on a log entry
func LogEntrySetFields(ctx context.Context, fields map[string]interface{}) {
	if entry, ok := ctx.Value(middleware.LogEntryCtxKey).(*structuredLoggerEntry); ok {
		entry.logger = entry.logger.WithFields(fields)
	}
}
