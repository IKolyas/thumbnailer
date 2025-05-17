package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	LevelError LogLevel = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

type Logger struct {
	mu     sync.Mutex
	level  LogLevel
	output io.Writer
}

func New(ctx context.Context, level string, outputFile string) (*Logger, error) {
	lvl, err := parseLogLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	var output io.Writer = os.Stdout

	if outputFile != "" && outputFile != "stdout" {
		file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		go func() {
			<-ctx.Done()
			if err := file.Close(); err != nil {
				log.Printf("failed to close log file: %v", err)
			}
		}()
		output = file
	}

	return &Logger{
		level:  lvl,
		output: output,
	}, nil
}

func parseLogLevel(level string) (LogLevel, error) {
	switch strings.ToLower(level) {
	case "error":
		return LevelError, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	default:
		return LevelInfo, fmt.Errorf("unknown log level: %s", level)
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	l.level = level
	l.mu.Unlock()
}

func (l *Logger) GetLevel() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *Logger) Log(level LogLevel, msg string) {
	if level > l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	_, err := fmt.Fprintf(l.output, "[%s] [%s] %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		levelToString(level),
		msg)
	if err != nil {
		log.Printf("failed to write log: %v", err)
	}
}

func levelToString(level LogLevel) string {
	switch level {
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

func (l *Logger) Error(msg string) {
	l.Log(LevelError, msg)
}

func (l *Logger) Warn(msg string) {
	l.Log(LevelWarn, msg)
}

func (l *Logger) Info(msg string) {
	l.Log(LevelInfo, msg)
}

func (l *Logger) Debug(msg string) {
	l.Log(LevelDebug, msg)
}
