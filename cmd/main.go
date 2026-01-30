package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"testovoe/internal/application"
	"testovoe/internal/config"
	"testovoe/internal/domain"
	"testovoe/internal/http/handlers"
	"testovoe/internal/http/router"
	"testovoe/internal/storage"
	"testovoe/internal/usecase"

	_ "testovoe/docs"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := setupLogger(cfg.Env)

	db, err := storage.New(ctx, cfg.Storage.Addr)
	if err != nil {
		log.Error("Failed to connect to storage", "error", err)
		return
	}
	defer db.Close()

	httpRouter := chi.NewRouter()

	useCase := usecase.New(log, db, cfg)

	httpHandlers := handlers.New(log, useCase)

	router.Router(httpRouter, httpHandlers, log)

	app := application.New(ctx, cfg, log, httpRouter)

	app.MustRun()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown

	app.Shutdown()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case domain.EnvLocal:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case domain.EnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case domain.EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
