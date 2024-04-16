package gorm

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

func mapSlogLogLevel(logLevelString string) slog.Level {
	switch strings.ToLower(logLevelString) {
	case "silent":
		return slog.LevelError + 10
	case "info":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelWarn
	}
}

type SlogLogger struct {
	delegate *slog.Logger
	logLevel slog.Level
}

func NewSlogLogger(logLevel string) *SlogLogger {
	return &SlogLogger{
		delegate: slog.Default(),
		logLevel: mapSlogLogLevel(logLevel),
	}
}

func (sl *SlogLogger) LogMode(level logger.LogLevel) logger.Interface {
	l := *sl
	switch level {
	case logger.Silent:
		l.logLevel = slog.LevelError + 10
	case logger.Info:
		l.logLevel = slog.LevelInfo
	case logger.Warn:
		l.logLevel = slog.LevelWarn
	case logger.Error:
		l.logLevel = slog.LevelError
	}
	return &l
}

func (sl *SlogLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if sl.logLevel <= slog.LevelInfo {
		sl.delegate.DebugContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (sl *SlogLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if sl.logLevel <= slog.LevelWarn {
		sl.delegate.WarnContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (sl *SlogLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if sl.logLevel <= slog.LevelError {
		sl.delegate.ErrorContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (sl *SlogLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if sl.logLevel <= slog.LevelDebug {
		elapsed := time.Since(begin)
		sql, rows := fc()
		if err != nil {
			sl.delegate.ErrorContext(ctx, fmt.Sprintf("Query: %v, took: %v, rows affected: %v, error: %v", sql, float64(elapsed.Nanoseconds())/1e6, rows, err))
		} else {
			sl.delegate.InfoContext(ctx, fmt.Sprintf("Query: %v, took: %v, rows affected: %v", sql, float64(elapsed.Nanoseconds())/1e6, rows))
		}
	}
}
