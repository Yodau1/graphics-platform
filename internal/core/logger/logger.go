package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// loggerContextKey — приватный тип ключа для context.WithValue.
// Использование отдельного типа (а не string) исключает коллизии ключей
// с другими пакетами, которые тоже хранят данные в контексте.
type loggerContextKey struct{}

type Logger struct {
	*zap.Logger
	file *os.File
}

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

func FromContext(ctx context.Context) *Logger {
	logger, ok := ctx.Value(loggerContextKey{}).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return logger
}

func NewLogger(config Config) (*Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("parse logger level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("create log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	file, err := os.CreateTemp(config.Folder, timestamp+"-*.log")
	if err != nil {
		return nil, fmt.Errorf("create log file in %q: %w", filepath.Clean(config.Folder), err)
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(file), level),
	)

	return &Logger{
		Logger: zap.New(core, zap.AddCaller()),
		file:   file,
	}, nil
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		file:   l.file,
	}
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "close log file:", err)
	}
}
