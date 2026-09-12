package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
	"time"
)

type loggerContextKey struct{}

var key = loggerContextKey{}

type Logger struct {
	*zap.Logger
	file *os.File
}

func ToContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, key, log)
}
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(key).(*Logger)
	if !ok {
		panic("no logger in context")
	}
	return log
}

func NewLogger(config Config) (*Logger, error) {
	zapLevel := zap.NewAtomicLevel()
	if err := zapLevel.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshall log level: %w", err)
	}
	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}
	logFilePath := filepath.Join(config.Folder, fmt.Sprintf("%s.log", time.Now().UTC().Format("2006-01-02T15-04-05.000000")))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}
	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15-04-05.000000")
	encoder := zapcore.NewConsoleEncoder(zapConfig)
	core := zapcore.NewTee(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel), zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapLevel))
	return &Logger{Logger: zap.New(core, zap.AddCaller()), file: logFile}, nil
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...), file: l.file}
}
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		log.Println("failed to close application logger")
	}
}
