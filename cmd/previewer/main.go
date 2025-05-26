package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IKolyas/thumbnailer/internal/app"
	"github.com/IKolyas/thumbnailer/internal/config"
	"github.com/davidbyttow/govips/v2/vips"
)

func init() {
	vips.Startup(nil)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	path, err := config.ParseFlags()
	if err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	// Load config.
	cfg, err := config.Load(path)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	defer cancel()

	// Гарантируем завершение vips при выходе.
	defer func() {
		vips.Shutdown()
	}()

	go func() {
		if err := application.Run(); err != nil {
			application.Logger.Error("Application failed")
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case <-quit:
		application.Logger.Info("Shutting down gracefully...")
		cancel()
	case <-ctx.Done():
		application.Logger.Info("Context canceled, shutting down...")
	}

	application.Stop()
	application.Logger.Info("Application stopped")
}
