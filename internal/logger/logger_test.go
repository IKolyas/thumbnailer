package logger

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ctx := context.Background()

	// Тестирование корректного уровня логирования
	logger, err := New(ctx, "info", "")
	assert.NoError(t, err)
	assert.Equal(t, LevelInfo, logger.GetLevel())

	// Тестирование некорректного уровня логирования
	logger, err = New(ctx, "unknown", "")
	assert.Error(t, err)
	assert.Nil(t, logger)
}

func TestSetLevel(t *testing.T) {
	ctx := context.Background()
	logger, _ := New(ctx, "info", "")

	logger.SetLevel(LevelDebug)
	assert.Equal(t, LevelDebug, logger.GetLevel())
}

func TestLog(t *testing.T) {
	ctx := context.Background()
	logger, _ := New(ctx, "info", "")

	// Тестирование вывода сообщения
	var buf bytes.Buffer
	logger.output = &buf
	logger.Log(LevelInfo, "test message")
	assert.Contains(t, buf.String(), "test message")

	// Тестирование игнорирования сообщения
	buf.Reset()
	logger.Log(LevelDebug, "test message")
	assert.Empty(t, buf.String())
}

func TestError(t *testing.T) {
	ctx := context.Background()
	logger, _ := New(ctx, "info", "")

	var buf bytes.Buffer
	logger.output = &buf
	logger.Error("test message")
	assert.Contains(t, buf.String(), "test message")
}

func TestWarn(t *testing.T) {
	ctx := context.Background()
	logger, _ := New(ctx, "info", "")

	var buf bytes.Buffer
	logger.output = &buf
	logger.Warn("test message")
	assert.Contains(t, buf.String(), "test message")
}

func TestInfo(t *testing.T) {
	ctx := context.Background()
	logger, _ := New(ctx, "info", "")

	var buf bytes.Buffer
	logger.output = &buf
	logger.Info("test message")
	assert.Contains(t, buf.String(), "test message")
}

func TestDebug(t *testing.T) {
	ctx := context.Background()
	logger, _ := New(ctx, "info", "")

	var buf bytes.Buffer
	logger.output = &buf
	logger.Debug("test message")
	assert.Empty(t, buf.String())
}
