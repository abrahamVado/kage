package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kage/backend/internal/app"
)

func main() {
	logger := log.New(os.Stdout, "backend ", log.LstdFlags)

	// 1.- Load configuration, construct dependencies, and build the HTTP server.
	cfg, err := app.LoadConfig()
	if err != nil {
		logger.Fatalf("load config: %v", err)
	}

	application, err := app.Build(cfg, logger)
	if err != nil {
		logger.Fatalf("build application: %v", err)
	}

	srv := &http.Server{Addr: ":" + cfg.HTTPPort, Handler: application.Engine}

	go func() {
		// 2.- Start the HTTP listener and report fatal boot errors.
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server error: %v", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Printf("shutdown error: %v", err)
	}
	if err := application.Cleanup(ctx); err != nil {
		logger.Printf("cleanup error: %v", err)
	}
}
