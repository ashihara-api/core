package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	global *slog.Logger
)

func init() {
	global = NewWith(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
}

func New() *slog.Logger {
	return NewWithLevel(slog.LevelInfo)
}

func NewWithLevel(lv slog.Level) *slog.Logger {
	return NewWith(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     lv,
	}))
}

func NewWith(h slog.Handler) *slog.Logger {
	return slog.New(h)
}

func Debug(msg string, args ...any) {
	global.Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	global.DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	global.Info(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	global.InfoContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	global.Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	global.WarnContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	global.Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	global.ErrorContext(ctx, msg, args...)
}
