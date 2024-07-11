package logger

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/google/uuid"
	slogmulti "github.com/samber/slog-multi"
)

type LoggerInterface interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, err error, args ...any)
	WithOperation(operationID string) LoggerInterface
}

type Logger struct {
	logger      *slog.Logger
	operationID string
}

var _ LoggerInterface = &Logger{}

func NewLogger(logFilePath string) (LoggerInterface, error) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return nil, err
	}

	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Set stdout to INFO level
	})

	multiHandler := slogmulti.Fanout(fileHandler, stdoutHandler)

	slogLogger := slog.New(multiHandler)

	return &Logger{logger: slogLogger}, nil
}

func (l *Logger) log(level slog.Level, msg string, args ...any) {
	if l.operationID != "" {
		args = append(args, "operationID", l.operationID)
	}
	l.logger.Log(context.TODO(), level, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(slog.LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(slog.LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, err error, args ...any) {
	args = append(args, "error", err)
	l.log(slog.LevelError, msg, args...)
}

func (l *Logger) WithOperation(operationID string) LoggerInterface {
	return &Logger{
		logger:      l.logger,
		operationID: operationID,
	}
}

// Helper function to generate a new operation ID
func NewOperationID() string {
	return uuid.New().String()
}
