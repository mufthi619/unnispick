package logger

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

type Config struct {
	Environment string
	LogLevel    string
}

func NewLogger(cfg Config) (*Logger, error) {
	config := zap.NewProductionConfig()

	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		return nil, fmt.Errorf("parsing log level: %w", err)
	}
	config.Level = zap.NewAtomicLevelAt(level)

	// Determine Env
	if cfg.Environment == "development" {
		config.Encoding = "console"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.Development = true
	} else {
		config.Encoding = "json"
	}

	zapLogger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("building logger: %w", err)
	}

	return &Logger{
		logger: zapLogger,
	}, nil
}

func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return l.logger
	}

	spanContext := span.SpanContext()
	return l.logger.With(
		zap.String("trace_id", spanContext.TraceID().String()),
		zap.String("span_id", spanContext.SpanID().String()),
	)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Debug(msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Info(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Warn(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Error(msg, fields...)
}

func (l *Logger) ErrorWithStack(ctx context.Context, err error, msg string, fields ...zap.Field) {
	allFields := append(fields, zap.Error(err))
	l.WithContext(ctx).Error(msg, allFields...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.logger.Sync()
}
