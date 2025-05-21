package app

import (
	"context"
	"log"
	"time"

	"github.com/IKolyas/image-previewer/internal/config"
	"github.com/IKolyas/image-previewer/internal/logger"
	"github.com/IKolyas/image-previewer/internal/server/http"
	"github.com/IKolyas/image-previewer/internal/storage/memory"
)

type App struct {
	cfg     *config.Config
	server  *http.Server
	storage *memory.LRUStorage
	Logger  *logger.Logger
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	logger, err := logger.New(ctx, cfg.Logger.Level, cfg.Logger.Output)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	storage, err := memory.NewLRUStorage(cfg.CacheCapacity, cfg.StorageDir)
	if err != nil {
		log.Fatalf("Error create lru storage: %v", err)
	}

	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		log.Fatalf("Error parsing duration: %v", err)
	}

	server, err := http.NewServer(
		cfg.Host,
		storage,
		logger,
		http.WithMaxBodySize(cfg.MaxBodySize),
		http.WithTimeout(timeout),
	)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:     cfg,
		server:  server,
		storage: storage,
		Logger:  logger,
	}, nil
}

func (a *App) Run() error {
	a.Logger.Info("Starting application")
	return a.server.Start()
}

func (a *App) Stop() {
	if err := a.server.Stop(); err != nil {
		a.Logger.Error("Failed to stop server")
	}
	a.Logger.Info("Stop application")
	if err := a.storage.Clear(); err != nil {
		a.Logger.Error("Failed to clear cache")
	}
}
