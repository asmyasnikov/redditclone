package logger

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/asmyasnikov/redditclone/api"
)

// ILogger interface
type ILogger interface {
	// With returns a logger based off the root logger and decorates it with the given context and arguments.
	With(ctx context.Context, args ...interface{}) *Logger

	// Debug uses fmt.Sprint to construct and logger a message at DEBUG level
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and logger a message at INFO level
	Info(args ...interface{})
	// Error uses fmt.Sprint to construct and logger a message at ERROR level
	Error(args ...interface{})

	// Debugf uses fmt.Sprintf to construct and logger a message at DEBUG level
	Debugf(format string, args ...interface{})
	// Infof uses fmt.Sprintf to construct and logger a message at INFO level
	Infof(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and logger a message at ERROR level
	Errorf(format string, args ...interface{})
	// Sync synchronises logging
	Sync() error
	// Print uses fmt.Sprint to construct and logger a message at DEBUG level
	Print(v ...interface{})
}

// Logger struct
type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) Print(v ...interface{}) {
	l.Debug(v)
}

type contextKey int

const (
	requestIDKey contextKey = iota
	correlationIDKey
)

var defaultZapConfig = zap.Config{
	Encoding: "json",
	EncoderConfig: zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "",
		LineEnding:     "",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: nil,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     nil,
	},
}

// New creates a new logger
func New(conf api.Log) (*Logger, error) {
	cfg, err := configToZapConfig(conf)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not convert conf to zap conf;\nconf: %v", conf)
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, errors.Wrapf(err, "Can not build loger by cfg: %#v", cfg)
	}

	logger := NewWithZap(zapLogger)

	logger.Info("Logger construction succeeded")
	return logger, nil
}

func configToZapConfig(conf api.Log) (zap.Config, error) {
	cfg := defaultZapConfig
	cfg.OutputPaths = conf.OutputPaths
	cfg.Encoding = conf.Encoding
	cfg.InitialFields = make(map[string]interface{}, len(conf.InitialFields))

	for key, val := range conf.InitialFields {
		cfg.InitialFields[key] = val
	}

	if err := cfg.Level.UnmarshalText([]byte(conf.Level)); err != nil {
		return cfg, errors.Wrapf(err, "Can not unmarshal text %q, expected one of zapcore.Levels", conf.Level)
	}

	return cfg, nil
}

// NewByDefault creates a new logger using the default configuration.
func NewByDefault() *Logger {
	l, _ := zap.NewProduction()
	return NewWithZap(l)
}

// NewWithZap creates a new logger using the preconfigured zap logger.
func NewWithZap(l *zap.Logger) *Logger {
	return &Logger{l.Sugar()}
}

// With returns a logger based off the root logger and decorates it with the given context and arguments.
//
// If the context contains request ID and/or correlation ID information (recorded via WithRequestID()
// and WithCorrelationID()), they will be added to every logger message generated by the new logger.
//
// The arguments should be specified as a sequence of name, value pairs with names being strings.
// The arguments will also be added to every logger message generated by the logger.
func (l *Logger) With(ctx context.Context, args ...interface{}) *Logger {
	if ctx != nil {
		if id, ok := ctx.Value(requestIDKey).(string); ok {
			args = append(args, zap.String("request_id", id))
		}
		if id, ok := ctx.Value(correlationIDKey).(string); ok {
			args = append(args, zap.String("correlation_id", id))
		}
	}
	if len(args) > 0 {
		return &Logger{l.SugaredLogger.With(args...)}
	}
	return l
}

// WithRequest returns a context which knows the request ID and correlation ID in the given request.
func WithRequest(ctx context.Context, req *http.Request) context.Context {
	id := getRequestID(req)
	if id == "" {
		id = uuid.New().String()
	}
	ctx = context.WithValue(ctx, requestIDKey, id)
	if id := getCorrelationID(req); id != "" {
		ctx = context.WithValue(ctx, correlationIDKey, id)
	}
	return ctx
}

// getCorrelationID extracts the correlation ID from the HTTP request
func getCorrelationID(req *http.Request) string {
	return req.Header.Get("X-Correlation-ID")
}

// getRequestID extracts the correlation ID from the HTTP request
func getRequestID(req *http.Request) string {
	return req.Header.Get("X-Request-ID")
}
